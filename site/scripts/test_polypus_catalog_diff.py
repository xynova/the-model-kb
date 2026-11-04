"""Tests for polypus-catalog-diff.py (no network)."""

from __future__ import annotations

import importlib.util
import json
import unittest
from pathlib import Path

SCRIPTS = Path(__file__).resolve().parent
SPEC = importlib.util.spec_from_file_location(
    "polypus_catalog_diff", SCRIPTS / "polypus-catalog-diff.py"
)
assert SPEC and SPEC.loader
pcd = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(pcd)


def catalog(*hosts: dict, name: str = "gemma-4-26b-a4b") -> dict:
    return {"cards": [{"name": name, "hosts": list(hosts)}]}


class TestModelsUrl(unittest.TestCase):
    def test_appends_models_path(self) -> None:
        self.assertEqual(
            pcd.models_url("http://127.0.0.1:1320"),
            "http://127.0.0.1:1320/v1/apis/openai/models",
        )

    def test_accepts_full_models_url(self) -> None:
        full = "https://example.ts.net/v1/apis/openai/models"
        self.assertEqual(pcd.models_url(full), full)

    def test_accepts_v1_base(self) -> None:
        self.assertEqual(
            pcd.models_url("http://127.0.0.1:1320/v1"),
            "http://127.0.0.1:1320/v1/apis/openai/models",
        )


class TestMatchCloudflareRoute(unittest.TestCase):
    def test_cf_local_matches_at_cf_id(self) -> None:
        cat = catalog(
            {
                "id": "cloudflare",
                "runtime_model_id": "@cf/google/gemma-4-26b-a4b-it",
            },
            {
                "id": "openrouter",
                "runtime_model_id": "google/gemma-4-26b-a4b-it",
            },
        )
        report = pcd.classify_models(
            ["cf_local/@cf/google/gemma-4-26b-a4b-it"], cat
        )
        self.assertEqual(report["counts"]["matched"], 1)
        self.assertEqual(report["counts"]["proposed"], 0)
        self.assertEqual(report["matched"][0]["cards"][0]["name"], "gemma-4-26b-a4b")
        self.assertEqual(report["matched"][0]["cards"][0]["host_id"], "cloudflare")
        self.assertEqual(len(report["matched"][0]["cards"]), 1)


class TestMatchDeepinfraAndMaasStem(unittest.TestCase):
    def test_deepinfra_prefix_and_maas_stem(self) -> None:
        cat = catalog(
            {
                "id": "deepinfra",
                "runtime_model_id": "google/gemma-4-26B-A4B-it",
            },
            {
                "id": "google",
                "runtime_model_id": "gemma-4-26b-a4b-it-maas",
            },
        )
        report = pcd.classify_models(
            [
                "deepinfra/google/gemma-4-26B-A4B-it",
                "gemini_free/gemma-4-26b-a4b-it",
            ],
            cat,
        )
        self.assertEqual(report["counts"]["matched"], 2)
        self.assertEqual(report["counts"]["proposed"], 0)


class TestSkipSpecialty(unittest.TestCase):
    def test_embeddings_and_audio_are_skipped(self) -> None:
        cat = catalog({"id": "cloudflare", "runtime_model_id": "@cf/other"})
        report = pcd.classify_models(
            [
                "cf_local/@cf/baai/bge-m3",
                "cf_local/@cf/openai/whisper-large-v3-turbo",
                "cf_local/@cf/deepgram/nova-3",
                "openrouter/hexgrad/kokoro-82m",
                "router/investigator",
            ],
            cat,
        )
        self.assertEqual(report["counts"]["skipped"], 5)
        self.assertEqual(report["counts"]["proposed"], 0)


class TestProposeUnmatchedLlm(unittest.TestCase):
    def test_unmatched_chat_model_is_proposed(self) -> None:
        cat = catalog(
            {
                "id": "cloudflare",
                "runtime_model_id": "@cf/google/gemma-4-26b-a4b-it",
            }
        )
        report = pcd.classify_models(
            ["deepinfra/deepseek-ai/DeepSeek-V4-Flash-0731"], cat
        )
        self.assertEqual(len(report["proposed"]), 1)
        row = report["proposed"][0]
        self.assertEqual(row["host_id"], "deepinfra")
        self.assertEqual(
            row["runtime_model_id"], "deepseek-ai/DeepSeek-V4-Flash-0731"
        )


class TestGeminiPreviewNotCollapsedToCursor(unittest.TestCase):
    def test_preview_id_stays_proposed(self) -> None:
        cat = {
            "cards": [
                {
                    "name": "cursor-gemini-3-flash",
                    "hosts": [
                        {
                            "id": "cursor",
                            "runtime_model_id": "gemini-3-flash",
                        }
                    ],
                }
            ]
        }
        report = pcd.classify_models(
            ["gemini_paid/gemini-3-flash-preview"], cat
        )
        self.assertEqual(report["counts"]["proposed"], 1)
        self.assertEqual(report["counts"]["matched"], 0)


class TestCatalogUnseen(unittest.TestCase):
    def test_catalog_id_missing_from_gateway(self) -> None:
        cat = catalog(
            {
                "id": "cloudflare",
                "runtime_model_id": "@cf/zai-org/glm-5.3-flash",
            },
            name="glm-5.3-flash",
        )
        report = pcd.classify_models(
            ["cf_local/@cf/google/gemma-4-26b-a4b-it"], cat
        )
        self.assertEqual(report["counts"]["catalog_unseen"], 1)
        self.assertEqual(report["catalog_unseen"][0]["name"], "glm-5.3-flash")


class TestLoadModelsInjectedFetch(unittest.TestCase):
    def test_ok_payload(self) -> None:
        body = json.dumps(
            {"object": "list", "data": [{"id": "cf_local/@cf/google/gemma-4-26b-a4b-it"}]}
        )

        def fetch(_url: str) -> tuple[int, str]:
            return 200, body

        ids = pcd.load_models("http://example/v1/apis/openai/models", fetch=fetch)
        self.assertEqual(ids, ["cf_local/@cf/google/gemma-4-26b-a4b-it"])

    def test_http_error(self) -> None:
        def fetch(_url: str) -> tuple[int, str]:
            return 500, ""

        with self.assertRaises(RuntimeError):
            pcd.load_models("http://example/v1/apis/openai/models", fetch=fetch)


if __name__ == "__main__":
    unittest.main()
