"""Tests for generate-roster.py."""

from __future__ import annotations

import importlib.util
import unittest
from pathlib import Path

SCRIPTS = Path(__file__).resolve().parent
SPEC = importlib.util.spec_from_file_location(
    "generate_roster", SCRIPTS / "generate-roster.py"
)
assert SPEC and SPEC.loader
gr = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(gr)

FIXTURE_CARD = """---
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
"""


class GenerateRosterTests(unittest.TestCase):
    def test_remap_kind_only_top_level(self) -> None:
        fm, _ = gr.split_frontmatter(FIXTURE_CARD)
        self.assertIsNotNone(fm)
        dumped = gr.dump_hugo_frontmatter(fm)
        self.assertIn("roster_kind: sovereign", dumped)
        self.assertNotIn("\nkind:", dumped)

    def test_hosts_survive_generation(self) -> None:
        fm, body = gr.split_frontmatter(FIXTURE_CARD)
        self.assertIsNotNone(fm)
        gr.validate_card(fm, Path("cheap/test-model.md"))
        text = "---\n" + gr.dump_hugo_frontmatter(fm) + "\n---\n\n" + body
        self.assertIn("hosts:", text)
        out_fm, _ = gr.split_frontmatter(text)
        self.assertIsNotNone(out_fm)
        self.assertEqual(len(out_fm["hosts"]), 1)
        self.assertEqual(out_fm["hosts"][0]["id"], "cloudflare")

    def test_duplicate_host_id_fails(self) -> None:
        bad = FIXTURE_CARD.replace(
            "capacity: hobby",
            "capacity: hobby\n  - id: cloudflare\n    runtime: cloudflare\n    capacity: hobby",
        )
        fm, _ = gr.split_frontmatter(bad)
        self.assertIsNotNone(fm)
        with self.assertRaises(ValueError):
            gr.validate_card(fm, Path("x.md"))

    def test_objective_in_skip_list(self) -> None:
        self.assertIn("OBJECTIVE.md", gr.SKIP_NAMES)

    def test_hugo_frontmatter_quotes_aa_intelligence(self) -> None:
        fm, _ = gr.split_frontmatter(FIXTURE_CARD)
        self.assertIsNotNone(fm)
        fm["aa_intelligence"] = 43.6
        dumped = gr.dump_hugo_frontmatter(fm)
        self.assertRegex(dumped, r"aa_intelligence: ['\"]43\.6['\"]")
        fm["aa_intelligence"] = 52
        dumped = gr.dump_hugo_frontmatter(fm)
        self.assertRegex(dumped, r"aa_intelligence: ['\"]52['\"]")

    def test_allowed_categories(self) -> None:
        self.assertEqual(
            gr.ALLOWED_CATEGORIES,
            frozenset({"cheap", "mid", "spend", "unpriced", "classify"}),
        )
        fm, _ = gr.split_frontmatter(FIXTURE_CARD)
        assert fm
        fm["category"] = "classify"
        gr.validate_card(fm, Path("classify/x.md"))
        fm["category"] = "budget"
        with self.assertRaises(ValueError):
            gr.validate_card(fm, Path("budget/x.md"))

    def test_assign_category(self) -> None:
        def fm(**kwargs: object) -> dict:
            base = {
                "kind": "sovereign",
                "jobs": ["chat"],
                "hosts": [{"id": "x", "price_in": 1.0, "price_out": 0.55}],
            }
            base.update(kwargs)
            return base

        self.assertEqual(gr.assign_category(fm(jobs=["classify", "cheap"])), "classify")
        self.assertEqual(gr.assign_category(fm(kind="service", jobs=["coding"])), "unpriced")
        self.assertEqual(
            gr.assign_category(fm(hosts=[{"price_in": 1.0, "price_out": 0.55}])),
            "cheap",
        )
        self.assertEqual(
            gr.assign_category(fm(hosts=[{"price_in": 10.0, "price_out": 4.0}])),
            "mid",
        )
        self.assertEqual(
            gr.assign_category(fm(hosts=[{"price_in": 10.0, "price_out": 10.0}])),
            "spend",
        )
        self.assertEqual(gr.assign_category(fm(hosts=[{"price_note": "on request"}])), "unpriced")

    def test_card_to_json(self) -> None:
        fm, _ = gr.split_frontmatter(FIXTURE_CARD)
        self.assertIsNotNone(fm)
        j = gr.card_to_json(fm)
        self.assertEqual(j["name"], "test-model")
        self.assertEqual(j["hosts"][0]["capacity"], "hobby")
        self.assertEqual(j["preferred_host"], "cloudflare")


if __name__ == "__main__":
    unittest.main()
