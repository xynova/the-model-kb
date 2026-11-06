package limits

import (
	"testing"
	"time"
)

var now = time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)

func catalogWithHost(host map[string]any, kind, name string) map[string]any {
	if kind == "" {
		kind = "sovereign"
	}
	if name == "" {
		name = "test-card"
	}
	return map[string]any{
		"cards": []any{
			map[string]any{"name": name, "kind": kind, "hosts": []any{host}},
		},
	}
}

func TestNeuronsSuggestHobby(t *testing.T) {
	html := "<p>Workers AI includes 10,000 Neurons per day on the free tier.</p>"
	cat := catalogWithHost(map[string]any{
		"id": "cloudflare", "capacity": "unknown",
		"links": map[string]any{"cloudflare_pricing": "https://developers.cloudflare.com/workers-ai/pricing/"},
	}, "", "gemma-4-26b-a4b")
	out := Retrieve(cat, func(string) (int, string) { return 200, html }, now, func(time.Duration) {})
	row := out["rows"].([]any)[0].(map[string]any)
	if row["status"] != "ok" || row["suggested_capacity"] != "hobby" {
		t.Fatalf("%v", row)
	}
	if row["evidence"] == "" {
		t.Fatal("expected evidence")
	}
}

func TestOpenRouterMarketplaceHobby(t *testing.T) {
	html := "<p>OpenRouter routes to many providers. Pay per token.</p>"
	cat := catalogWithHost(map[string]any{
		"id": "openrouter", "capacity": "hobby",
		"links": map[string]any{"openrouter": "https://openrouter.ai/docs/pricing"},
	}, "", "")
	out := Retrieve(cat, func(string) (int, string) { return 200, html }, now, func(time.Duration) {})
	row := out["rows"].([]any)[0].(map[string]any)
	if row["suggested_capacity"] != "hobby" || row["suggestion_reason"] != "openrouter_marketplace" {
		t.Fatalf("%v", row)
	}
}

func TestGoogleNoRpmUnknown(t *testing.T) {
	html := "<p>Vertex AI pricing varies by region and model.</p>"
	cat := catalogWithHost(map[string]any{
		"id": "google", "capacity": "unknown",
		"links": map[string]any{"provider_pricing": "https://cloud.google.com/vertex-ai/generative-ai/pricing"},
	}, "", "")
	out := Retrieve(cat, func(string) (int, string) { return 200, html }, now, func(time.Duration) {})
	row := out["rows"].([]any)[0].(map[string]any)
	if row["suggested_capacity"] != "unknown" || row["suggestion_reason"] != "google_no_rpm" {
		t.Fatalf("%v", row)
	}
}

func TestDeepinfraFlexUnknown(t *testing.T) {
	html := "<p>Priority and Flex tiers for burst throughput.</p>"
	cat := catalogWithHost(map[string]any{
		"id": "deepinfra", "capacity": "unknown",
		"links": map[string]any{"deepinfra": "https://deepinfra.com/pricing"},
	}, "", "")
	out := Retrieve(cat, func(string) (int, string) { return 200, html }, now, func(time.Duration) {})
	row := out["rows"].([]any)[0].(map[string]any)
	if row["suggested_capacity"] != "unknown" || row["suggestion_reason"] != "deepinfra_no_rpm" {
		t.Fatalf("%v", row)
	}
}

func TestDeepinfraProductionIndustrial(t *testing.T) {
	html := "<p>Our production API offers high-throughput for enterprise.</p>"
	cat := catalogWithHost(map[string]any{
		"id": "deepinfra", "capacity": "unknown",
		"links": map[string]any{"deepinfra": "https://deepinfra.com/pricing"},
	}, "", "")
	out := Retrieve(cat, func(string) (int, string) { return 200, html }, now, func(time.Duration) {})
	row := out["rows"].([]any)[0].(map[string]any)
	if row["suggested_capacity"] != "industrial" || row["suggestion_reason"] != "deepinfra_production_text" {
		t.Fatalf("%v", row)
	}
}

func TestServiceAndCursorNeverFetch(t *testing.T) {
	called := false
	cat := map[string]any{"cards": []any{
		map[string]any{
			"name": "cursor-pro", "kind": "service",
			"hosts": []any{map[string]any{
				"id": "cursor", "capacity": "unknown",
				"links": map[string]any{"service_docs": "https://cursor.com/docs"},
			}},
		},
	}}
	Retrieve(cat, func(string) (int, string) { called = true; return 200, "neurons" }, now, func(time.Duration) {})
	if called {
		t.Fatal("fetch should not run")
	}
}

func TestXinferenceUnknown(t *testing.T) {
	html := "<p>Model API served from Australia. Enterprise hardens it for production. Gemma on request.</p>"
	cat := catalogWithHost(map[string]any{
		"id": "xinference", "capacity": "unknown",
		"links": map[string]any{"xinference": "https://xinference.co/pricing.html"},
	}, "", "")
	out := Retrieve(cat, func(string) (int, string) { return 200, html }, now, func(time.Duration) {})
	row := out["rows"].([]any)[0].(map[string]any)
	if row["status"] != "ok" || row["suggested_capacity"] != "unknown" || row["suggestion_reason"] != "xinference_no_rpm" {
		t.Fatalf("%v", row)
	}
}

func TestWrongHostnameSkipped(t *testing.T) {
	called := false
	cat := catalogWithHost(map[string]any{
		"id": "cloudflare", "capacity": "hobby",
		"links": map[string]any{"cloudflare_pricing": "https://example.com/pricing"},
	}, "", "")
	out := Retrieve(cat, func(string) (int, string) { called = true; return 200, "rpm" }, now, func(time.Duration) {})
	row := out["rows"].([]any)[0].(map[string]any)
	if row["status"] != "skipped" || called {
		t.Fatalf("%v called=%v", row, called)
	}
}
