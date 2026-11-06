package apply

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xynova/roster-mill/internal/card"
)

var now = time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

const gemmaCard = `---
name: gemma-4-26b-a4b
title: Gemma 4 26B A4B
kind: sovereign
category: multimodal
jobs:
  - chat
tags:
  - llm
status: watch
lab: google
hosts:
  - id: cloudflare
    runtime: cloudflare
    runtime_model_id: "@cf/google/gemma"
    price_in: 0.10
    price_out: 0.30
    capacity: hobby
    limits_note: old note
  - id: deepinfra
    runtime: deepinfra
    runtime_model_id: google/gemma
    price_in: 0.07
    price_out: 0.34
    capacity: hobby
refreshed: 2026-01-01
---

# Gemma
`

const aaCard = `---
name: claude-haiku-5.5
title: Claude Haiku 5.5
kind: provider
category: budget
jobs:
  - cheap
tags:
  - llm
status: watch
lab: anthropic
hosts:
  - id: openrouter
    runtime: openrouter
    capacity: hobby
aa_intelligence: unknown
links:
  models_dev: https://models.dev/
source: test
refreshed: 2026-11-04
---

# Claude Haiku 5.5

## My experience

Keep me.
`

func TestUnconfirmedKeepsCapacity(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "multimodal", "gemma.md")
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	_ = os.WriteFile(path, []byte(gemmaCard), 0o644)
	review := map[string]any{"rows": []any{
		map[string]any{
			"status": "ok", "name": "gemma-4-26b-a4b", "host_id": "cloudflare",
			"evidence": "10k Neurons/day free tier", "price_in": nil, "price_out": nil,
		},
	}}
	code, errs := ApplyMine(review, nil, root, now)
	if code != 0 {
		t.Fatalf("%v", errs)
	}
	fm := card.Split(read(t, path)).AsMap()
	cf := host(fm, "cloudflare")
	if cf["capacity"] != "hobby" {
		t.Fatalf("capacity=%v", cf["capacity"])
	}
	if cf["limits_note"] != "10k Neurons/day free tier" {
		t.Fatalf("note=%v", cf["limits_note"])
	}
}

func TestConfirmedSetsCapacity(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "gemma.md")
	_ = os.WriteFile(path, []byte(gemmaCard), 0o644)
	review := map[string]any{"rows": []any{
		map[string]any{
			"status": "ok", "name": "gemma-4-26b-a4b", "host_id": "deepinfra",
			"evidence": "production throughput", "price_in": 0.08, "price_out": 0.35,
			"capacity_applied": false,
		},
	}}
	decisions := []Decision{{Name: "gemma-4-26b-a4b", HostID: "deepinfra", Capacity: "industrial"}}
	code, _ := ApplyMine(review, decisions, root, now)
	if code != 0 {
		t.Fatal(code)
	}
	fm := card.Split(read(t, path)).AsMap()
	di := host(fm, "deepinfra")
	if di["capacity"] != "industrial" {
		t.Fatalf("capacity=%v", di["capacity"])
	}
	row := review["rows"].([]any)[0].(map[string]any)
	if row["capacity_applied"] != true {
		t.Fatalf("applied=%v", row["capacity_applied"])
	}
}

func TestNullPricesPreserved(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "gemma.md")
	_ = os.WriteFile(path, []byte(gemmaCard), 0o644)
	review := map[string]any{"rows": []any{
		map[string]any{
			"status": "ok", "name": "gemma-4-26b-a4b", "host_id": "cloudflare",
			"evidence": "neurons", "price_in": nil, "price_out": nil,
		},
	}}
	ApplyMine(review, nil, root, now)
	fm := card.Split(read(t, path)).AsMap()
	cf := host(fm, "cloudflare")
	got, _ := cf["price_in"].(float64)
	if got != 0.10 && got != 0.1 {
		t.Fatalf("price_in=%v", cf["price_in"])
	}
}

func TestDumpKeepsKind(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "gemma.md")
	_ = os.WriteFile(path, []byte(gemmaCard), 0o644)
	review := map[string]any{"rows": []any{
		map[string]any{
			"status": "ok", "name": "gemma-4-26b-a4b", "host_id": "cloudflare", "evidence": "rpm",
		},
	}}
	ApplyMine(review, nil, root, now)
	text := read(t, path)
	if !strings.Contains(text, "kind: sovereign") {
		t.Fatal(text)
	}
	if strings.Contains(text, "roster_kind") {
		t.Fatal(text)
	}
}

