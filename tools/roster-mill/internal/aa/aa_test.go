package aa

import (
	"testing"
	"time"
)

func aaModel(name, slug string, score any, opts map[string]any) map[string]any {
	evals := map[string]any{"artificial_analysis_intelligence_index": score}
	if extra, ok := opts["evaluations"].(map[string]any); ok {
		for k, v := range extra {
			evals[k] = v
		}
	}
	payload := map[string]any{
		"name":          name,
		"slug":          slug,
		"model_creator": map[string]any{"name": "Lab"},
		"evaluations":   evals,
		"artificial_analysis_intelligence_index_cost": nil,
		"pricing":     map[string]any{},
		"performance": map[string]any{},
	}
	if creator, ok := opts["creator"].(string); ok {
		payload["model_creator"] = map[string]any{"name": creator}
	}
	if cost, ok := opts["cost"]; ok && cost != nil {
		payload["artificial_analysis_intelligence_index_cost"] = map[string]any{
			"cost_per_task": map[string]any{"total_cost": cost},
		}
	}
	if p, ok := opts["pricing"].(map[string]any); ok {
		payload["pricing"] = p
	}
	if p, ok := opts["performance"].(map[string]any); ok {
		payload["performance"] = p
	}
	return payload
}

func catalogCard(name, title, kind string) map[string]any {
	return map[string]any{"name": name, "title": title, "kind": kind}
}

func TestPrefersMaxOverHigh(t *testing.T) {
	catalog := map[string]any{"cards": []any{
		catalogCard("claude-haiku-5.5", "Claude Haiku 5.5", "provider"),
	}}
	models := []map[string]any{
		aaModel("Claude Haiku 5.5 (high)", "claude-haiku-5-5", 30, nil),
		aaModel("Claude Haiku 5.5 (max)", "claude-haiku-5-5", 36, nil),
	}
	rows := MatchCards(catalog, models)
	if rows[0]["status"] != "ok" || rows[0]["aa_intelligence"] != 36 || rows[0]["effort"] != "max" {
		t.Fatalf("%v", rows[0])
	}
}

func TestSkipsNonReasoningWhenThinkingExists(t *testing.T) {
	catalog := map[string]any{"cards": []any{
		catalogCard("gemini-3.8-flash", "Gemini 3.8 Flash", "provider"),
	}}
	models := []map[string]any{
		aaModel("Gemini 3.8 Flash (Non-reasoning)", "gemini-3-8-flash", 12, nil),
		aaModel("Gemini 3.8 Flash (High)", "gemini-3-8-flash", 31, nil),
	}
	rows := MatchCards(catalog, models)
	if rows[0]["aa_intelligence"] != 31 || rows[0]["effort"] != "high" {
		t.Fatalf("%v", rows[0])
	}
}

func TestUnmatchedWhenOnlyNonReasoning(t *testing.T) {
	catalog := map[string]any{"cards": []any{
		catalogCard("gemini-2.5-flash", "Gemini 2.5 Flash", "provider"),
	}}
	models := []map[string]any{
		aaModel("Gemini 2.5 Flash (Non-reasoning)", "gemini-2-5-flash", 20, nil),
	}
	rows := MatchCards(catalog, models)
	if rows[0]["status"] != "unmatched" || rows[0]["aa_intelligence"] != nil {
		t.Fatalf("%v", rows[0])
	}
}

func TestGLMFlashIgnoresGLMMax(t *testing.T) {
	catalog := map[string]any{"cards": []any{
		catalogCard("glm-5.3-flash", "GLM-5.3-Flash", "sovereign"),
	}}
	models := []map[string]any{aaModel("GLM-5.3 (Max)", "glm-5-3", 32, nil)}
	rows := MatchCards(catalog, models)
	if rows[0]["status"] != "unmatched" {
		t.Fatalf("%v", rows[0])
	}
}

func TestCursorCardSkipped(t *testing.T) {
	catalog := map[string]any{"cards": []any{
		map[string]any{
			"name": "cursor-gemini-3-flash", "title": "Gemini 3 Flash",
			"kind": "service", "service": "cursor",
		},
	}}
	models := []map[string]any{aaModel("Gemini 3 Flash (high)", "gemini-3-flash", 20, nil)}
	rows := MatchCards(catalog, models)
	if rows[0]["status"] != "skipped_service" {
		t.Fatalf("%v", rows[0])
	}
}

