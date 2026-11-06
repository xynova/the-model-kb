package limits

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/xynova/roster-mill/internal/aa"
	derrors "github.com/xynova/roster-mill/internal/errors"
)

const (
	UserAgent   = "the-model-kb-roster-mine/1.0"
	MaxBody     = 2 * 1024 * 1024
	HTTPTimeout = 20 * time.Second
)

var skipHostIDs = map[string]struct{}{"cursor": {}, "local": {}}

// LinkKeys maps host link keys to allowed hostnames.
var LinkKeys = map[string]string{
	"cloudflare_pricing": "developers.cloudflare.com",
	"openrouter":         "openrouter.ai",
	"provider_pricing":   "cloud.google.com",
	"deepinfra":          "deepinfra.com",
	"xinference":         "xinference.co",
}

var (
	evidencePattern = regexp.MustCompile(`(?is).{0,240}(?:neuron|rpm|tpm|rate limit|free tier|quota|production|throughput|flex|priority).{0,240}`)
	htmlTag         = regexp.MustCompile(`<[^>]+>`)
	productionRPM   = regexp.MustCompile(`(?is)production.{0,80}(?:rpm|throughput|high.?throughput)`)
	highThroughput  = regexp.MustCompile(`(?i)high[- ]throughput`)
	productionWord  = regexp.MustCompile(`(?i)\bproduction\b`)
	priceRe         = regexp.MustCompile(`\$\s*(\d+(?:\.\d+)?)`)
)

// FetchURL returns HTTP status and body.
type FetchURL func(url string) (int, string)

func stripHTML(html string) string {
	text := htmlTag.ReplaceAllString(html, " ")
	return strings.Join(strings.Fields(text), " ")
}

func extractEvidence(text string) string {
	match := evidencePattern.FindString(text)
	if match == "" {
		return ""
	}
	snippet := strings.TrimSpace(match)
	if len(snippet) > 240 {
		snippet = strings.TrimRight(snippet[:240], " \t")
	}
	return snippet
}

func hostnameAllowed(rawURL, expectedHost string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	host = strings.TrimPrefix(host, "www.")
	expected := strings.ToLower(expectedHost)
	return host == expected || strings.HasSuffix(host, "."+expected)
}

func pickAllowlistedURL(links map[string]any) (key, u string, ok bool) {
	if links == nil {
		return "", "", false
	}
	// Stable preference matches Python dict insertion order of LINK_KEYS.
	order := []string{
		"cloudflare_pricing", "openrouter", "provider_pricing", "deepinfra", "xinference",
	}
	for _, k := range order {
		expected := LinkKeys[k]
		raw, _ := links[k].(string)
		if raw == "" {
			continue
		}
		if hostnameAllowed(raw, expected) {
			return k, raw, true
		}
	}
	return "", "", false
}

func parseUSDPrices(text string) (in *float64, out *float64) {
	lower := strings.ToLower(text)
	for _, m := range priceRe.FindAllStringSubmatchIndex(text, -1) {
		start := m[0]
		from := start - 40
		if from < 0 {
			from = 0
		}
		to := start + 40
		if to > len(lower) {
			to = len(lower)
		}
		window := lower[from:to]
		val := parseFloat(text[m[2]:m[3]])
		if containsAny(window, "input", "prompt", "in ") {
			in = &val
		} else if containsAny(window, "output", "completion", "out ") {
			out = &val
		}
	}
	if in != nil && out != nil {
		return in, out
	}
	return nil, nil
}

func parseFloat(s string) float64 {
	var f float64
	_, _ = fmt.Sscanf(s, "%f", &f)
	return f
}

func containsAny(s string, needles ...string) bool {
	for _, n := range needles {
		if strings.Contains(s, n) {
			return true
		}
	}
	return false
}

func suggestCapacity(linkKey, text, evidence string) (string, string) {
	lower := strings.ToLower(text)
	evLower := strings.ToLower(evidence)
	switch linkKey {
	case "cloudflare_pricing":
		_ = evLower
		return "hobby", "cloudflare_free_tier"
	case "openrouter":
		if productionRPM.MatchString(text) || highThroughput.MatchString(text) {
			return "industrial", "deepinfra_production_text"
		}
		return "hobby", "openrouter_marketplace"
	case "provider_pricing":
		if productionRPM.MatchString(text) {
			return "industrial", "deepinfra_production_text"
		}
		return "unknown", "google_no_rpm"
	case "deepinfra":
		if productionRPM.MatchString(text) || highThroughput.MatchString(text) {
			return "industrial", "deepinfra_production_text"
		}
		if productionWord.MatchString(text) {
			return "industrial", "deepinfra_production_text"
		}
		_ = lower
		return "unknown", "deepinfra_no_rpm"
	case "xinference":
		return "unknown", "xinference_no_rpm"
	default:
		return "unknown", "fetch_error"
	}
}

