package aa

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	derrors "github.com/xynova/roster-mill/internal/errors"
)

const (
	FreeURL             = "https://artificialanalysis.ai/api/v2/language/models/free"
	UserAgent           = "the-model-kb-roster-mine/1.0"
	MaxBody             = 2 * 1024 * 1024
	MaxPages            = 50
	SupportedIndexMajor = 5
	CandleMinN          = 3
	HTTPTimeout         = 30 * time.Second
)

var (
	effortParen = regexp.MustCompile(`(?i)[\s_\-]*\((max|xhigh|high|medium|low|reasoning|non-reasoning)\)\s*$`)
	effortSlug  = regexp.MustCompile(`(?i)-(max|xhigh|high|medium|low|reasoning|non-reasoning)$`)
	modelPage   = regexp.MustCompile(`(?i)https?://(?:www\.)?artificialanalysis\.ai/models/(?:releases/)?([a-z0-9\-]+)`)
	nonAlnum    = regexp.MustCompile(`[^a-z0-9]+`)
)

var effortRank = map[string]int{
	"max": 50, "xhigh": 40, "high": 30, "reasoning": 25,
	"": 20, "medium": 10, "low": 0,
}

var skipEfforts = map[string]struct{}{"non-reasoning": {}}

// EvaluationKeys is the AA evaluations object field order.
var EvaluationKeys = []string{
	"artificial_analysis_intelligence_index",
	"artificial_analysis_coding_index",
	"artificial_analysis_agentic_index",
	"artificial_analysis_finance_and_accounting_index",
	"artificial_analysis_strategy_and_ops_index",
	"artificial_analysis_legal_index",
	"artificial_analysis_healthcare_and_medical_index",
	"artificial_analysis_engineering_index",
	"artificial_analysis_economics_index",
}

// PerformanceKeys is the AA performance object field order.
var PerformanceKeys = []string{
	"median_output_tokens_per_second",
	"median_time_to_first_token_seconds",
	"median_time_to_first_answer_token_seconds",
	"median_end_to_end_response_time_seconds",
}

// PricingKeys is the AA pricing object field order.
var PricingKeys = []string{
	"price_1m_input_tokens",
	"price_1m_output_tokens",
}

// FetchPage loads one AA free-API page.
type FetchPage func(page int, apiKey string) (map[string]any, error)

// YAMLScore returns an int when the value is an integer number, else float64, else nil.
func YAMLScore(value any) any {
	if value == nil {
		return nil
	}
	switch t := value.(type) {
	case bool:
		return nil
	case int:
		return t
	case int8:
		return int(t)
	case int16:
		return int(t)
	case int32:
		return int(t)
	case int64:
		return int(t)
	case uint:
		return int(t)
	case uint64:
		return int(t)
	case float32:
		return compactNum(float64(t))
	case float64:
		return compactNum(t)
	case json.Number:
		f, err := t.Float64()
		if err != nil {
			return nil
		}
		return compactNum(f)
	default:
		return nil
	}
}

func compactNum(value float64) any {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return nil
	}
	if value == math.Trunc(value) && math.Abs(value) < 1e15 {
		return int(value)
	}
	return value
}

func CompactNum(value float64) any {
	if value == math.Trunc(value) {
		return int(value)
	}
	return math.Round(value*10000) / 10000
}

func asFloat(v any) (float64, bool) {
	s := YAMLScore(v)
	if s == nil {
		return 0, false
	}
	switch t := s.(type) {
	case int:
		return float64(t), true
	case float64:
		return t, true
	default:
		return 0, false
	}
}

// NormalizeKey lowercases, turns dots into hyphens, then strips other punctuation.
func NormalizeKey(value string) string {
	text := strings.ToLower(strings.ReplaceAll(value, ".", "-"))
	return strings.Trim(nonAlnum.ReplaceAllString(text, "-"), "-")
}

// ParseEffort splits an AA display name / slug into a match key and effort label.
func ParseEffort(name, slug string) (key, effort string) {
	label := strings.TrimSpace(name)
	if m := effortParen.FindStringSubmatchIndex(label); m != nil {
		effort = strings.ToLower(label[m[2]:m[3]])
		label = label[:m[0]]
	}
	slugText := strings.TrimSpace(slug)
	if m := effortSlug.FindStringSubmatchIndex(slugText); m != nil {
		if effort == "" {
			effort = strings.ToLower(slugText[m[2]:m[3]])
		}
		slugText = slugText[:m[0]]
	}
	key = NormalizeKey(slugText)
	if key == "" {
		key = NormalizeKey(label)
	}
	return key, effort
}

