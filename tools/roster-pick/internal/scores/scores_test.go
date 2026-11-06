package scores_test

import (
	"encoding/json"
	"testing"

	"github.com/xynova/roster-pick/internal/scores"
)

func TestLeafNotApplicable(t *testing.T) {
	t.Parallel()
	const snippet = `{
  "scores_schema": 1,
  "generated_at": "2026-01-01",
  "source_catalog_schema": 1,
  "offers": [{
    "name": "jev",
    "title": "Jev",
    "kind": "provider",
    "status": "watch",
    "jobs": ["classify"],
    "host_id": "local",
    "capacity": "unknown",
    "industrial_ok": false,
    "quality": {
      "quality.general_problem_solving.benchmark_composite": {
        "status": "not_applicable"
      }
    }
  }]
}`
	var doc scores.Document
	if err := json.Unmarshal([]byte(snippet), &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(doc.Offers) != 1 {
		t.Fatalf("offers: got %d", len(doc.Offers))
	}
	leaf := doc.Offers[0].Quality["quality.general_problem_solving.benchmark_composite"]
	if leaf.Status != "not_applicable" {
		t.Fatalf("status: got %q", leaf.Status)
	}
	if leaf.Value != nil {
		t.Fatalf("value: want nil, got %v", *leaf.Value)
	}
}

func TestLeafMeasuredAA(t *testing.T) {
	t.Parallel()
	const snippet = `{
  "scores_schema": 1,
  "generated_at": "2026-01-01",
  "source_catalog_schema": 1,
  "offers": [{
    "name": "deepseek-v4-flash",
    "title": "DeepSeek V4 Flash",
    "kind": "sovereign",
    "status": "watch",
    "jobs": ["chat"],
    "host_id": "deepinfra",
    "capacity": "industrial",
    "industrial_ok": true,
    "quality": {
      "quality.general_problem_solving.benchmark_composite": {
        "status": "measured",
        "value": 52,
        "attribution": "model_only",
        "comparison_key": "aa_intelligence_index:v5.0.1"
      }
    }
  }]
}`
	var doc scores.Document
	if err := json.Unmarshal([]byte(snippet), &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	leaf := doc.Offers[0].Quality["quality.general_problem_solving.benchmark_composite"]
	if leaf.Status != "measured" {
		t.Fatalf("status: got %q", leaf.Status)
	}
	if leaf.Value == nil || *leaf.Value != 52 {
		t.Fatalf("value: got %v", leaf.Value)
	}
}

func TestLoadOmitsQuality(t *testing.T) {
	t.Parallel()
	const snippet = `{
  "scores_schema": 1,
  "generated_at": "2026-01-01",
  "source_catalog_schema": 1,
  "offers": [{
    "name": "x",
    "title": "X",
    "kind": "sovereign",
    "status": "watch",
    "jobs": [],
    "host_id": "local",
    "capacity": "hobby",
    "industrial_ok": false
  }]
}`
	var doc scores.Document
	if err := json.Unmarshal([]byte(snippet), &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if doc.Offers[0].Quality != nil {
		t.Fatalf("quality: want nil map, got %#v", doc.Offers[0].Quality)
	}
}
