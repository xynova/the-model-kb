"""Tests for apply-mine.py."""

from __future__ import annotations

import importlib.util
import json
import tempfile
import unittest
from datetime import datetime, timezone
from pathlib import Path

SCRIPTS = Path(__file__).resolve().parent
APPLY_SPEC = importlib.util.spec_from_file_location(
    "apply_mine", SCRIPTS / "apply-mine.py"
)
assert APPLY_SPEC and APPLY_SPEC.loader
am = importlib.util.module_from_spec(APPLY_SPEC)
APPLY_SPEC.loader.exec_module(am)

GR_SPEC = importlib.util.spec_from_file_location(
    "generate_roster", SCRIPTS / "generate-roster.py"
)
assert GR_SPEC and GR_SPEC.loader
gr = importlib.util.module_from_spec(GR_SPEC)
GR_SPEC.loader.exec_module(gr)

NOW = datetime(2026, 3, 1, tzinfo=timezone.utc)

GEMMA_CARD = """---
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
"""


class TestUnconfirmedKeepsCapacity(unittest.TestCase):
    def test_cloudflare_stays_hobby(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            card = root / "multimodal" / "gemma.md"
            card.parent.mkdir(parents=True)
            card.write_text(GEMMA_CARD, encoding="utf-8")
            review = {
                "rows": [
                    {
                        "status": "ok",
                        "name": "gemma-4-26b-a4b",
                        "host_id": "cloudflare",
                        "evidence": "10k Neurons/day free tier",
                        "price_in": None,
                        "price_out": None,
                    }
                ]
            }
            code, errs = am.apply(review, [], root, now=NOW)
            self.assertEqual(code, 0, errs)
            text = card.read_text(encoding="utf-8")
            fm, _ = gr.split_frontmatter(text)
            assert fm
            cf = next(h for h in fm["hosts"] if h["id"] == "cloudflare")
            self.assertEqual(cf["capacity"], "hobby")
            self.assertEqual(cf["limits_note"], "10k Neurons/day free tier")


class TestConfirmedSetsCapacity(unittest.TestCase):
    def test_deepinfra_industrial(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            card = root / "gemma.md"
            card.write_text(GEMMA_CARD, encoding="utf-8")
            review = {
                "rows": [
                    {
                        "status": "ok",
                        "name": "gemma-4-26b-a4b",
                        "host_id": "deepinfra",
                        "evidence": "production throughput",
                        "price_in": 0.08,
                        "price_out": 0.35,
                        "capacity_applied": False,
                    }
                ]
            }
            decisions = [
                {
                    "name": "gemma-4-26b-a4b",
                    "host_id": "deepinfra",
                    "capacity": "industrial",
                }
            ]
            code, _ = am.apply(review, decisions, root, now=NOW)
            self.assertEqual(code, 0)
            fm, _ = gr.split_frontmatter(card.read_text(encoding="utf-8"))
            assert fm
            di = next(h for h in fm["hosts"] if h["id"] == "deepinfra")
            self.assertEqual(di["capacity"], "industrial")
            self.assertTrue(review["rows"][0].get("capacity_applied"))


class TestMissingPriceDoesNotClear(unittest.TestCase):
    def test_null_prices_preserved(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            card = root / "gemma.md"
            card.write_text(GEMMA_CARD, encoding="utf-8")
            review = {
                "rows": [
                    {
                        "status": "ok",
                        "name": "gemma-4-26b-a4b",
                        "host_id": "cloudflare",
                        "evidence": "neurons",
                        "price_in": None,
                        "price_out": None,
                    }
                ]
            }
            am.apply(review, [], root, now=NOW)
            fm, _ = gr.split_frontmatter(card.read_text(encoding="utf-8"))
            assert fm
            cf = next(h for h in fm["hosts"] if h["id"] == "cloudflare")
            self.assertEqual(cf["price_in"], 0.10)


class TestDumpKeepsKind(unittest.TestCase):
    def test_no_roster_kind(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            card = root / "gemma.md"
            card.write_text(GEMMA_CARD, encoding="utf-8")
            review = {
                "rows": [
                    {
                        "status": "ok",
                        "name": "gemma-4-26b-a4b",
                        "host_id": "cloudflare",
                        "evidence": "rpm",
                    }
                ]
            }
            am.apply(review, [], root, now=NOW)
            text = card.read_text(encoding="utf-8")
            self.assertIn("kind: sovereign", text)
            self.assertNotIn("roster_kind", text)


class TestIllegalCapacityFails(unittest.TestCase):
    def test_load_decisions_rejects_hosted(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            path = Path(tmp) / "decisions.yaml"
            path.write_text(
                "decisions:\n"
                "  - name: gemma-4-26b-a4b\n"
                "    host_id: deepinfra\n"
                "    capacity: hosted\n",
                encoding="utf-8",
            )
            with self.assertRaises(ValueError):
                am.load_decisions(path)


if __name__ == "__main__":
    unittest.main()
