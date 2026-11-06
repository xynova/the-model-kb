"""Tests for apply-aa.py."""

from __future__ import annotations

import importlib.util
import tempfile
import unittest
from datetime import datetime, timezone
from pathlib import Path

SCRIPTS = Path(__file__).resolve().parent
SPEC = importlib.util.spec_from_file_location("apply_aa", SCRIPTS / "apply-aa.py")
assert SPEC and SPEC.loader
aa = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(aa)

GR_SPEC = importlib.util.spec_from_file_location(
    "generate_roster", SCRIPTS / "generate-roster.py"
)
assert GR_SPEC and GR_SPEC.loader
gr = importlib.util.module_from_spec(GR_SPEC)
GR_SPEC.loader.exec_module(gr)

NOW = datetime(2026, 11, 5, tzinfo=timezone.utc)

CARD = """---
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
"""


class TestApplyScore(unittest.TestCase):
    def test_files_max_row_and_keeps_body(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            card = root / "budget" / "claude-haiku-5.5.md"
            card.parent.mkdir(parents=True)
            card.write_text(CARD, encoding="utf-8")
            review = {
                "status": "ok",
                "rows": [
                    {
                        "name": "claude-haiku-5.5",
                        "status": "ok",
                        "aa_slug": "claude-haiku-5-5",
                        "aa_intelligence": 36,
                        "applied": False,
                    }
                ],
            }
            code, errs = aa.apply(review, root, now=NOW)
            self.assertEqual(code, 0, errs)
            self.assertTrue(review["rows"][0]["applied"])
            text = card.read_text(encoding="utf-8")
            fm, body = gr.split_frontmatter(text)
            assert fm
            self.assertEqual(fm["aa_intelligence"], 36)
            self.assertEqual(fm["refreshed"], "2026-11-05")
            self.assertEqual(
                fm["links"]["artificial_analysis"],
                "https://artificialanalysis.ai/models/claude-haiku-5-5",
            )
            self.assertIn("Keep me.", body)

    def test_skips_write_when_score_unchanged(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            card = root / "haiku.md"
            text = CARD.replace("aa_intelligence: unknown", "aa_intelligence: 36")
            card.write_text(text, encoding="utf-8")
            review = {
                "status": "ok",
                "rows": [
                    {
                        "name": "claude-haiku-5.5",
                        "status": "ok",
                        "aa_slug": "claude-haiku-5-5",
                        "aa_intelligence": 36.0,
                        "applied": False,
                    }
                ],
            }
            code, _ = aa.apply(review, root, now=NOW)
            self.assertEqual(code, 0)
            self.assertFalse(review["rows"][0]["applied"])
            fm, _ = gr.split_frontmatter(card.read_text(encoding="utf-8"))
            assert fm
            self.assertEqual(str(fm["refreshed"]), "2026-11-04")

    def test_clears_no_score_when_previous_number(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            card = root / "general" / "deepseek-v4-flash.md"
            card.parent.mkdir(parents=True)
            text = CARD.replace("name: claude-haiku-5.5", "name: deepseek-v4-flash")
            text = text.replace("aa_intelligence: unknown", "aa_intelligence: 52")
            card.write_text(text, encoding="utf-8")
            review = {
                "status": "ok",
                "rows": [
                    {
                        "name": "deepseek-v4-flash",
                        "status": "no_score",
                        "aa_intelligence": None,
                        "applied": False,
                    }
                ],
            }
            code, errs = aa.apply(review, root, now=NOW)
            self.assertEqual(code, 0, errs)
            self.assertTrue(review["rows"][0]["applied"])
            fm, body = gr.split_frontmatter(card.read_text(encoding="utf-8"))
            assert fm
            self.assertEqual(fm["aa_intelligence"], "unknown")
            self.assertEqual(fm["refreshed"], "2026-11-05")
            self.assertIn("Keep me.", body)

    def test_clears_unmatched_when_previous_number(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            card = root / "multimodal" / "llama-4-scout-17b.md"
            card.parent.mkdir(parents=True)
            text = CARD.replace("name: claude-haiku-5.5", "name: llama-4-scout-17b")
            text = text.replace("aa_intelligence: unknown", "aa_intelligence: 10")
            card.write_text(text, encoding="utf-8")
            review = {
                "status": "ok",
                "rows": [
                    {
                        "name": "llama-4-scout-17b",
                        "status": "unmatched",
                        "aa_intelligence": None,
                        "applied": False,
                    }
                ],
            }
            code, _ = aa.apply(review, root, now=NOW)
            self.assertEqual(code, 0)
            self.assertTrue(review["rows"][0]["applied"])
            fm, _ = gr.split_frontmatter(card.read_text(encoding="utf-8"))
            assert fm
            self.assertEqual(fm["aa_intelligence"], "unknown")

    def test_skips_clear_when_already_unknown(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            card = root / "haiku.md"
            card.write_text(CARD, encoding="utf-8")
            review = {
                "status": "ok",
                "rows": [
                    {
                        "name": "claude-haiku-5.5",
                        "status": "no_score",
                        "aa_intelligence": None,
                        "applied": False,
                    }
                ],
            }
            code, _ = aa.apply(review, root, now=NOW)
            self.assertEqual(code, 0)
            self.assertFalse(review["rows"][0]["applied"])
            fm, _ = gr.split_frontmatter(card.read_text(encoding="utf-8"))
            assert fm
            self.assertEqual(str(fm["refreshed"]), "2026-11-04")

    def test_noop_when_review_skipped(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            card = root / "haiku.md"
            card.write_text(CARD, encoding="utf-8")
            code, errs = aa.apply({"status": "skipped_no_key", "rows": []}, root)
            self.assertEqual(code, 0, errs)
            fm, _ = gr.split_frontmatter(card.read_text(encoding="utf-8"))
            assert fm
            self.assertEqual(fm["aa_intelligence"], "unknown")


if __name__ == "__main__":
    unittest.main()