func rankEffort(effort string) int {
	if v, ok := effortRank[effort]; ok {
		return v
	}
	return 15
}

func yamlCostPerTask(value any) any {
	if m, ok := value.(map[string]any); ok {
		if task, ok := m["cost_per_task"].(map[string]any); ok {
			if s := YAMLScore(task["total_cost"]); s != nil {
				f, _ := asFloat(s)
				return f
			}
		}
		if s := YAMLScore(m["total_cost"]); s != nil {
			f, _ := asFloat(s)
			return f
		}
		return nil
	}
	s := YAMLScore(value)
	if s == nil {
		return nil
	}
	f, _ := asFloat(s)
	return f
}

func creatorName(model map[string]any) string {
	creator, _ := model["model_creator"].(map[string]any)
	if creator == nil {
		return ""
	}
	name, _ := creator["name"].(string)
	return strings.TrimSpace(name)
}

// ModelPageURL is the public AA model page.
func ModelPageURL(slug string) string {
	return "https://artificialanalysis.ai/models/" + slug
}

func extractOverview(model map[string]any) map[string]any {
	evaluations, _ := model["evaluations"].(map[string]any)
	if evaluations == nil {
		evaluations = map[string]any{}
	}
	pricing, _ := model["pricing"].(map[string]any)
	if pricing == nil {
		pricing = map[string]any{}
	}
	performance, _ := model["performance"].(map[string]any)
	if performance == nil {
		performance = map[string]any{}
	}
	evals := map[string]any{}
	for _, key := range EvaluationKeys {
		evals[key] = YAMLScore(evaluations[key])
	}
	prices := map[string]any{}
	for _, key := range PricingKeys {
		prices[key] = YAMLScore(pricing[key])
	}
	perf := map[string]any{}
	for _, key := range PerformanceKeys {
		perf[key] = YAMLScore(performance[key])
	}
	return map[string]any{
		"evaluations": evals,
		"pricing":     prices,
		"performance": perf,
	}
}

func intelligencePerCost(score, cost any) *float64 {
	sf, okS := asFloat(score)
	cf, okC := asFloat(cost)
	if !okS || !okC || cf <= 0 {
		return nil
	}
	v := sf / cf
	return &v
}

func comparePeer(candidate map[string]any) map[string]any {
	return map[string]any{
		"name":            candidate["aa_name"],
		"creator":         candidate["creator"],
		"slug":            candidate["aa_slug"],
		"effort":          candidate["effort"],
		"aa_intelligence": candidate["aa_intelligence"],
		"cost_per_task":   candidate["cost_per_task"],
		"url":             ModelPageURL(fmt.Sprint(candidate["aa_slug"])),
	}
}

func thinkingCandidates(models []map[string]any) []map[string]any {
	var out []map[string]any
	for _, model := range models {
		slug := fmt.Sprint(model["slug"])
		if model["slug"] == nil {
			slug = ""
		}
		name := fmt.Sprint(model["name"])
		if model["name"] == nil {
			name = ""
		}
		key, effort := ParseEffort(name, slug)
		if key == "" {
			continue
		}
		if _, skip := skipEfforts[effort]; skip {
			continue
		}
		evaluations, _ := model["evaluations"].(map[string]any)
		var score any
		if evaluations != nil {
			score = YAMLScore(evaluations["artificial_analysis_intelligence_index"])
		}
		if score == nil {
			continue
		}
		out = append(out, map[string]any{
			"key":             key,
			"aa_slug":         slug,
			"aa_name":         name,
			"effort":          effort,
			"creator":         creatorName(model),
			"aa_intelligence": score,
			"cost_per_task":   yamlCostPerTask(model["artificial_analysis_intelligence_index_cost"]),
			"rank":            rankEffort(effort),
		})
	}
	return out
}