func TestListedNullStaysNoScore(t *testing.T) {
	catalog := map[string]any{"cards": []any{
		map[string]any{
			"name": "gemma-4-26b-a4b", "title": "Gemma 4 26B A4B",
			"kind": "sovereign", "aa_intelligence": 31,
		},
	}}
	models := []map[string]any{aaModel("Gemma 4 26B A4B", "gemma-4-26b-a4b", nil, nil)}
	rows := MatchCards(catalog, models)
	if rows[0]["status"] != "no_score" || rows[0]["aa_intelligence"] != nil {
		t.Fatalf("%v", rows[0])
	}
}

func TestSkipWithoutKey(t *testing.T) {
	out := Retrieve(map[string]any{"cards": []any{}}, "", nil, time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC), nil)
	if out["status"] != "skipped_no_key" {
		t.Fatalf("%v", out["status"])
	}
}

func TestSkipOtherIndexMajor(t *testing.T) {
	fetch := func(page int, key string) (map[string]any, error) {
		return map[string]any{
			"intelligence_index_version": 6,
			"pagination":                 map[string]any{"has_more": false},
			"data":                       []any{aaModel("Claude Haiku 5.5 (max)", "claude-haiku-5-5", 36, nil)},
		}, nil
	}
	out := Retrieve(map[string]any{"cards": []any{}}, "test-key", fetch, time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC), func(time.Duration) {})
	if out["status"] != "skipped_version" {
		t.Fatalf("%v", out["status"])
	}
}

func TestPaginatesAndMatches(t *testing.T) {
	pages := map[int]map[string]any{
		1: {
			"intelligence_index_version": 5,
			"pagination":                 map[string]any{"has_more": true},
			"data":                       []any{aaModel("Granite 4.2 8B", "granite-4-2-8b", 5, nil)},
		},
		2: {
			"intelligence_index_version": 5,
			"pagination":                 map[string]any{"has_more": false},
			"data":                       []any{aaModel("Gemini 3.8 Flash (high)", "gemini-3-8-flash", 31, nil)},
		},
	}
	fetch := func(page int, key string) (map[string]any, error) {
		if key != "test-key" {
			t.Fatalf("key=%s", key)
		}
		return pages[page], nil
	}
	catalog := map[string]any{"cards": []any{
		catalogCard("granite-4.2-8b", "Granite 4.2 8B", "sovereign"),
		catalogCard("gemini-3.8-flash", "Gemini 3.8 Flash", "provider"),
	}}
	out := Retrieve(catalog, "test-key", fetch, time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC), func(time.Duration) {})
	if out["status"] != "ok" {
		t.Fatalf("%v", out["status"])
	}
	byName := map[string]map[string]any{}
	for _, item := range out["rows"].([]any) {
		row := item.(map[string]any)
		byName[row["name"].(string)] = row
	}
	if byName["granite-4.2-8b"]["aa_intelligence"] != 5 {
		t.Fatalf("%v", byName["granite-4.2-8b"])
	}
	if byName["gemini-3.8-flash"]["aa_intelligence"] != 31 || byName["gemini-3.8-flash"]["effort"] != "high" {
		t.Fatalf("%v", byName["gemini-3.8-flash"])
	}
}

func TestAlternativeBeatsIntelligencePerCost(t *testing.T) {
	catalog := map[string]any{"cards": []any{
		catalogCard("claude-haiku-5.5", "Claude Haiku 5.5", "provider"),
	}}
	models := []map[string]any{
		aaModel("Claude Haiku 5.5 (Max)", "claude-haiku-5-5", 36.5, map[string]any{"creator": "Anthropic", "cost": 0.41}),
		aaModel("Better Peer (Max)", "better-peer", 40, map[string]any{"creator": "DeepSeek", "cost": 0.20}),
		aaModel("Smarter Spend (Max)", "smarter-spend", 50, map[string]any{"creator": "Google", "cost": 0.50}),
		aaModel("Cheaper Dumber (Max)", "cheaper-dumber", 10, map[string]any{"creator": "IBM", "cost": 0.01}),
	}
	rows := MatchCards(catalog, models)
	alt := rows[0]["alternative"].(map[string]any)
	if alt["name"] != "Better Peer (Max)" || alt["creator"] != "DeepSeek" || alt["aa_intelligence"] != 40 {
		t.Fatalf("%v", alt)
	}
	if alt["cost_per_task"] != 0.2 {
		t.Fatalf("cost=%v", alt["cost_per_task"])
	}
	payload := CompareData(map[string]any{"generated_at": "2026-11-05", "rows": rows}, nil)
	entry := payload["cards"].(map[string]map[string]any)["claude-haiku-5.5"]
	if entry["aa_intelligence"] != 36.5 {
		t.Fatalf("%v", entry)
	}
}

