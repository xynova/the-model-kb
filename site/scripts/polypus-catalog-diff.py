#!/usr/bin/env python3
"""Diff Polypus models lists against roster/catalog.json.

GETs /v1/apis/openai/models and /v1/apis/systemone/models by default. Prints
proposed model-roster keep/drop candidates. Does not write cards, mint hosts,
or treat the gateway as a priced host. Tests inject a models payload; live GET
is opt-in via main().
"""

from __future__ import annotations

import argparse
import json
import os
import sys
from pathlib import Path
from typing import Callable
from urllib import error, request

DEFAULT_BASE_URL = "http://127.0.0.1:1320"
API_PATHS: dict[str, str] = {
    "openai": "/v1/apis/openai/models",
    "systemone": "/v1/apis/systemone/models",
}
USER_AGENT = "the-model-kb-polypus-catalog-diff/1.0"
TIMEOUT_S = 20

ROUTE_PREFIX_TO_HOST: dict[str, str] = {
    "cf_local": "cloudflare",
    "deepinfra": "deepinfra",
    "openrouter": "openrouter",
    "gemini_free": "google",
    "gemini_paid": "google",
}

SKIP_NEEDLES: tuple[str, ...] = (
    "/bge-",
    "/bge/",
    "bge-m3",
    "bge-small",
    "whisper",
    "deepgram",
    "kokoro",
    "fish-audio",
    "transcribe",
    "aura-2",
    "smart-turn",
    "investigator",
)

STEM_SUFFIXES: tuple[str, ...] = ("-maas", "-it", "-instruct")


def models_url(base_url: str, api: str = "openai") -> str:
    raw = (base_url or "").strip().rstrip("/")
    if not raw:
        raw = DEFAULT_BASE_URL
    if raw.endswith("/apis/openai/models") or raw.endswith("/apis/systemone/models"):
        return raw
    path = API_PATHS.get(api, API_PATHS["openai"])
    if raw.endswith("/v1"):
        return raw + path.removeprefix("/v1")
    return raw + path


def default_models_urls(base_url: str) -> list[str]:
    return [models_url(base_url, "openai"), models_url(base_url, "systemone")]


def split_route(model_id: str) -> tuple[str, str, str]:
    text = (model_id or "").strip()
    if "/" not in text:
        return "", "", text
    prefix, remainder = text.split("/", 1)
    host = ROUTE_PREFIX_TO_HOST.get(prefix, "")
    if host:
        return prefix, host, remainder
    return "", "", text


def last_component(value: str) -> str:
    text = (value or "").strip().lower()
    if "/" in text:
        return text.rsplit("/", 1)[-1]
    return text


def stem(value: str) -> str:
    text = last_component(value)
    changed = True
    while changed:
        changed = False
        for suffix in STEM_SUFFIXES:
            if text.endswith(suffix) and len(text) > len(suffix):
                text = text[: -len(suffix)]
                changed = True
    return text


def is_skipped(model_id: str) -> bool:
    lowered = (model_id or "").lower()
    if lowered.startswith("router/"):
        return True
    return any(needle in lowered for needle in SKIP_NEEDLES)


def parse_models_payload(payload: dict) -> list[str]:
    if not isinstance(payload, dict):
        raise ValueError("models payload must be an object")
    data = payload.get("data")
    if not isinstance(data, list):
        raise ValueError("models payload missing data list")
    ids: list[str] = []
    seen: set[str] = set()
    for item in data:
        if not isinstance(item, dict):
            continue
        model_id = item.get("id")
        if not isinstance(model_id, str) or not model_id.strip():
            continue
        model_id = model_id.strip()
        if model_id in seen:
            continue
        seen.add(model_id)
        ids.append(model_id)
    return ids


def catalog_runtime_index(catalog: dict) -> list[dict]:
    rows: list[dict] = []
    cards = catalog.get("cards") or []
    if not isinstance(cards, list):
        return rows
    for card in cards:
        if not isinstance(card, dict):
            continue
        name = card.get("name")
        if not isinstance(name, str) or not name:
            continue
        for host in card.get("hosts") or []:
            if not isinstance(host, dict):
                continue
            rid = host.get("runtime_model_id")
            if not isinstance(rid, str) or not rid.strip():
                continue
            host_id = host.get("id")
            rows.append(
                {
                    "name": name,
                    "host_id": host_id if isinstance(host_id, str) else "",
                    "runtime_model_id": rid.strip(),
                }
            )
    return rows


