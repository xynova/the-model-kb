package polypus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	derrors "github.com/xynova/roster-mill/internal/errors"
)

const (
	DefaultBaseURL = "http://127.0.0.1:1320"
	UserAgent      = "the-model-kb-polypus-catalog-diff/1.0"
	Timeout        = 20 * time.Second
)

var apiPaths = map[string]string{
	"openai":    "/v1/apis/openai/models",
	"systemone": "/v1/apis/systemone/models",
}

var routePrefixToHost = map[string]string{
	"cf_local":    "cloudflare",
	"deepinfra":   "deepinfra",
	"openrouter":  "openrouter",
	"gemini_free": "google",
	"gemini_paid": "google",
}

var skipNeedles = []string{
	"/bge-", "/bge/", "bge-m3", "bge-small", "whisper", "deepgram",
	"kokoro", "fish-audio", "transcribe", "aura-2", "smart-turn", "investigator",
}

var stemSuffixes = []string{"-maas", "-it", "-instruct"}

// FetchURL returns HTTP status and body.
type FetchURL func(url string) (int, string)

// ModelsURL builds a Polypus models list URL.
func ModelsURL(baseURL, api string) string {
	raw := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if raw == "" {
		raw = DefaultBaseURL
	}
	if strings.HasSuffix(raw, "/apis/openai/models") || strings.HasSuffix(raw, "/apis/systemone/models") {
		return raw
	}
	path := apiPaths[api]
	if path == "" {
		path = apiPaths["openai"]
	}
	if strings.HasSuffix(raw, "/v1") {
		return raw + strings.TrimPrefix(path, "/v1")
	}
	return raw + path
}

// DefaultModelsURLs returns the openai + systemone pair.
func DefaultModelsURLs(baseURL string) []string {
	return []string{ModelsURL(baseURL, "openai"), ModelsURL(baseURL, "systemone")}
}

func lastComponent(value string) string {
	text := strings.ToLower(strings.TrimSpace(value))
	if i := strings.LastIndex(text, "/"); i >= 0 {
		return text[i+1:]
	}
	return text
}

func stem(value string) string {
	text := lastComponent(value)
	changed := true
	for changed {
		changed = false
		for _, suffix := range stemSuffixes {
			if strings.HasSuffix(text, suffix) && len(text) > len(suffix) {
				text = text[:len(text)-len(suffix)]
				changed = true
			}
		}
	}
	return text
}

func isSkipped(modelID string) bool {
	lowered := strings.ToLower(modelID)
	if strings.HasPrefix(lowered, "router/") {
		return true
	}
	for _, needle := range skipNeedles {
		if strings.Contains(lowered, needle) {
			return true
		}
	}
	return false
}

func splitRoute(modelID string) (prefix, hostID, remainder string) {
	text := strings.TrimSpace(modelID)
	if !strings.Contains(text, "/") {
		return "", "", text
	}
	prefix, remainder, _ = strings.Cut(text, "/")
	host := routePrefixToHost[prefix]
	if host != "" {
		return prefix, host, remainder
	}
	return "", "", text
}

// ParseModelsPayload extracts unique model ids from an OpenAI-style list.
func ParseModelsPayload(payload map[string]any) ([]string, error) {
	if payload == nil {
		return nil, fmt.Errorf("models payload must be an object")
	}
	data, ok := payload["data"].([]any)
	if !ok {
		return nil, fmt.Errorf("models payload missing data list")
	}
	var ids []string
	seen := map[string]struct{}{}
	for _, item := range data {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		modelID, _ := m["id"].(string)
		modelID = strings.TrimSpace(modelID)
		if modelID == "" {
			continue
		}
		if _, dup := seen[modelID]; dup {
			continue
		}
		seen[modelID] = struct{}{}
		ids = append(ids, modelID)
	}
	return ids, nil
}

func catalogRuntimeIndex(catalog map[string]any) []map[string]any {
	var rows []map[string]any
	cards, _ := catalog["cards"].([]any)
	for _, item := range cards {
		cardMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		name, _ := cardMap["name"].(string)
		if name == "" {
			continue
		}
		hosts, _ := cardMap["hosts"].([]any)
		for _, h := range hosts {
			host, ok := h.(map[string]any)
			if !ok {
				continue
			}
			rid, _ := host["runtime_model_id"].(string)
			rid = strings.TrimSpace(rid)
			if rid == "" {
				continue
			}
			hostID, _ := host["id"].(string)
			rows = append(rows, map[string]any{
				"name":             name,
				"host_id":          hostID,
				"runtime_model_id": rid,
			})
		}
	}
	return rows
}