func TestAlternativeSkipsSmarterButCostlier(t *testing.T) {
	catalog := map[string]any{"cards": []any{
		catalogCard("claude-haiku-5.5", "Claude Haiku 5.5", "provider"),
	}}
	models := []map[string]any{
		aaModel("Claude Haiku 5.5 (Max)", "claude-haiku-5-5", 36.5, map[string]any{"creator": "Anthropic", "cost": 0.41}),
		aaModel("Smarter Spend (Max)", "smarter-spend", 50, map[string]any{"creator": "Google", "cost": 0.50}),
	}
	rows := MatchCards(catalog, models)
	if rows[0]["alternative"] != nil {
		t.Fatalf("%v", rows[0]["alternative"])
	}
}

func TestSimilarIsClosestOffRoster(t *testing.T) {
	catalog := map[string]any{"cards": []any{
		catalogCard("claude-haiku-5.5", "Claude Haiku 5.5", "provider"),
	}}
	models := []map[string]any{
		aaModel("Claude Haiku 5.5 (Max)", "claude-haiku-5-5", 36.5, map[string]any{"creator": "Anthropic", "cost": 0.41}),
		aaModel("Better Peer (Max)", "better-peer", 40, map[string]any{"creator": "DeepSeek", "cost": 0.20}),
		aaModel("Near Peer (High)", "near-peer", 38, map[string]any{"creator": "Z AI", "cost": 1.0}),
	}
	rows := MatchCards(catalog, models)
	if rows[0]["alternative"].(map[string]any)["slug"] != "better-peer" {
		t.Fatalf("alt=%v", rows[0]["alternative"])
	}
	similar := rows[0]["similar"].(map[string]any)
	if similar["name"] != "Near Peer (High)" || similar["creator"] != "Z AI" {
		t.Fatalf("%v", similar)
	}
}

func TestSkipsOtherRosterCards(t *testing.T) {
	catalog := map[string]any{"cards": []any{
		catalogCard("muse-wonder", "Muse Wonder", "provider"),
		catalogCard("claude-haiku-5.5", "Claude Haiku 5.5", "provider"),
	}}
	models := []map[string]any{
		aaModel("Muse Wonder (Max)", "muse-wonder", 43.6, map[string]any{"creator": "Meta", "cost": 6.95}),
		aaModel("Claude Haiku 5.5 (Max)", "claude-haiku-5-5", 36.5, map[string]any{"creator": "Anthropic", "cost": 0.41}),
	}
	rows := MatchCards(catalog, models)
	byName := map[string]map[string]any{}
	for _, row := range rows {
		byName[row["name"].(string)] = row
	}
	if byName["muse-wonder"]["alternative"] != nil || byName["muse-wonder"]["similar"] != nil {
		t.Fatalf("%v", byName["muse-wonder"])
	}
}

func TestCompareKeepsFreeAPIOverview(t *testing.T) {
	catalog := map[string]any{"cards": []any{
		catalogCard("muse-wonder", "Muse Wonder", "provider"),
	}}
	models := []map[string]any{
		aaModel("Muse Wonder (Max)", "muse-wonder", 43.6, map[string]any{
			"creator": "Meta", "cost": 6.95,
			"evaluations": map[string]any{
				"artificial_analysis_coding_index":      40.1,
				"artificial_analysis_engineering_index": 41.2,
			},
			"pricing": map[string]any{
				"price_1m_input_tokens": 4, "price_1m_output_tokens": 20,
			},
			"performance": map[string]any{
				"median_output_tokens_per_second":    80.07,
				"median_time_to_first_token_seconds": 2.22,
			},
		}),
	}
	rows := MatchCards(catalog, models)
	overview := rows[0]["overview"].(map[string]any)
	evals := overview["evaluations"].(map[string]any)
	if evals["artificial_analysis_intelligence_index"] != 43.6 {
		t.Fatalf("%v", evals)
	}
	payload := CompareData(map[string]any{"generated_at": "2026-11-05", "rows": rows}, nil)
	intel := payload["stats"].(map[string]any)["evaluations"].(map[string]any)["artificial_analysis_intelligence_index"].(map[string]any)
	if intel["min"] != 43.6 || intel["n"] != 1 {
		t.Fatalf("%v", intel)
	}
	if payload["candle_min_n"] != 3 {
		t.Fatalf("candle_min_n=%v", payload["candle_min_n"])
	}
	if len(payload["stats_by_category"].(map[string]any)) != 0 {
		t.Fatalf("%v", payload["stats_by_category"])
	}
}