def runtime_matches(remainder: str, runtime_model_id: str) -> bool:
    left = remainder.strip().lower()
    right = runtime_model_id.strip().lower()
    if not left or not right:
        return False
    if left == right:
        return True
    if last_component(left) == last_component(right) and last_component(left):
        return True
    left_stem = stem(left)
    right_stem = stem(right)
    return bool(left_stem and left_stem == right_stem)


def match_catalog_rows(remainder: str, catalog_rows: list[dict]) -> list[dict]:
    hits: list[dict] = []
    seen: set[tuple[str, str]] = set()
    for row in catalog_rows:
        if not runtime_matches(remainder, row["runtime_model_id"]):
            continue
        key = (row["name"], row["host_id"])
        if key in seen:
            continue
        seen.add(key)
        hits.append(row)
    return hits


def classify_models(model_ids: list[str], catalog: dict) -> dict:
    catalog_rows = catalog_runtime_index(catalog)
    matched: list[dict] = []
    proposed: list[dict] = []
    skipped: list[dict] = []
    seen_polypus: set[str] = set()

    for model_id in model_ids:
        if model_id in seen_polypus:
            continue
        seen_polypus.add(model_id)
        prefix, host_id, remainder = split_route(model_id)
        if is_skipped(model_id):
            skipped.append(
                {
                    "id": model_id,
                    "reason": "specialty",
                    "host_id": host_id,
                    "runtime_model_id": remainder,
                }
            )
            continue
        hits = match_catalog_rows(remainder, catalog_rows)
        if not hits and remainder != model_id:
            hits = match_catalog_rows(model_id, catalog_rows)
        if host_id:
            same_host = [row for row in hits if row["host_id"] == host_id]
            if same_host:
                hits = same_host
        if hits:
            matched.append(
                {
                    "id": model_id,
                    "host_id": host_id,
                    "runtime_model_id": remainder,
                    "cards": hits,
                }
            )
            continue
        proposed.append(
            {
                "id": model_id,
                "prefix": prefix,
                "host_id": host_id or "unknown",
                "runtime_model_id": remainder,
                "action": "model-roster keep/drop; do not mint polypus as a host",
            }
        )

    catalog_unseen: list[dict] = []
    for row in catalog_rows:
        found = False
        for model_id in model_ids:
            _prefix, _host, remainder = split_route(model_id)
            if runtime_matches(remainder, row["runtime_model_id"]):
                found = True
                break
            if runtime_matches(model_id, row["runtime_model_id"]):
                found = True
                break
        if not found:
            catalog_unseen.append(row)

    return {
        "diff_schema": 1,
        "matched": matched,
        "proposed": proposed,
        "skipped": skipped,
        "catalog_unseen": catalog_unseen,
        "counts": {
            "polypus": len(model_ids),
            "matched": len(matched),
            "proposed": len(proposed),
            "skipped": len(skipped),
            "catalog_runtime_ids": len(catalog_rows),
            "catalog_unseen": len(catalog_unseen),
        },
    }


def default_fetch(url: str) -> tuple[int, str]:
    req = request.Request(url, headers={"User-Agent": USER_AGENT})
    try:
        with request.urlopen(req, timeout=TIMEOUT_S) as resp:
            body = resp.read().decode("utf-8", errors="replace")
            return resp.getcode() or 0, body
    except error.HTTPError as exc:
        body = exc.read().decode("utf-8", errors="replace") if exc.fp else ""
        return exc.code, body
    except error.URLError as exc:
        raise RuntimeError(f"polypus fetch failed: {exc.reason}") from exc


def load_models(
    url: str,
    *,
    fetch: Callable[[str], tuple[int, str]] | None = None,
) -> list[str]:
    if fetch is None:
        fetch = default_fetch
    status, body = fetch(url)
    if status != 200 or not body.strip():
        raise RuntimeError(f"polypus models HTTP {status} at {url}")
    try:
        payload = json.loads(body)
    except json.JSONDecodeError as exc:
        raise RuntimeError(f"polypus models JSON parse failed: {exc}") from exc
    return parse_models_payload(payload)