func TestLoadDecisionsRejectsHosted(t *testing.T) {
	path := filepath.Join(t.TempDir(), "decisions.yaml")
	_ = os.WriteFile(path, []byte("decisions:\n  - name: gemma-4-26b-a4b\n    host_id: deepinfra\n    capacity: hosted\n"), 0o644)
	if _, err := LoadDecisions(path); err == nil {
		t.Fatal("expected error")
	}
}

func TestApplyScoreKeepsBody(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "budget", "claude-haiku-5.5.md")
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	_ = os.WriteFile(path, []byte(aaCard), 0o644)
	aaNow := time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC)
	review := map[string]any{"status": "ok", "rows": []any{
		map[string]any{
			"name": "claude-haiku-5.5", "status": "ok",
			"aa_slug": "claude-haiku-5-5", "aa_intelligence": 36, "applied": false,
		},
	}}
	code, errs := ApplyAA(review, root, aaNow)
	if code != 0 {
		t.Fatalf("%v", errs)
	}
	if review["rows"].([]any)[0].(map[string]any)["applied"] != true {
		t.Fatal("expected applied")
	}
	f := card.Split(read(t, path))
	fm := f.AsMap()
	got := fm["aa_intelligence"]
	switch v := got.(type) {
	case int:
		if v != 36 {
			t.Fatalf("aa_intelligence=%v", v)
		}
	case int64:
		if v != 36 {
			t.Fatalf("aa_intelligence=%v", v)
		}
	case uint64:
		if v != 36 {
			t.Fatalf("aa_intelligence=%v", v)
		}
	default:
		t.Fatalf("aa_intelligence=%T %v", got, got)
	}
	if !strings.Contains(f.Body, "Keep me.") {
		t.Fatalf("body=%q", f.Body)
	}
	link, _ := fm["links"].(map[string]any)["artificial_analysis"].(string)
	if link != "https://artificialanalysis.ai/models/claude-haiku-5-5" {
		t.Fatalf("links=%v", fm["links"])
	}
}

func TestSkipsWriteWhenScoreUnchanged(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "haiku.md")
	text := strings.ReplaceAll(aaCard, "aa_intelligence: unknown", "aa_intelligence: 36")
	_ = os.WriteFile(path, []byte(text), 0o644)
	review := map[string]any{"status": "ok", "rows": []any{
		map[string]any{
			"name": "claude-haiku-5.5", "status": "ok",
			"aa_slug": "claude-haiku-5-5", "aa_intelligence": 36.0, "applied": false,
		},
	}}
	code, _ := ApplyAA(review, root, time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC))
	if code != 0 || review["rows"].([]any)[0].(map[string]any)["applied"] != false {
		t.Fatal("should skip")
	}
}

func TestClearsNoScoreWhenPreviousNumber(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "general", "deepseek-v4-flash.md")
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	text := strings.ReplaceAll(aaCard, "name: claude-haiku-5.5", "name: deepseek-v4-flash")
	text = strings.ReplaceAll(text, "aa_intelligence: unknown", "aa_intelligence: 52")
	_ = os.WriteFile(path, []byte(text), 0o644)
	review := map[string]any{"status": "ok", "rows": []any{
		map[string]any{"name": "deepseek-v4-flash", "status": "no_score", "aa_intelligence": nil, "applied": false},
	}}
	code, errs := ApplyAA(review, root, time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC))
	if code != 0 {
		t.Fatalf("%v", errs)
	}
	fm := card.Split(read(t, path)).AsMap()
	if fm["aa_intelligence"] != "unknown" {
		t.Fatalf("%v", fm["aa_intelligence"])
	}
}

func TestNoopWhenReviewSkipped(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "haiku.md")
	_ = os.WriteFile(path, []byte(aaCard), 0o644)
	code, errs := ApplyAA(map[string]any{"status": "skipped_no_key", "rows": []any{}}, root, time.Time{})
	if code != 0 {
		t.Fatalf("%v", errs)
	}
	fm := card.Split(read(t, path)).AsMap()
	if fm["aa_intelligence"] != "unknown" {
		t.Fatalf("%v", fm["aa_intelligence"])
	}
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func host(fm map[string]any, id string) map[string]any {
	for _, h := range fm["hosts"].([]any) {
		hm := h.(map[string]any)
		if hm["id"] == id {
			return hm
		}
	}
	return nil
}