func pickAlternative(current map[string]any, pool []map[string]any, rosterKeys map[string]struct{}) map[string]any {
	currentScore, okS := asFloat(current["aa_intelligence"])
	currentCost, okC := asFloat(current["cost_per_task"])
	if !okS || !okC || currentCost <= 0 {
		return nil
	}
	var best map[string]any
	var bestIPC *float64
	curKey, _ := current["key"].(string)
	for _, cand := range pool {
		ck, _ := cand["key"].(string)
		if _, onRoster := rosterKeys[ck]; onRoster || ck == curKey {
			continue
		}
		score := YAMLScore(cand["aa_intelligence"])
		cost := cand["cost_per_task"]
		ipc := intelligencePerCost(score, cost)
		cf, okCost := asFloat(cost)
		sf, okScore := asFloat(score)
		if score == nil || ipc == nil || !okCost {
			continue
		}
		if !okScore || sf < currentScore || cf >= currentCost {
			continue
		}
		if best == nil || bestIPC == nil || *ipc > *bestIPC {
			best = cand
			bestIPC = ipc
			continue
		}
		if *ipc == *bestIPC {
			bestScore, _ := asFloat(best["aa_intelligence"])
			if sf > bestScore {
				best = cand
			}
		}
	}
	if best == nil {
		return nil
	}
	return comparePeer(best)
}

func pickSimilar(current map[string]any, pool []map[string]any, rosterKeys map[string]struct{}, skipSlug string) map[string]any {
	currentScore, ok := asFloat(current["aa_intelligence"])
	if !ok {
		return nil
	}
	var best map[string]any
	var bestDelta *float64
	curKey, _ := current["key"].(string)
	for _, cand := range pool {
		ck, _ := cand["key"].(string)
		if _, onRoster := rosterKeys[ck]; onRoster || ck == curKey {
			continue
		}
		if skipSlug != "" && fmt.Sprint(cand["aa_slug"]) == skipSlug {
			continue
		}
		score := YAMLScore(cand["aa_intelligence"])
		sf, okS := asFloat(score)
		if score == nil || !okS {
			continue
		}
		delta := math.Abs(sf - currentScore)
		if best == nil || bestDelta == nil || delta < *bestDelta {
			best = cand
			d := delta
			bestDelta = &d
			continue
		}
		if delta == *bestDelta {
			bestScore, _ := asFloat(best["aa_intelligence"])
			if sf > bestScore {
				best = cand
			}
		}
	}
	if best == nil {
		return nil
	}
	return comparePeer(best)
}

func linkSlug(cardMap map[string]any) string {
	links, _ := cardMap["links"].(map[string]any)
	if links == nil {
		return ""
	}
	u, _ := links["artificial_analysis"].(string)
	m := modelPage.FindStringSubmatch(u)
	if m == nil {
		return ""
	}
	key, _ := ParseEffort("", m[1])
	return key
}

func pickBest(models []map[string]any) map[string]map[string]any {
	best := map[string]map[string]any{}
	for _, model := range models {
		slug := strOrEmpty(model["slug"])
		name := strOrEmpty(model["name"])
		key, effort := ParseEffort(name, slug)
		if key == "" {
			continue
		}
		if _, skip := skipEfforts[effort]; skip {
			continue
		}
		evaluations, _ := model["evaluations"].(map[string]any)
		var score any
		if evaluations != nil {
			score = YAMLScore(evaluations["artificial_analysis_intelligence_index"])
		}
		candidate := map[string]any{
			"key":             key,
			"aa_slug":         slug,
			"aa_name":         name,
			"effort":          effort,
			"creator":         creatorName(model),
			"aa_intelligence": score,
			"cost_per_task":   yamlCostPerTask(model["artificial_analysis_intelligence_index_cost"]),
			"rank":            rankEffort(effort),
			"overview":        extractOverview(model),
		}
		current := best[key]
		if current == nil {
			best[key] = candidate
			continue
		}
		if candidate["rank"].(int) > current["rank"].(int) {
			best[key] = candidate
			continue
		}
		if candidate["rank"].(int) == current["rank"].(int) &&
			current["aa_intelligence"] == nil && candidate["aa_intelligence"] != nil {
			best[key] = candidate
		}
	}
	return best
}

func cardKeys(cardMap map[string]any) []string {
	var keys []string
	if name, ok := cardMap["name"].(string); ok && name != "" {
		keys = append(keys, NormalizeKey(name))
	}
	if title, ok := cardMap["title"].(string); ok && title != "" {
		keys = append(keys, NormalizeKey(title))
	}
	if linked := linkSlug(cardMap); linked != "" {
		keys = append(keys, linked)
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, key)
	}
	return out
}

func asMaps(v any) []map[string]any {
	switch t := v.(type) {
	case []map[string]any:
		return t
	case []any:
		out := make([]map[string]any, 0, len(t))
		for _, item := range t {
			if m, ok := item.(map[string]any); ok {
				out = append(out, m)
			}
		}
		return out
	default:
		return nil
	}
}

func catalogCards(catalog map[string]any) []map[string]any {
	if catalog == nil {
		return nil
	}
	return asMaps(catalog["cards"])
}

