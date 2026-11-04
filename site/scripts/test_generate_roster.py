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
category: general
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
        gr.validate_card(fm, Path("general/test-model.md"))
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

    def test_card_to_json(self) -> None:
        fm, _ = gr.split_frontmatter(FIXTURE_CARD)
        self.assertIsNotNone(fm)
        j = gr.card_to_json(fm)
        self.assertEqual(j["name"], "test-model")
        self.assertEqual(j["hosts"][0]["capacity"], "hobby")
        self.assertEqual(j["preferred_host"], "cloudflare")


if __name__ == "__main__":
    unittest.main()
