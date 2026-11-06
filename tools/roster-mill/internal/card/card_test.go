package card

import (
	"path/filepath"
	"regexp"
	"testing"
)

const fixtureCard = `---
name: test-model
title: Test Model
kind: sovereign
category: cheap
jobs:
  - chat
tags:
  - llm
status: watch
lab: meta
preferred_host: cloudflare
hosts:
  - id: cloudflare
    runtime: cloudflare
    runtime_model_id: "@cf/test/model"
    price_in: 0.1
    price_out: 0.2
    capacity: hobby
source: test
refreshed: 2026-01-01
---

# Test Model

Body.
`

func TestRemapKindOnlyTopLevel(t *testing.T) {
	f := Split(fixtureCard)
	if f.Front == nil {
		t.Fatal("expected frontmatter")
	}
	dumped, err := DumpHugo(f.Front)
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`roster_kind:\s*sovereign`).MatchString(dumped) {
		t.Fatalf("missing roster_kind: %s", dumped)
	}
	if regexp.MustCompile(`\nkind:`).MatchString(dumped) {
		t.Fatalf("kind survived: %s", dumped)
	}
}

func TestHostsSurviveGeneration(t *testing.T) {
	f := Split(fixtureCard)
	fm := f.AsMap()
	if err := ValidateCard(fm, "cheap/test-model.md"); err != nil {
		t.Fatal(err)
	}
	hugo, err := DumpHugo(f.Front)
	if err != nil {
		t.Fatal(err)
	}
	text := "---\n" + hugo + "\n---\n\n" + f.Body
	out := Split(text)
	outFM := out.AsMap()
	hosts, _ := outFM["hosts"].([]any)
	if len(hosts) != 1 {
		t.Fatalf("hosts=%d", len(hosts))
	}
	h, _ := hosts[0].(map[string]any)
	if str(h["id"]) != "cloudflare" {
		t.Fatalf("id=%v", h["id"])
	}
}

func TestDuplicateHostIDFails(t *testing.T) {
	bad := `---
name: test-model
title: Test Model
kind: sovereign
category: cheap
jobs:
  - chat
tags:
  - llm
status: watch
lab: meta
hosts:
  - id: cloudflare
    runtime: cloudflare
    capacity: hobby
  - id: cloudflare
    runtime: cloudflare
    capacity: hobby
source: test
refreshed: 2026-01-01
---
`
	fm := Split(bad).AsMap()
	if err := ValidateCard(fm, "x.md"); err == nil {
		t.Fatal("expected duplicate host error")
	}
}

func TestObjectiveInSkipList(t *testing.T) {
	if _, ok := SkipNames["OBJECTIVE.md"]; !ok {
		t.Fatal("OBJECTIVE.md missing from SkipNames")
	}
}

func TestHugoFrontmatterQuotesAAIntelligence(t *testing.T) {
	f := Split(fixtureCard)
	m := Mapping(f.Front)
	Set(m, "aa_intelligence", floatNodeForTest(43.6))
	dumped, err := DumpHugo(f.Front)
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`aa_intelligence: ['"]43\.6['"]`).MatchString(dumped) {
		t.Fatalf("43.6 not quoted: %s", dumped)
	}
	Set(m, "aa_intelligence", intNodeForTest(52))
	dumped, err = DumpHugo(f.Front)
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`aa_intelligence: ['"]52['"]`).MatchString(dumped) {
		t.Fatalf("52 not quoted: %s", dumped)
	}
}

func TestAllowedCategories(t *testing.T) {
	want := []string{"cheap", "mid", "spend", "unpriced", "classify", "tts"}
	if len(AllowedCategories) != len(want) {
		t.Fatalf("len=%d", len(AllowedCategories))
	}
	for _, c := range want {
		if _, ok := AllowedCategories[c]; !ok {
			t.Fatalf("missing %s", c)
		}
	}
	fm := Split(fixtureCard).AsMap()
	fm["category"] = "classify"
	if err := ValidateCard(fm, filepath.Join("classify", "x.md")); err != nil {
		t.Fatal(err)
	}
	fm["category"] = "budget"
	if err := ValidateCard(fm, "budget/x.md"); err == nil {
		t.Fatal("expected invalid category")
	}
}

func TestAssignCategory(t *testing.T) {
	base := func(extra map[string]any) map[string]any {
		fm := map[string]any{
			"kind": "sovereign",
			"jobs": []any{"chat"},
			"hosts": []any{
				map[string]any{"id": "x", "price_in": 1.0, "price_out": 0.55},
			},
		}
		for k, v := range extra {
			fm[k] = v
		}
		return fm
	}
	if got := AssignCategory(base(map[string]any{"jobs": []any{"classify", "cheap"}})); got != "classify" {
		t.Fatalf("classify: %s", got)
	}
	if got := AssignCategory(base(map[string]any{"jobs": []any{"tts", "cheap"}})); got != "tts" {
		t.Fatalf("tts: %s", got)
	}
	if got := AssignCategory(base(map[string]any{"kind": "service", "jobs": []any{"coding"}})); got != "unpriced" {
		t.Fatalf("service: %s", got)
	}
	if got := AssignCategory(base(map[string]any{
		"hosts": []any{map[string]any{"price_in": 1.0, "price_out": 0.55}},
	})); got != "cheap" {
		t.Fatalf("cheap: %s", got)
	}
	if got := AssignCategory(base(map[string]any{
		"hosts": []any{map[string]any{"price_in": 10.0, "price_out": 4.0}},
	})); got != "mid" {
		t.Fatalf("mid: %s", got)
	}
	if got := AssignCategory(base(map[string]any{
		"hosts": []any{map[string]any{"price_in": 10.0, "price_out": 10.0}},
	})); got != "spend" {
		t.Fatalf("spend: %s", got)
	}
	if got := AssignCategory(base(map[string]any{
		"hosts": []any{map[string]any{"price_note": "on request"}},
	})); got != "unpriced" {
		t.Fatalf("unpriced: %s", got)
	}
}

func TestCardToJSON(t *testing.T) {
	fm := Split(fixtureCard).AsMap()
	j := ToJSON(fm)
	if j["name"] != "test-model" {
		t.Fatalf("name=%v", j["name"])
	}
	hosts, _ := j["hosts"].([]any)
	h, _ := hosts[0].(map[string]any)
	if h["capacity"] != "hobby" {
		t.Fatalf("capacity=%v", h["capacity"])
	}
	if j["preferred_host"] != "cloudflare" {
		t.Fatalf("preferred_host=%v", j["preferred_host"])
	}
}

func floatNodeForTest(f float64) *yamlScalar {
	return yamlFloat(f)
}

func intNodeForTest(n int) *yamlScalar {
	return yamlInt(n)
}