// MatchCards pairs catalog cards with AA models.
func MatchCards(catalog map[string]any, aaModels []map[string]any) []map[string]any {
	best := pickBest(aaModels)
	pool := thinkingCandidates(aaModels)
	rosterKeys := map[string]struct{}{}
	for _, c := range catalogCards(catalog) {
		if strOrEmpty(c["kind"]) == "service" {
			continue
		}
		for _, k := range cardKeys(c) {
			rosterKeys[k] = struct{}{}
		}
	}
	var rows []map[string]any
	for _, c := range catalogCards(catalog) {
		name, _ := c["name"].(string)
		if name == "" {
			continue
		}
		row := map[string]any{
			"name":             name,
			"status":           "unmatched",
			"aa_slug":          "",
			"aa_name":          "",
			"effort":           "",
			"aa_creator":       "",
			"aa_cost_per_task": nil,
			"aa_intelligence":  nil,
			"alternative":      nil,
			"similar":          nil,
			"overview":         nil,
			"previous":         c["aa_intelligence"],
			"applied":          false,
		}
		if strOrEmpty(c["kind"]) == "service" {
			row["status"] = "skipped_service"
			rows = append(rows, row)
			continue
		}
		var hit map[string]any
		for _, key := range cardKeys(c) {
			if h, ok := best[key]; ok {
				hit = h
				break
			}
		}
		if hit == nil {
			rows = append(rows, row)
			continue
		}
		row["aa_slug"] = hit["aa_slug"]
		row["aa_name"] = hit["aa_name"]
		row["effort"] = hit["effort"]
		if creator, ok := hit["creator"].(string); ok {
			row["aa_creator"] = creator
		} else {
			row["aa_creator"] = ""
		}
		row["aa_cost_per_task"] = hit["cost_per_task"]
		row["aa_intelligence"] = hit["aa_intelligence"]
		row["overview"] = hit["overview"]
		if hit["aa_intelligence"] == nil {
			row["status"] = "no_score"
		} else {
			row["status"] = "ok"
			alt := pickAlternative(hit, pool, rosterKeys)
			row["alternative"] = asAny(alt)
			skipSlug := ""
			if alt != nil {
				skipSlug = fmt.Sprint(alt["slug"])
			}
			row["similar"] = asAny(pickSimilar(hit, pool, rosterKeys, skipSlug))
		}
		rows = append(rows, row)
	}
	return rows
}

func strOrEmpty(v any) string {
	s, _ := v.(string)
	return s
}

func asAny(m map[string]any) any {
	if m == nil {
		return nil
	}
	return m
}

