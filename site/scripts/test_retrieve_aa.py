"""Tests for retrieve-aa.py (no network)."""

from __future__ import annotations

import importlib.util
import unittest
from datetime import datetime, timezone
from pathlib import Path

SCRIPTS = Path(__file__).resolve().parent
SPEC = importlib.util.spec_from_file_location("retrieve_aa", SCRIPTS / "retrieve-aa.py")
assert SPEC and SPEC.loader
ra = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(ra)

NOW = datetime(2026, 11, 5, tzinfo=timezone.utc)


def aa_model(
    name: str,
    slug: str,
    score: float | None,
    *,
    creator: str = "Lab",
    cost: float | None = None,
) -> dict:
    payload: dict = {
        "name": name,
        "slug": slug,
        "model_creator": {"name": creator},
        "evaluations": {"artificial_analysis_intelligence_index": score},
        "artificial_analysis_intelligence_index_cost": None,
    }
    if cost is not None:
        payload["artificial_analysis_intelligence_index_cost"] = {
            "cost_per_task": {"total_cost": cost}
        }
    return payload


class TestEffortPick(unittest.TestCase):
    def test_prefers_max_over_high(self) -> None:
        catalog = {
            "cards": [
                {
                    "name": "claude-haiku-5.5",
                    "title": "Claude Haiku 5.5",
                    "kind": "provider",
                    "aa_intelligence": "unknown",
                }
            ]
        }
        models = [
            aa_model("Claude Haiku 5.5 (high)", "claude-haiku-5-5", 30),
            aa_model("Claude Haiku 5.5 (max)", "claude-haiku-5-5", 36),
        ]
        rows = ra.match_cards(catalog, models)
        self.assertEqual(rows[0]["status"], "ok")
        self.assertEqual(rows[0]["aa_intelligence"], 36)
        self.assertEqual(rows[0]["effort"], "max")

    def test_skips_non_reasoning_when_thinking_exists(self) -> None:
        catalog = {
            "cards": [
                {
                    "name": "gemini-3.8-flash",
                    "title": "Gemini 3.8 Flash",
                    "kind": "provider",
                }
            ]
        }
        models = [
            aa_model("Gemini 3.8 Flash (Non-reasoning)", "gemini-3-8-flash", 12),
            aa_model("Gemini 3.8 Flash (High)", "gemini-3-8-flash", 31),
        ]
        rows = ra.match_cards(catalog, models)
        self.assertEqual(rows[0]["status"], "ok")
        self.assertEqual(rows[0]["aa_intelligence"], 31)
        self.assertEqual(rows[0]["effort"], "high")

    def test_unmatched_when_only_non_reasoning(self) -> None:
        catalog = {
            "cards": [
                {
                    "name": "gemini-2.5-flash",
                    "title": "Gemini 2.5 Flash",
                    "kind": "provider",
                }
            ]
        }
        models = [
            aa_model(
                "Gemini 2.5 Flash (Non-reasoning)",
                "gemini-2-5-flash",
                20,
            )
        ]
        rows = ra.match_cards(catalog, models)
        self.assertEqual(rows[0]["status"], "unmatched")
        self.assertIsNone(rows[0]["aa_intelligence"])


class TestFlashDoesNotTakeBaseMax(unittest.TestCase):
    def test_glm_flash_ignores_glm_max(self) -> None:
        catalog = {
            "cards": [
                {
                    "name": "glm-5.3-flash",
                    "title": "GLM-5.3-Flash",
                    "kind": "sovereign",
                }
            ]
        }
        models = [aa_model("GLM-5.3 (Max)", "glm-5-3", 32)]
        rows = ra.match_cards(catalog, models)
        self.assertEqual(rows[0]["status"], "unmatched")


