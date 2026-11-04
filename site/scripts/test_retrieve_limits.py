"""Tests for retrieve-limits.py (no network)."""

from __future__ import annotations

import importlib.util
import unittest
from datetime import datetime, timezone
from pathlib import Path

SCRIPTS = Path(__file__).resolve().parent
SPEC = importlib.util.spec_from_file_location(
    "retrieve_limits", SCRIPTS / "retrieve-limits.py"
)
assert SPEC and SPEC.loader
rl = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(rl)

NOW = datetime(2026, 3, 1, 12, 0, 0, tzinfo=timezone.utc)


def catalog_with_host(host: dict, *, kind: str = "sovereign", name: str = "test-card") -> dict:
    return {
        "cards": [
            {
                "name": name,
                "kind": kind,
                "hosts": [host],
            }
        ]
    }


class TestCloudflareNeuronHobby(unittest.TestCase):
    def test_neurons_suggest_hobby(self) -> None:
        html = "<p>Workers AI includes 10,000 Neurons per day on the free tier.</p>"

        def fetch(_url: str) -> tuple[int, str]:
            return 200, html

        cat = catalog_with_host(
            {
                "id": "cloudflare",
                "capacity": "unknown",
                "links": {
                    "cloudflare_pricing": "https://developers.cloudflare.com/workers-ai/pricing/"
                },
            },
            name="gemma-4-26b-a4b",
        )
        out = rl.retrieve(cat, fetch=fetch, now=NOW, sleep=lambda _: None)
        row = out["rows"][0]
        self.assertEqual(row["status"], "ok")
        self.assertEqual(row["suggested_capacity"], "hobby")
        self.assertTrue(row["evidence"])


class TestOpenRouterMarketplaceHobby(unittest.TestCase):
    def test_no_production_rpm(self) -> None:
        html = "<p>OpenRouter routes to many providers. Pay per token.</p>"

        def fetch(_url: str) -> tuple[int, str]:
            return 200, html

        cat = catalog_with_host(
            {
                "id": "openrouter",
                "capacity": "hobby",
                "links": {"openrouter": "https://openrouter.ai/docs/pricing"},
            }
        )
        out = rl.retrieve(cat, fetch=fetch, now=NOW, sleep=lambda _: None)
        row = out["rows"][0]
        self.assertEqual(row["suggested_capacity"], "hobby")
        self.assertEqual(row["suggestion_reason"], "openrouter_marketplace")


class TestGoogleNoRpmUnknown(unittest.TestCase):
    def test_provider_pricing_unknown(self) -> None:
        html = "<p>Vertex AI pricing varies by region and model.</p>"

        def fetch(_url: str) -> tuple[int, str]:
            return 200, html

        cat = catalog_with_host(
            {
                "id": "google",
                "capacity": "unknown",
                "links": {
                    "provider_pricing": "https://cloud.google.com/vertex-ai/generative-ai/pricing"
                },
            }
        )
        out = rl.retrieve(cat, fetch=fetch, now=NOW, sleep=lambda _: None)
        row = out["rows"][0]
        self.assertEqual(row["suggested_capacity"], "unknown")
        self.assertEqual(row["suggestion_reason"], "google_no_rpm")


class TestDeepinfraFlexUnknown(unittest.TestCase):
    def test_priority_flex_without_production(self) -> None:
        html = "<p>Priority and Flex tiers for burst throughput.</p>"

        def fetch(_url: str) -> tuple[int, str]:
            return 200, html

        cat = catalog_with_host(
            {
                "id": "deepinfra",
                "capacity": "unknown",
                "links": {"deepinfra": "https://deepinfra.com/pricing"},
            }
        )
        out = rl.retrieve(cat, fetch=fetch, now=NOW, sleep=lambda _: None)
        row = out["rows"][0]
        self.assertEqual(row["suggested_capacity"], "unknown")
        self.assertEqual(row["suggestion_reason"], "deepinfra_no_rpm")


class TestDeepinfraProductionIndustrial(unittest.TestCase):
    def test_production_text(self) -> None:
        html = "<p>Our production API offers high-throughput for enterprise.</p>"

        def fetch(_url: str) -> tuple[int, str]:
            return 200, html

        cat = catalog_with_host(
            {
                "id": "deepinfra",
                "capacity": "unknown",
                "links": {"deepinfra": "https://deepinfra.com/pricing"},
            }
        )
        out = rl.retrieve(cat, fetch=fetch, now=NOW, sleep=lambda _: None)
        row = out["rows"][0]
        self.assertEqual(row["suggested_capacity"], "industrial")
        self.assertEqual(row["suggestion_reason"], "deepinfra_production_text")


class TestSkipCursor(unittest.TestCase):
    def test_service_and_cursor_never_fetch(self) -> None:
        called = False

        def fetch(_url: str) -> tuple[int, str]:
            nonlocal called
            called = True
            return 200, "neurons"

        cat_service = {
            "cards": [
                {
                    "name": "cursor-pro",
                    "kind": "service",
                    "hosts": [
                        {
                            "id": "cursor",
                            "capacity": "unknown",
                            "links": {"service_docs": "https://cursor.com/docs"},
                        }
                    ],
                }
            ]
        }
        rl.retrieve(cat_service, fetch=fetch, now=NOW, sleep=lambda _: None)
        self.assertFalse(called)


class TestRejectOffAllowlist(unittest.TestCase):
    def test_wrong_hostname_skipped(self) -> None:
        called = False

        def fetch(_url: str) -> tuple[int, str]:
            nonlocal called
            called = True
            return 200, "rpm"

        cat = catalog_with_host(
            {
                "id": "cloudflare",
                "capacity": "hobby",
                "links": {"cloudflare_pricing": "https://example.com/pricing"},
            }
        )
        out = rl.retrieve(cat, fetch=fetch, now=NOW, sleep=lambda _: None)
        row = out["rows"][0]
        self.assertEqual(row["status"], "skipped")
        self.assertFalse(called)


if __name__ == "__main__":
    unittest.main()