// DefaultFetchPage GETs one AA free page.
func DefaultFetchPage(page int, apiKey string) (map[string]any, error) {
	u := FreeURL + "?" + url.Values{"page": {fmt.Sprintf("%d", page)}}.Encode()
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", apiKey)
	client := &http.Client{Timeout: HTTPTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, MaxBody+1))
	if err != nil {
		return nil, err
	}
	if len(raw) > MaxBody {
		raw = raw[:MaxBody]
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func collectModels(apiKey string, fetchPage FetchPage, sleep func(time.Duration)) (any, []map[string]any, error) {
	if fetchPage == nil {
		fetchPage = DefaultFetchPage
	}
	if sleep == nil {
		sleep = time.Sleep
	}
	var models []map[string]any
	var version any
	page := 1
	for page <= MaxPages {
		payload, err := fetchPage(page, apiKey)
		if err != nil {
			return nil, nil, err
		}
		if version == nil {
			version = payload["intelligence_index_version"]
		}
		chunk, _ := payload["data"].([]any)
		for _, item := range chunk {
			if m, ok := item.(map[string]any); ok {
				models = append(models, m)
			}
		}
		pagination, _ := payload["pagination"].(map[string]any)
		more, _ := pagination["has_more"].(bool)
		if !more {
			break
		}
		page++
		sleep(200 * time.Millisecond)
	}
	return version, models, nil
}

// Retrieve builds aa-intelligence-review.json from catalog + AA API.
func Retrieve(catalog map[string]any, apiKey string, fetchPage FetchPage, now time.Time, sleep func(time.Duration)) map[string]any {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	generated := now.UTC().Format("2006-01-02")
	if strings.TrimSpace(apiKey) == "" {
		return map[string]any{
			"review_schema":              1,
			"generated_at":               generated,
			"source":                     FreeURL,
			"status":                     "skipped_no_key",
			"intelligence_index_version": nil,
			"rows":                       []any{},
		}
	}
	version, models, err := collectModels(apiKey, fetchPage, sleep)
	if err != nil {
		return map[string]any{
			"review_schema":              1,
			"generated_at":               generated,
			"source":                     FreeURL,
			"status":                     "error",
			"error":                      err.Error(),
			"intelligence_index_version": nil,
			"rows":                       []any{},
		}
	}
	major := indexMajor(version)
	if major != SupportedIndexMajor {
		return map[string]any{
			"review_schema":              1,
			"generated_at":               generated,
			"source":                     FreeURL,
			"status":                     "skipped_version",
			"intelligence_index_version": version,
			"rows":                       []any{},
		}
	}
	rows := MatchCards(catalog, models)
	anyRows := make([]any, len(rows))
	for i, r := range rows {
		anyRows[i] = r
	}
	return map[string]any{
		"review_schema":              1,
		"generated_at":               generated,
		"source":                     FreeURL,
		"status":                     "ok",
		"intelligence_index_version": version,
		"rows":                       anyRows,
	}
}

func indexMajor(version any) int {
	f, ok := asFloat(version)
	if !ok {
		return 0
	}
	return int(f)
}

// MetricSummary is min/max/median/mean/n for a numeric series.
func MetricSummary(values []float64) map[string]any {
	if len(values) == 0 {
		return nil
	}
	nums := append([]float64(nil), values...)
	sort.Float64s(nums)
	n := len(nums)
	var median float64
	mid := n / 2
	if n%2 == 1 {
		median = nums[mid]
	} else {
		median = (nums[mid-1] + nums[mid]) / 2
	}
	sum := 0.0
	for _, v := range nums {
		sum += v
	}
	return map[string]any{
		"min":    CompactNum(nums[0]),
		"max":    CompactNum(nums[n-1]),
		"median": CompactNum(median),
		"mean":   CompactNum(sum / float64(n)),
		"n":      n,
	}
}

func catalogCategories(catalog map[string]any) map[string]string {
	out := map[string]string{}
	if catalog == nil {
		return out
	}
	for _, c := range catalogCards(catalog) {
		name, _ := c["name"].(string)
		category, _ := c["category"].(string)
		if name != "" && category != "" {
			out[name] = category
		}
	}
	return out
}

func summaries(valuesByKey map[string][]float64, minN int) map[string]any {
	out := map[string]any{}
	for key, vals := range valuesByKey {
		summary := MetricSummary(vals)
		if summary != nil {
			n, _ := summary["n"].(int)
			if n >= minN {
				out[key] = summary
			}
		}
	}
	return out
}

func overviewStats(cards map[string]map[string]any, minN int) map[string]any {
	evalValues := map[string][]float64{}
	for _, key := range EvaluationKeys {
		evalValues[key] = nil
	}
	priceValues := map[string][]float64{}
	for _, key := range PricingKeys {
		priceValues[key] = nil
	}
	priceValues["cost_per_task"] = nil
	perfValues := map[string][]float64{}
	for _, key := range PerformanceKeys {
		perfValues[key] = nil
	}
	for _, card := range cards {
		overview, _ := card["overview"].(map[string]any)
		if overview == nil {
			overview = map[string]any{}
		}
		if evaluations, ok := overview["evaluations"].(map[string]any); ok {
			for _, key := range EvaluationKeys {
				if f, ok := asFloat(evaluations[key]); ok {
					evalValues[key] = append(evalValues[key], f)
				}
			}
		}
		if pricing, ok := overview["pricing"].(map[string]any); ok {
			for _, key := range PricingKeys {
				if f, ok := asFloat(pricing[key]); ok {
					priceValues[key] = append(priceValues[key], f)
				}
			}
		}
		if f, ok := asFloat(card["cost_per_task"]); ok {
			priceValues["cost_per_task"] = append(priceValues["cost_per_task"], f)
		}
		if performance, ok := overview["performance"].(map[string]any); ok {
			for _, key := range PerformanceKeys {
				if f, ok := asFloat(performance[key]); ok {
					perfValues[key] = append(perfValues[key], f)
				}
			}
		}
	}
	return map[string]any{
		"evaluations": summaries(evalValues, minN),
		"pricing":     summaries(priceValues, minN),
		"performance": summaries(perfValues, minN),
	}
}

func statsByCategory(cards map[string]map[string]any) map[string]any {
	grouped := map[string]map[string]map[string]any{}
	for name, c := range cards {
		category, _ := c["category"].(string)
		if category == "" {
			continue
		}
		if grouped[category] == nil {
			grouped[category] = map[string]map[string]any{}
		}
		grouped[category][name] = c
	}
	out := map[string]any{}
	for category, group := range grouped {
		stats := overviewStats(group, CandleMinN)
		evals, _ := stats["evaluations"].(map[string]any)
		pricing, _ := stats["pricing"].(map[string]any)
		perf, _ := stats["performance"].(map[string]any)
		if len(evals) > 0 || len(pricing) > 0 || len(perf) > 0 {
			out[category] = stats
		}
	}
	return out
}

func cardCompareEntry(row map[string]any, category string) map[string]any {
	slug := strOrEmpty(row["aa_slug"])
	urlStr := ""
	if slug != "" {
		urlStr = ModelPageURL(slug)
	}
	return map[string]any{
		"aa_intelligence": row["aa_intelligence"],
		"cost_per_task":   row["aa_cost_per_task"],
		"creator":         strOrEmpty(row["aa_creator"]),
		"effort":          strOrEmpty(row["effort"]),
		"category":        category,
		"url":             urlStr,
		"alternative":     row["alternative"],
		"similar":         row["similar"],
		"overview":        row["overview"],
	}
}

func rowMaps(review map[string]any) []map[string]any {
	if review == nil {
		return nil
	}
	return asMaps(review["rows"])
}

// CompareData builds site/data/aa_compare.json from a review + catalog.
func CompareData(review, catalog map[string]any) map[string]any {
	categories := catalogCategories(catalog)
	cards := map[string]map[string]any{}
	for _, row := range rowMaps(review) {
		name, _ := row["name"].(string)
		if name == "" {
			continue
		}
		category := categories[name]
		if category == "" {
			if raw, ok := row["category"].(string); ok {
				category = raw
			}
		}
		cards[name] = cardCompareEntry(row, category)
	}
	return map[string]any{
		"generated_at":      review["generated_at"],
		"cards":             cards,
		"stats":             overviewStats(cards, 1),
		"stats_by_category": statsByCategory(cards),
		"candle_min_n":      CandleMinN,
	}
}

// WriteJSON writes payload atomically as indented JSON.
func WriteJSON(path string, payload any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, path)
}