func TestMetricSummary(t *testing.T) {
	summary := MetricSummary([]float64{4.4, 43.6, 27.8, 36.5, 5.1})
	if summary["min"] != 4.4 || summary["max"] != 43.6 || summary["median"] != 27.8 || summary["mean"] != 23.48 || summary["n"] != 5 {
		t.Fatalf("%v", summary)
	}
}

func intelRow(name string, score float64) map[string]any {
	evaluations := map[string]any{}
	for _, key := range EvaluationKeys {
		evaluations[key] = nil
	}
	evaluations["artificial_analysis_intelligence_index"] = score
	pricing := map[string]any{}
	for _, key := range PricingKeys {
		pricing[key] = nil
	}
	performance := map[string]any{}
	for _, key := range PerformanceKeys {
		performance[key] = nil
	}
	return map[string]any{
		"name": name, "aa_slug": name, "aa_intelligence": score,
		"overview": map[string]any{
			"evaluations": evaluations, "pricing": pricing, "performance": performance,
		},
	}
}

func TestCompareDataAcceptsTypedCardSlice(t *testing.T) {
	catalog := map[string]any{
		"cards": []map[string]any{
			{"name": "cheap-a", "category": "cheap"},
			{"name": "cheap-b", "category": "cheap"},
			{"name": "cheap-c", "category": "cheap"},
		},
	}
	rows := []map[string]any{
		intelRow("cheap-a", 10), intelRow("cheap-b", 20), intelRow("cheap-c", 30),
	}
	payload := CompareData(map[string]any{"generated_at": "2026-11-06", "rows": rows}, catalog)
	byCat := payload["stats_by_category"].(map[string]any)
	if _, ok := byCat["cheap"]; !ok {
		t.Fatalf("typed []map catalog dropped categories: %v", payload["stats_by_category"])
	}
}

func TestCheapStatsExcludeMidAndNeedThreePeers(t *testing.T) {
	catalog := map[string]any{"cards": []any{
		map[string]any{"name": "cheap-a", "category": "cheap"},
		map[string]any{"name": "cheap-b", "category": "cheap"},
		map[string]any{"name": "cheap-c", "category": "cheap"},
		map[string]any{"name": "mid-a", "category": "mid"},
	}}
	rows := []map[string]any{
		intelRow("cheap-a", 10), intelRow("cheap-b", 20),
		intelRow("cheap-c", 30), intelRow("mid-a", 80),
	}
	payload := CompareData(map[string]any{"generated_at": "2026-11-06", "rows": rows}, catalog)
	cards := payload["cards"].(map[string]map[string]any)
	if cards["cheap-a"]["category"] != "cheap" || cards["mid-a"]["category"] != "mid" {
		t.Fatalf("%v", cards)
	}
	cheap := payload["stats_by_category"].(map[string]any)["cheap"].(map[string]any)["evaluations"].(map[string]any)["artificial_analysis_intelligence_index"].(map[string]any)
	if cheap["min"] != 10 || cheap["max"] != 30 || cheap["n"] != 3 {
		t.Fatalf("%v", cheap)
	}
	if _, ok := payload["stats_by_category"].(map[string]any)["mid"]; ok {
		t.Fatal("mid should be omitted")
	}
	global := payload["stats"].(map[string]any)["evaluations"].(map[string]any)["artificial_analysis_intelligence_index"].(map[string]any)
	if global["max"] != 80 || global["n"] != 4 {
		t.Fatalf("%v", global)
	}
}

func TestTwoCheapCardsOmitBandCandle(t *testing.T) {
	catalog := map[string]any{"cards": []any{
		map[string]any{"name": "cheap-a", "category": "cheap"},
		map[string]any{"name": "cheap-b", "category": "cheap"},
	}}
	rows := []map[string]any{intelRow("cheap-a", 10), intelRow("cheap-b", 20)}
	payload := CompareData(map[string]any{"generated_at": "2026-11-06", "rows": rows}, catalog)
	if len(payload["stats_by_category"].(map[string]any)) != 0 {
		t.Fatalf("%v", payload["stats_by_category"])
	}
}
