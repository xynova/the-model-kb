package catalog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFormatISOWeek(t *testing.T) {
	t.Parallel()
	got := FormatISOWeek(time.Date(2026, 9, 16, 0, 0, 0, 0, time.UTC))
	if got != "2026-W38" {
		t.Fatalf("FormatISOWeek = %q, want 2026-W38", got)
	}
	if got := FormatWeekRange("2026-W38"); got != "Sep 14–20, 2026" {
		t.Fatalf("FormatWeekRange = %q", got)
	}
}

func TestRefreshIndexesAndAdditions(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "models", "auk.md"), `---
name: auk
title: AuK
category: models
does: speech
tags:
  - audio
  - open
status: watch
run: local
license: open
source: test
added: 2026-09-16
---

# AuK

## What it is

Speech model for local TTS and voice edit.
`)
	mustWrite(t, filepath.Join(root, "creative", "suno-v6.md"), `---
name: suno-v6
title: Suno v6
category: creative
does: music
tags:
  - audio
status: watch
run: api
license: closed
source: test
added: 2026-09-16
---

# Suno

## What it is

Closed music generator over an API.
`)
	if err := os.MkdirAll(filepath.Join(root, "video"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := RefreshIndexes(root); err != nil {
		t.Fatal(err)
	}
	global := mustRead(t, filepath.Join(root, "INDEX.md"))
	if !strings.Contains(global, "| auk |") || !strings.Contains(global, "Speech model") {
		t.Fatalf("global INDEX missing auk/summary:\n%s", global)
	}
	if !strings.Contains(global, "| name | title | category | summary | tags | run | license |") {
		t.Fatalf("global INDEX columns:\n%s", global)
	}
	if strings.Contains(global, "| added |") {
		t.Fatalf("global INDEX should not show added column:\n%s", global)
	}
	modelsIdx := mustRead(t, filepath.Join(root, "models", "INDEX.md"))
	if !strings.Contains(modelsIdx, "| auk |") {
		t.Fatalf("models INDEX missing auk:\n%s", modelsIdx)
	}
	videoIdx := mustRead(t, filepath.Join(root, "video", "INDEX.md"))
	if !strings.Contains(videoIdx, "| name | title |") {
		t.Fatalf("video INDEX missing header:\n%s", videoIdx)
	}

	now := time.Date(2026, 9, 17, 12, 0, 0, 0, time.Local)
	if err := RebuildAdditions(root, now); err != nil {
		t.Fatal(err)
	}
	week := mustRead(t, filepath.Join(root, "additions", "2026-W38.md"))
	if !strings.Contains(week, "status: open") {
		t.Fatalf("expected open week:\n%s", week)
	}
	if !strings.Contains(week, "weight: 1") {
		t.Fatalf("newest week should be weight 1 (sidebar top):\n%s", week)
	}
	if !strings.Contains(week, "Sep 14–20, 2026") {
		t.Fatalf("week title should show calendar range:\n%s", week)
	}
	if !strings.Contains(week, "week_start: 2026-09-14") || !strings.Contains(week, "week_end: 2026-09-20") {
		t.Fatalf("week frontmatter missing start/end:\n%s", week)
	}
	if !strings.Contains(week, "## models") || !strings.Contains(week, "## creative") {
		t.Fatalf("week missing section indexes:\n%s", week)
	}
	if strings.Contains(week, "| added |") {
		t.Fatalf("week tables should not show added column:\n%s", week)
	}
	addIdx := mustRead(t, filepath.Join(root, "additions", "INDEX.md"))
	if !strings.Contains(addIdx, "| 2026-W38 | Sep 14–20, 2026 | open |") {
		t.Fatalf("additions INDEX:\n%s", addIdx)
	}

	mustWrite(t, filepath.Join(root, "models", "old.md"), `---
name: old
title: Old
category: models
does: llm
tags:
  - llm
status: watch
run: local
license: open
source: test
added: 2026-08-31
---

# Old

## What it is

Older model kept for week ordering tests.
`)
	if err := RebuildAdditions(root, now); err != nil {
		t.Fatal(err)
	}
	past := mustRead(t, filepath.Join(root, "additions", "2026-W36.md"))
	if !strings.Contains(past, "status: closed") {
		t.Fatalf("expected closed past week:\n%s", past)
	}
	if !strings.Contains(past, "weight: 2") {
		t.Fatalf("older week should be weight 2:\n%s", past)
	}
	week = mustRead(t, filepath.Join(root, "additions", "2026-W38.md"))
	if !strings.Contains(week, "status: open") || !strings.Contains(week, "weight: 1") {
		t.Fatalf("current week should stay open with weight 1:\n%s", week)
	}
}

func mustWrite(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