// RunRetrieve is the CLI retrieve-aa command.
func RunRetrieve(repoRoot, apiKey string, stdout io.Writer) error {
	catalogPath := filepath.Join(repoRoot, "roster", "catalog.json")
	reviewPath := filepath.Join(repoRoot, "roster", "retrieval", "aa-intelligence-review.json")
	if strings.TrimSpace(apiKey) == "" {
		fmt.Fprintln(os.Stderr, "AA_API_KEY unset; skip AA retrieve")
		return nil
	}
	raw, err := os.ReadFile(catalogPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "catalog missing: %s\n", catalogPath)
		fmt.Fprintln(os.Stderr, "Run: make -C site generate-roster")
		return derrors.Wrap(err, derrors.CodeNotFound, "retrieve-aa", "catalog missing")
	}
	var catalog map[string]any
	if err := json.Unmarshal(raw, &catalog); err != nil {
		return derrors.Wrap(err, derrors.CodeFailed, "retrieve-aa", "parse catalog")
	}
	payload := Retrieve(catalog, apiKey, nil, time.Now().UTC(), nil)
	status, _ := payload["status"].(string)
	if status == "error" {
		msg, _ := payload["error"].(string)
		if msg == "" {
			msg = "AA retrieve failed"
		}
		fmt.Fprintln(os.Stderr, msg)
		return derrors.New(derrors.CodeUnavailable, "retrieve-aa", msg)
	}
	if status == "skipped_version" {
		fmt.Fprintf(os.Stderr, "AA intelligence_index_version %v is not %d\n",
			payload["intelligence_index_version"], SupportedIndexMajor)
		_ = WriteJSON(reviewPath, payload)
		return derrors.New(derrors.CodeFailed, "retrieve-aa", "unsupported intelligence index version")
	}
	if err := WriteJSON(reviewPath, payload); err != nil {
		return err
	}
	comparePath := filepath.Join(repoRoot, "site", "data", "aa_compare.json")
	if err := WriteJSON(comparePath, CompareData(payload, catalog)); err != nil {
		return err
	}
	rows := rowMaps(payload)
	fmt.Fprintf(stdout, "wrote %s (%d rows)\n", reviewPath, len(rows))
	return nil
}