class TestServiceSkipped(unittest.TestCase):
    def test_cursor_card_skipped(self) -> None:
        catalog = {
            "cards": [
                {
                    "name": "cursor-gemini-3-flash",
                    "title": "Gemini 3 Flash",
                    "kind": "service",
                    "service": "cursor",
                }
            ]
        }
        models = [aa_model("Gemini 3 Flash (high)", "gemini-3-flash", 20)]
        rows = ra.match_cards(catalog, models)
        self.assertEqual(rows[0]["status"], "skipped_service")


class TestNoScore(unittest.TestCase):
    def test_listed_null_stays_no_score(self) -> None:
        catalog = {
            "cards": [
                {
                    "name": "gemma-4-26b-a4b",
                    "title": "Gemma 4 26B A4B",
                    "kind": "sovereign",
                    "aa_intelligence": 31,
                }
            ]
        }
        models = [aa_model("Gemma 4 26B A4B", "gemma-4-26b-a4b", None)]
        rows = ra.match_cards(catalog, models)
        self.assertEqual(rows[0]["status"], "no_score")
        self.assertIsNone(rows[0]["aa_intelligence"])


class TestRetrieveStatus(unittest.TestCase):
    def test_skip_without_key(self) -> None:
        out = ra.retrieve({"cards": []}, api_key="", now=NOW)
        self.assertEqual(out["status"], "skipped_no_key")
        self.assertEqual(out["rows"], [])

    def test_skip_other_index_major(self) -> None:
        def fetch_page(_page: int, _key: str) -> dict:
            return {
                "intelligence_index_version": 6,
                "pagination": {"has_more": False},
                "data": [aa_model("Claude Haiku 5.5 (max)", "claude-haiku-5-5", 36)],
            }

        out = ra.retrieve(
            {"cards": []},
            api_key="test-key",
            fetch_page=fetch_page,
            now=NOW,
            sleep=lambda _: None,
        )
        self.assertEqual(out["status"], "skipped_version")

    def test_paginates_and_matches(self) -> None:
        pages = {
            1: {
                "intelligence_index_version": 5,
                "pagination": {"has_more": True},
                "data": [aa_model("Granite 4.2 8B", "granite-4-2-8b", 5)],
            },
            2: {
                "intelligence_index_version": 5,
                "pagination": {"has_more": False},
                "data": [
                    aa_model("Gemini 3.8 Flash (high)", "gemini-3-8-flash", 31)
                ],
            },
        }

        def fetch_page(page: int, key: str) -> dict:
            self.assertEqual(key, "test-key")
            return pages[page]

        catalog = {
            "cards": [
                {
                    "name": "granite-4.2-8b",
                    "title": "Granite 4.2 8B",
                    "kind": "sovereign",
                },
                {
                    "name": "gemini-3.8-flash",
                    "title": "Gemini 3.8 Flash",
                    "kind": "provider",
                },
            ]
        }
        out = ra.retrieve(
            catalog,
            api_key="test-key",
            fetch_page=fetch_page,
            now=NOW,
            sleep=lambda _: None,
        )
        self.assertEqual(out["status"], "ok")
        by_name = {row["name"]: row for row in out["rows"]}
        self.assertEqual(by_name["granite-4.2-8b"]["aa_intelligence"], 5)
        self.assertEqual(by_name["gemini-3.8-flash"]["aa_intelligence"], 31)
        self.assertEqual(by_name["gemini-3.8-flash"]["effort"], "high")