func runtimeMatches(remainder, runtimeModelID string) bool {
	left := strings.ToLower(strings.TrimSpace(remainder))
	right := strings.ToLower(strings.TrimSpace(runtimeModelID))
	if left == "" || right == "" {
		return false
	}
	if left == right {
		return true
	}
	if lastComponent(left) == lastComponent(right) && lastComponent(left) != "" {
		return true
	}
	ls, rs := stem(left), stem(right)
	return ls != "" && ls == rs
}

func matchCatalogRows(remainder string, catalogRows []map[string]any) []map[string]any {
	var hits []map[string]any
	seen := map[string]struct{}{}
	for _, row := range catalogRows {
		if !runtimeMatches(remainder, str(row["runtime_model_id"])) {
			continue
		}
		key := str(row["name"]) + "\x00" + str(row["host_id"])
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		hits = append(hits, row)
	}
	return hits
}

// ClassifyModels diffs Polypus ids against catalog runtime_model_id values.
func ClassifyModels(modelIDs []string, catalog map[string]any) map[string]any {
	catalogRows := catalogRuntimeIndex(catalog)
	var matched, proposed, skipped []any
	seenPolypus := map[string]struct{}{}
	for _, modelID := range modelIDs {
		if _, dup := seenPolypus[modelID]; dup {
			continue
		}
		seenPolypus[modelID] = struct{}{}
		prefix, hostID, remainder := splitRoute(modelID)
		if isSkipped(modelID) {
			skipped = append(skipped, map[string]any{
				"id":               modelID,
				"reason":           "specialty",
				"host_id":          hostID,
				"runtime_model_id": remainder,
			})
			continue
		}
		hits := matchCatalogRows(remainder, catalogRows)
		if len(hits) == 0 && remainder != modelID {
			hits = matchCatalogRows(modelID, catalogRows)
		}
		if hostID != "" {
			var sameHost []map[string]any
			for _, row := range hits {
				if str(row["host_id"]) == hostID {
					sameHost = append(sameHost, row)
				}
			}
			if len(sameHost) > 0 {
				hits = sameHost
			}
		}
		if len(hits) > 0 {
			cardAny := make([]any, len(hits))
			for i, h := range hits {
				cardAny[i] = h
			}
			matched = append(matched, map[string]any{
				"id":               modelID,
				"host_id":          hostID,
				"runtime_model_id": remainder,
				"cards":            cardAny,
			})
			continue
		}
		hostOut := hostID
		if hostOut == "" {
			hostOut = "unknown"
		}
		proposed = append(proposed, map[string]any{
			"id":               modelID,
			"prefix":           prefix,
			"host_id":          hostOut,
			"runtime_model_id": remainder,
			"action":           "model-roster keep/drop; do not mint polypus as a host",
		})
	}
	var catalogUnseen []any
	for _, row := range catalogRows {
		found := false
		for _, modelID := range modelIDs {
			_, _, remainder := splitRoute(modelID)
			if runtimeMatches(remainder, str(row["runtime_model_id"])) ||
				runtimeMatches(modelID, str(row["runtime_model_id"])) {
				found = true
				break
			}
		}
		if !found {
			catalogUnseen = append(catalogUnseen, row)
		}
	}
	if matched == nil {
		matched = []any{}
	}
	if proposed == nil {
		proposed = []any{}
	}
	if skipped == nil {
		skipped = []any{}
	}
	if catalogUnseen == nil {
		catalogUnseen = []any{}
	}
	return map[string]any{
		"diff_schema":    1,
		"matched":        matched,
		"proposed":       proposed,
		"skipped":        skipped,
		"catalog_unseen": catalogUnseen,
		"counts": map[string]any{
			"polypus":             len(modelIDs),
			"matched":             len(matched),
			"proposed":            len(proposed),
			"skipped":             len(skipped),
			"catalog_runtime_ids": len(catalogRows),
			"catalog_unseen":      len(catalogUnseen),
		},
	}
}

func defaultFetch(rawURL string) (int, string) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return 0, ""
	}
	req.Header.Set("User-Agent", UserAgent)
	client := &http.Client{Timeout: Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return 0, ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body)
}

func loadModels(url string, fetch FetchURL) ([]string, error) {
	if fetch == nil {
		fetch = defaultFetch
	}
	status, body := fetch(url)
	if status != 200 || strings.TrimSpace(body) == "" {
		return nil, fmt.Errorf("polypus models HTTP %d at %s", status, url)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return nil, fmt.Errorf("polypus models JSON parse failed: %w", err)
	}
	return ParseModelsPayload(payload)
}