def load_models_many(
    urls: list[str],
    *,
    fetch: Callable[[str], tuple[int, str]] | None = None,
) -> list[str]:
    ids: list[str] = []
    seen: set[str] = set()
    for url in urls:
        for model_id in load_models(url, fetch=fetch):
            if model_id in seen:
                continue
            seen.add(model_id)
            ids.append(model_id)
    return ids


def render_human(report: dict, *, url: str, catalog_path: Path) -> str:
    counts = report["counts"]
    lines = [
        "Polypus catalog diff",
        f"url: {url}",
        f"catalog: {catalog_path}",
        (
            f"polypus {counts['polypus']} | matched {counts['matched']} | "
            f"proposed {counts['proposed']} | skipped {counts['skipped']} | "
            f"catalog unseen {counts['catalog_unseen']}"
        ),
        "",
        "Prices and capacity stay on public host pages. This list is keep/drop only.",
        "",
        "## Proposed (not in catalog)",
    ]
    if not report["proposed"]:
        lines.append("(none)")
    for row in report["proposed"]:
        lines.append(f"- {row['id']}")
        lines.append(f"  host: {row['host_id']}")
        lines.append(f"  runtime_model_id: {row['runtime_model_id']}")
        lines.append(f"  action: {row['action']}")
    lines.append("")
    lines.append("## Matched")
    if not report["matched"]:
        lines.append("(none)")
    for row in report["matched"]:
        cards = ", ".join(
            f"{c['name']} / {c['host_id']}" for c in row["cards"]
        )
        lines.append(f"- {row['id']} -> {cards}")
    lines.append("")
    lines.append("## Catalog runtime ids not on this gateway")
    if not report["catalog_unseen"]:
        lines.append("(none)")
    for row in report["catalog_unseen"]:
        lines.append(
            f"- {row['name']} / {row['host_id']} ({row['runtime_model_id']})"
        )
    lines.append("")
    lines.append("## Skipped (embeddings, audio, router)")
    if not report["skipped"]:
        lines.append("(none)")
    for row in report["skipped"]:
        lines.append(f"- {row['id']}")
    lines.append("")
    return "\n".join(lines)


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(
        description="Diff Polypus model ids against roster/catalog.json"
    )
    parser.add_argument(
        "--catalog",
        type=Path,
        default=None,
        help="Path to catalog.json (default: roster/catalog.json)",
    )
    parser.add_argument(
        "--base-url",
        default=None,
        help="Polypus base URL (default: POLYPUS_BASE_URL or http://127.0.0.1:1320)",
    )
    parser.add_argument(
        "--url",
        default=None,
        help="Full models URL (overrides --base-url; skips the default openai+systemone pair)",
    )
    parser.add_argument(
        "--json",
        action="store_true",
        help="Print the report as JSON instead of sections",
    )
    args = parser.parse_args(argv)

    site_dir = Path(__file__).resolve().parents[1]
    repo_root = site_dir.parent
    catalog_path = args.catalog or (repo_root / "roster" / "catalog.json")
    if not catalog_path.is_file():
        print(f"catalog missing: {catalog_path}", file=sys.stderr)
        print("Run: make -C site generate-roster", file=sys.stderr)
        return 1

    with catalog_path.open(encoding="utf-8") as f:
        catalog = json.load(f)

    if args.url:
        urls = [args.url.strip()]
    else:
        base = (args.base_url or os.environ.get("POLYPUS_BASE_URL") or DEFAULT_BASE_URL)
        urls = default_models_urls(base)
    url_label = " ".join(urls)

    try:
        model_ids = load_models_many(urls)
    except RuntimeError as exc:
        print(str(exc), file=sys.stderr)
        return 1

    report = classify_models(model_ids, catalog)
    report["url"] = url_label
    report["urls"] = urls
    if args.json:
        print(json.dumps(report, indent=2, ensure_ascii=False))
        return 0
    print(render_human(report, url=url_label, catalog_path=catalog_path), end="")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