class TestComparePeers(unittest.TestCase):
    def test_alternative_beats_intelligence_per_cost(self) -> None:
        catalog = {
            "cards": [
                {
                    "name": "claude-haiku-5.5",
                    "title": "Claude Haiku 5.5",
                    "kind": "provider",
                }
            ]
        }
        models = [
            aa_model(
                "Claude Haiku 5.5 (Max)",
                "claude-haiku-5-5",
                36.5,
                creator="Anthropic",
                cost=0.41,
            ),
            aa_model(
                "Better Peer (Max)",
                "better-peer",
                40,
                creator="DeepSeek",
                cost=0.20,
            ),
            aa_model(
                "Smarter Spend (Max)",
                "smarter-spend",
                50,
                creator="Google",
                cost=0.50,
            ),
            aa_model(
                "Cheaper Dumber (Max)",
                "cheaper-dumber",
                10,
                creator="IBM",
                cost=0.01,
            ),
        ]
        rows = ra.match_cards(catalog, models)
        alt = rows[0]["alternative"]
        self.assertIsNotNone(alt)
        self.assertEqual(alt["name"], "Better Peer (Max)")
        self.assertEqual(alt["creator"], "DeepSeek")
        self.assertEqual(alt["aa_intelligence"], 40)
        self.assertEqual(alt["cost_per_task"], 0.2)
        self.assertEqual(rows[0]["aa_cost_per_task"], 0.41)
        self.assertEqual(rows[0]["aa_creator"], "Anthropic")
        payload = ra.compare_data({"generated_at": "2026-11-05", "rows": rows})
        entry = payload["cards"]["claude-haiku-5.5"]
        self.assertEqual(entry["cost_per_task"], 0.41)
        self.assertEqual(entry["aa_intelligence"], 36.5)
        self.assertEqual(entry["creator"], "Anthropic")
        self.assertEqual(entry["alternative"]["name"], "Better Peer (Max)")

    def test_alternative_skips_smarter_but_costlier(self) -> None:
        catalog = {
            "cards": [
                {
                    "name": "claude-haiku-5.5",
                    "title": "Claude Haiku 5.5",
                    "kind": "provider",
                }
            ]
        }
        models = [
            aa_model(
                "Claude Haiku 5.5 (Max)",
                "claude-haiku-5-5",
                36.5,
                creator="Anthropic",
                cost=0.41,
            ),
            aa_model(
                "Smarter Spend (Max)",
                "smarter-spend",
                50,
                creator="Google",
                cost=0.50,
            ),
        ]
        rows = ra.match_cards(catalog, models)
        self.assertIsNone(rows[0]["alternative"])

    def test_similar_is_closest_off_roster_index(self) -> None:
        catalog = {
            "cards": [
                {
                    "name": "claude-haiku-5.5",
                    "title": "Claude Haiku 5.5",
                    "kind": "provider",
                }
            ]
        }
        models = [
            aa_model(
                "Claude Haiku 5.5 (Max)",
                "claude-haiku-5-5",
                36.5,
                creator="Anthropic",
                cost=0.41,
            ),
            aa_model(
                "Better Peer (Max)",
                "better-peer",
                40,
                creator="DeepSeek",
                cost=0.20,
            ),
            aa_model(
                "Near Peer (High)",
                "near-peer",
                38,
                creator="Z AI",
                cost=1.0,
            ),
        ]
        rows = ra.match_cards(catalog, models)
        self.assertEqual(rows[0]["alternative"]["slug"], "better-peer")
        similar = rows[0]["similar"]
        self.assertEqual(similar["name"], "Near Peer (High)")
        self.assertEqual(similar["creator"], "Z AI")

    def test_skips_other_roster_cards(self) -> None:
        catalog = {
            "cards": [
                {
                    "name": "muse-wonder",
                    "title": "Muse Wonder",
                    "kind": "provider",
                },
                {
                    "name": "claude-haiku-5.5",
                    "title": "Claude Haiku 5.5",
                    "kind": "provider",
                },
            ]
        }
        models = [
            aa_model(
                "Muse Wonder (Max)",
                "muse-wonder",
                43.6,
                creator="Meta",
                cost=6.95,
            ),
            aa_model(
                "Claude Haiku 5.5 (Max)",
                "claude-haiku-5-5",
                36.5,
                creator="Anthropic",
                cost=0.41,
            ),
        ]
        rows = ra.match_cards(catalog, models)
        by_name = {row["name"]: row for row in rows}
        self.assertIsNone(by_name["muse-wonder"]["alternative"])
        self.assertIsNone(by_name["muse-wonder"]["similar"])


if __name__ == "__main__":
    unittest.main()