func loadModelsMany(urls []string, fetch FetchURL) ([]string, error) {
	var ids []string
	seen := map[string]struct{}{}
	for _, u := range urls {
		got, err := loadModels(u, fetch)
		if err != nil {
			return nil, err
		}
		for _, id := range got {
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func renderHuman(report map[string]any, urlLabel, catalogPath string) string {
	counts, _ := report["counts"].(map[string]any)
	var b strings.Builder
	fmt.Fprintln(&b, "Polypus catalog diff")
	fmt.Fprintf(&b, "url: %s\n", urlLabel)
	fmt.Fprintf(&b, "catalog: %s\n", catalogPath)
	fmt.Fprintf(&b, "polypus %v | matched %v | proposed %v | skipped %v | catalog unseen %v\n",
		counts["polypus"], counts["matched"], counts["proposed"], counts["skipped"], counts["catalog_unseen"])
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "Prices and capacity stay on public host pages. This list is keep/drop only.")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "## Proposed (not in catalog)")
	proposed, _ := report["proposed"].([]any)
	if len(proposed) == 0 {
		fmt.Fprintln(&b, "(none)")
	}
	for _, item := range proposed {
		row, _ := item.(map[string]any)
		fmt.Fprintf(&b, "- %s\n", row["id"])
		fmt.Fprintf(&b, "  host: %s\n", row["host_id"])
		fmt.Fprintf(&b, "  runtime_model_id: %s\n", row["runtime_model_id"])
		fmt.Fprintf(&b, "  action: %s\n", row["action"])
	}
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "## Matched")
	matched, _ := report["matched"].([]any)
	if len(matched) == 0 {
		fmt.Fprintln(&b, "(none)")
	}
	for _, item := range matched {
		row, _ := item.(map[string]any)
		cards, _ := row["cards"].([]any)
		var parts []string
		for _, c := range cards {
			cm, _ := c.(map[string]any)
			parts = append(parts, fmt.Sprintf("%s / %s", cm["name"], cm["host_id"]))
		}
		fmt.Fprintf(&b, "- %s -> %s\n", row["id"], strings.Join(parts, ", "))
	}
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "## Catalog runtime ids not on this gateway")
	unseen, _ := report["catalog_unseen"].([]any)
	if len(unseen) == 0 {
		fmt.Fprintln(&b, "(none)")
	}
	for _, item := range unseen {
		row, _ := item.(map[string]any)
		fmt.Fprintf(&b, "- %s / %s (%s)\n", row["name"], row["host_id"], row["runtime_model_id"])
	}
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "## Skipped (embeddings, audio, router)")
	skipped, _ := report["skipped"].([]any)
	if len(skipped) == 0 {
		fmt.Fprintln(&b, "(none)")
	}
	for _, item := range skipped {
		row, _ := item.(map[string]any)
		fmt.Fprintf(&b, "- %s\n", row["id"])
	}
	fmt.Fprintln(&b)
	return b.String()
}

// Run is the CLI polypus-diff command.
func Run(repoRoot, catalogPath, baseURL, fullURL string, asJSON bool, stdout io.Writer) error {
	if catalogPath == "" {
		catalogPath = repoRoot + "/roster/catalog.json"
	}
	raw, err := os.ReadFile(catalogPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "catalog missing: %s\n", catalogPath)
		fmt.Fprintln(os.Stderr, "Run: make -C site generate-roster")
		return derrors.Wrap(err, derrors.CodeNotFound, "polypus-diff", "catalog missing")
	}
	var catalog map[string]any
	if err := json.Unmarshal(raw, &catalog); err != nil {
		return derrors.Wrap(err, derrors.CodeFailed, "polypus-diff", "parse catalog")
	}
	var urls []string
	if strings.TrimSpace(fullURL) != "" {
		urls = []string{strings.TrimSpace(fullURL)}
	} else {
		base := strings.TrimSpace(baseURL)
		if base == "" {
			base = os.Getenv("POLYPUS_BASE_URL")
		}
		if base == "" {
			base = DefaultBaseURL
		}
		urls = DefaultModelsURLs(base)
	}
	urlLabel := strings.Join(urls, " ")
	modelIDs, err := loadModelsMany(urls, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return derrors.Wrap(err, derrors.CodeUnavailable, "polypus-diff", "fetch models")
	}
	report := ClassifyModels(modelIDs, catalog)
	report["url"] = urlLabel
	report["urls"] = urls
	if asJSON {
		enc := json.NewEncoder(stdout)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		return enc.Encode(report)
	}
	fmt.Fprint(stdout, renderHuman(report, urlLabel, catalogPath))
	return nil
}

func str(v any) string {
	s, _ := v.(string)
	return s
}
