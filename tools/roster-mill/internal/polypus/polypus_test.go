package polypus

import "testing"

func catalog(name string, hosts ...map[string]any) map[string]any {
	anyHosts := make([]any, len(hosts))
	for i, h := range hosts {
		anyHosts[i] = h
	}
	return map[string]any{"cards": []any{
		map[string]any{"name": name, "hosts": anyHosts},
	}}
}

func TestModelsURL(t *testing.T) {
	if got := ModelsURL("http://127.0.0.1:1320", "openai"); got != "http://127.0.0.1:1320/v1/apis/openai/models" {
		t.Fatal(got)
	}
	full := "https://example.ts.net/v1/apis/openai/models"
	if ModelsURL(full, "openai") != full {
		t.Fatal(ModelsURL(full, "openai"))
	}
	if got := ModelsURL("http://127.0.0.1:1320/v1", "openai"); got != "http://127.0.0.1:1320/v1/apis/openai/models" {
		t.Fatal(got)
	}
	if got := ModelsURL("http://127.0.0.1:1320", "systemone"); got != "http://127.0.0.1:1320/v1/apis/systemone/models" {
		t.Fatal(got)
	}
	urls := DefaultModelsURLs("http://127.0.0.1:1320")
	if len(urls) != 2 || urls[1] != "http://127.0.0.1:1320/v1/apis/systemone/models" {
		t.Fatalf("%v", urls)
	}
}

func TestCFLocalMatchesAtCFID(t *testing.T) {
	cat := catalog("gemma-4-26b-a4b",
		map[string]any{"id": "cloudflare", "runtime_model_id": "@cf/google/gemma-4-26b-a4b-it"},
		map[string]any{"id": "openrouter", "runtime_model_id": "google/gemma-4-26b-a4b-it"},
	)
	report := ClassifyModels([]string{"cf_local/@cf/google/gemma-4-26b-a4b-it"}, cat)
	counts := report["counts"].(map[string]any)
	if counts["matched"] != 1 || counts["proposed"] != 0 {
		t.Fatalf("%v", counts)
	}
	matched := report["matched"].([]any)[0].(map[string]any)
	cards := matched["cards"].([]any)
	if len(cards) != 1 {
		t.Fatalf("%v", cards)
	}
	if cards[0].(map[string]any)["host_id"] != "cloudflare" {
		t.Fatalf("%v", cards[0])
	}
}

func TestDeepinfraPrefixAndMaasStem(t *testing.T) {
	cat := catalog("gemma-4-26b-a4b",
		map[string]any{"id": "deepinfra", "runtime_model_id": "google/gemma-4-26b-a4b-it"},
	)
	report := ClassifyModels([]string{"deepinfra/google/gemma-4-26b-a4b-it-maas"}, cat)
	if report["counts"].(map[string]any)["matched"] != 1 {
		t.Fatalf("%v", report["counts"])
	}
}

func TestSkipNeedles(t *testing.T) {
	cat := catalog("x")
	report := ClassifyModels([]string{"openrouter/kokoro-82m", "router/foo"}, cat)
	if report["counts"].(map[string]any)["skipped"] != 2 {
		t.Fatalf("%v", report["counts"])
	}
}

func TestProposedUnknown(t *testing.T) {
	cat := catalog("gemma-4-26b-a4b",
		map[string]any{"id": "cloudflare", "runtime_model_id": "@cf/google/gemma-4-26b-a4b-it"},
	)
	report := ClassifyModels([]string{"openrouter/some/new-model"}, cat)
	if report["counts"].(map[string]any)["proposed"] != 1 {
		t.Fatalf("%v", report)
	}
}