// DefaultFetch GETs url with retry on transport error.
func DefaultFetch(rawURL string) (int, string) {
	doGet := func() (int, string, error) {
		req, err := http.NewRequest(http.MethodGet, rawURL, nil)
		if err != nil {
			return 0, "", err
		}
		req.Header.Set("User-Agent", UserAgent)
		client := &http.Client{Timeout: HTTPTimeout}
		resp, err := client.Do(req)
		if err != nil {
			return 0, "", err
		}
		defer resp.Body.Close()
		raw, err := io.ReadAll(io.LimitReader(resp.Body, MaxBody+1))
		if err != nil {
			return resp.StatusCode, "", err
		}
		if len(raw) > MaxBody {
			raw = raw[:MaxBody]
		}
		return resp.StatusCode, string(raw), nil
	}
	status, body, err := doGet()
	if err == nil {
		return status, body
	}
	time.Sleep(2 * time.Second)
	status, body, err = doGet()
	if err != nil {
		return 0, ""
	}
	return status, body
}

func processHostRow(cardMap, host map[string]any, fetch FetchURL, now time.Time) map[string]any {
	name := str(cardMap["name"])
	hostID := str(host["id"])
	base := map[string]any{
		"name":               name,
		"host_id":            hostID,
		"runtime":            host["runtime"],
		"runtime_model_id":   host["runtime_model_id"],
		"url":                "",
		"http_status":        0,
		"fetched_at":         now.Format(time.RFC3339),
		"status":             "skipped",
		"evidence":           "",
		"price_in":           nil,
		"price_out":          nil,
		"suggested_capacity": "unknown",
		"suggestion_reason":  "skip_no_link",
		"current_capacity":   host["capacity"],
		"capacity_applied":   false,
	}
	if cap, ok := host["capacity"].(string); ok {
		base["current_capacity"] = cap
	} else if host["capacity"] == nil {
		base["current_capacity"] = "unknown"
	}
	if str(cardMap["kind"]) == "service" {
		base["suggestion_reason"] = "skip_service"
		return base
	}
	if _, skip := skipHostIDs[hostID]; skip || str(host["capacity"]) == "local" {
		base["suggestion_reason"] = "skip_local"
		return base
	}
	links, _ := host["links"].(map[string]any)
	key, u, ok := pickAllowlistedURL(links)
	if !ok {
		base["suggestion_reason"] = "skip_no_link"
		return base
	}
	base["url"] = u
	status, body := fetch(u)
	base["http_status"] = status
	if status != 200 || body == "" {
		base["status"] = "error"
		base["suggestion_reason"] = "fetch_error"
		return base
	}
	text := stripHTML(body)
	evidence := extractEvidence(text)
	priceIn, priceOut := parseUSDPrices(text)
	suggested, reason := suggestCapacity(key, text, evidence)
	base["status"] = "ok"
	base["evidence"] = evidence
	if priceIn != nil {
		base["price_in"] = *priceIn
	}
	if priceOut != nil {
		base["price_out"] = *priceOut
	}
	base["suggested_capacity"] = suggested
	base["suggestion_reason"] = reason
	return base
}

func catalogCards(catalog map[string]any) []map[string]any {
	raw, _ := catalog["cards"].([]any)
	out := make([]map[string]any, 0, len(raw))
	for _, item := range raw {
		if m, ok := item.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

func hostMaps(cardMap map[string]any) []map[string]any {
	raw, _ := cardMap["hosts"].([]any)
	out := make([]map[string]any, 0, len(raw))
	for _, item := range raw {
		if m, ok := item.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

// Retrieve walks catalog hosts and fetches allowlisted public pages.
func Retrieve(catalog map[string]any, fetch FetchURL, now time.Time, sleep func(time.Duration)) map[string]any {
	if fetch == nil {
		fetch = DefaultFetch
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if sleep == nil {
		sleep = time.Sleep
	}
	var rows []any
	for _, c := range catalogCards(catalog) {
		for _, host := range hostMaps(c) {
			row := processHostRow(c, host, fetch, now)
			rows = append(rows, row)
			if str(row["status"]) == "ok" && str(row["url"]) != "" {
				sleep(time.Second)
			}
		}
	}
	return map[string]any{
		"review_schema": 1,
		"generated_at":  now.UTC().Format("2006-01-02"),
		"rows":          rows,
	}
}

// Run is the CLI retrieve-limits command.
func Run(repoRoot string, stdout io.Writer) error {
	catalogPath := filepath.Join(repoRoot, "roster", "catalog.json")
	reviewPath := filepath.Join(repoRoot, "roster", "retrieval", "limits-review.json")
	raw, err := os.ReadFile(catalogPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "catalog missing: %s\n", catalogPath)
		fmt.Fprintln(os.Stderr, "Run: make -C site generate-roster")
		return derrors.Wrap(err, derrors.CodeNotFound, "retrieve-limits", "catalog missing")
	}
	var catalog map[string]any
	if err := json.Unmarshal(raw, &catalog); err != nil {
		return derrors.Wrap(err, derrors.CodeFailed, "retrieve-limits", "parse catalog")
	}
	payload := Retrieve(catalog, nil, time.Now().UTC(), nil)
	if err := aa.WriteJSON(reviewPath, payload); err != nil {
		return err
	}
	rows, _ := payload["rows"].([]any)
	fmt.Fprintf(stdout, "wrote %s (%d rows)\n", reviewPath, len(rows))
	return nil
}

func str(v any) string {
	s, _ := v.(string)
	return s
}
