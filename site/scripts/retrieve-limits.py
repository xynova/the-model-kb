#!/usr/bin/env python3
"""Fetch public host pages and write roster/retrieval/limits-review.json."""

from __future__ import annotations

import json
import os
import re
import sys
import time
from datetime import datetime, timezone
from pathlib import Path
from typing import Callable
from urllib import error, parse, request

ALLOWED_CAPACITY = frozenset({"hobby", "industrial", "local", "unknown"})
SKIP_HOST_IDS = frozenset({"cursor", "local"})

LINK_KEYS: dict[str, str] = {
    "cloudflare_pricing": "developers.cloudflare.com",
    "openrouter": "openrouter.ai",
    "provider_pricing": "cloud.google.com",
    "deepinfra": "deepinfra.com",
    "xinference": "xinference.co",
}

EVIDENCE_PATTERN = re.compile(
    r".{0,240}(?:neuron|rpm|tpm|rate limit|free tier|quota|production|"
    r"throughput|flex|priority).{0,240}",
    re.IGNORECASE | re.DOTALL,
)
HTML_TAG = re.compile(r"<[^>]+>")
PRODUCTION_RPM = re.compile(
    r"production.{0,80}(?:rpm|throughput|high.?throughput)",
    re.IGNORECASE | re.DOTALL,
)
HIGH_THROUGHPUT = re.compile(r"high[- ]throughput", re.IGNORECASE)
PRODUCTION_WORD = re.compile(r"\bproduction\b", re.IGNORECASE)

USER_AGENT = "the-model-kb-roster-mine/1.0"
MAX_BODY = 2 * 1024 * 1024


def strip_html(html: str) -> str:
    text = HTML_TAG.sub(" ", html)
    return re.sub(r"\s+", " ", text).strip()


def extract_evidence(text: str) -> str:
    match = EVIDENCE_PATTERN.search(text)
    if not match:
        return ""
    snippet = match.group(0).strip()
    if len(snippet) > 240:
        snippet = snippet[:240].rstrip()
    return snippet


def hostname_allowed(url: str, expected_host: str) -> bool:
    parsed = parse.urlparse(url)
    host = (parsed.hostname or "").lower()
    if host.startswith("www."):
        host = host[4:]
    expected = expected_host.lower()
    return host == expected or host.endswith("." + expected)


def pick_allowlisted_url(links: dict) -> tuple[str, str] | None:
    if not isinstance(links, dict):
        return None
    for key, expected_host in LINK_KEYS.items():
        url = links.get(key)
        if not url or not isinstance(url, str):
            continue
        if hostname_allowed(url, expected_host):
            return key, url
    return None


def parse_usd_prices(text: str) -> tuple[float | None, float | None]:
    lower = text.lower()
    price_re = re.compile(r"\$\s*(\d+(?:\.\d+)?)")
    in_match = None
    out_match = None
    for m in price_re.finditer(text):
        start = m.start()
        window = lower[max(0, start - 40) : start + 40]
        if any(k in window for k in ("input", "prompt", "in ")):
            in_match = float(m.group(1))
        elif any(k in window for k in ("output", "completion", "out ")):
            out_match = float(m.group(1))
    if in_match is not None and out_match is not None:
        return in_match, out_match
    return None, None


def suggest_capacity(
    link_key: str,
    text: str,
    evidence: str,
) -> tuple[str, str]:
    lower = text.lower()
    ev_lower = evidence.lower()
    if link_key == "cloudflare_pricing":
        if "neuron" in ev_lower or "free tier" in ev_lower or "free tier" in lower:
            return "hobby", "cloudflare_free_tier"
        return "hobby", "cloudflare_free_tier"
    if link_key == "openrouter":
        if PRODUCTION_RPM.search(text) or HIGH_THROUGHPUT.search(text):
            return "industrial", "deepinfra_production_text"
        return "hobby", "openrouter_marketplace"
    if link_key == "provider_pricing":
        if PRODUCTION_RPM.search(text):
            return "industrial", "deepinfra_production_text"
        return "unknown", "google_no_rpm"
    if link_key == "deepinfra":
        if PRODUCTION_RPM.search(text) or HIGH_THROUGHPUT.search(text):
            return "industrial", "deepinfra_production_text"
        if PRODUCTION_WORD.search(text):
            return "industrial", "deepinfra_production_text"
        if "flex" in lower or "priority" in lower:
            return "unknown", "deepinfra_no_rpm"
        return "unknown", "deepinfra_no_rpm"
    if link_key == "xinference":
        return "unknown", "xinference_no_rpm"
    return "unknown", "fetch_error"


def default_fetch(url: str) -> tuple[int, str]:
    req = request.Request(url, headers={"User-Agent": USER_AGENT})

    def do_get() -> tuple[int, str]:
        with request.urlopen(req, timeout=20) as resp:
            status = resp.getcode() or 200
            raw = resp.read(MAX_BODY + 1)
            if len(raw) > MAX_BODY:
                raw = raw[:MAX_BODY]
            charset = resp.headers.get_content_charset() or "utf-8"
            return status, raw.decode(charset, errors="replace")

    try:
        return do_get()
    except error.HTTPError as exc:
        return exc.code, ""
    except (error.URLError, TimeoutError):
        time.sleep(2)
        try:
            return do_get()
        except error.HTTPError as exc:
            return exc.code, ""
        except (error.URLError, TimeoutError):
            return 0, ""


def process_host_row(
    card: dict,
    host: dict,
    *,
    fetch: Callable[[str], tuple[int, str]],
    now: datetime,
) -> dict:
    name = card["name"]
    host_id = host.get("id", "")
    runtime = host.get("runtime")
    runtime_model_id = host.get("runtime_model_id")
    current_capacity = host.get("capacity", "unknown")

    base = {
        "name": name,
        "host_id": host_id,
        "runtime": runtime,
        "runtime_model_id": runtime_model_id,
        "url": "",
        "http_status": 0,
        "fetched_at": now.isoformat(),
        "status": "skipped",
        "evidence": "",
        "price_in": None,
        "price_out": None,
        "suggested_capacity": "unknown",
        "suggestion_reason": "skip_no_link",
        "current_capacity": current_capacity,
        "capacity_applied": False,
    }

    kind = card.get("kind")
    if kind == "service":
        base["suggestion_reason"] = "skip_service"
        return base
    if host_id in SKIP_HOST_IDS or current_capacity == "local":
        base["suggestion_reason"] = "skip_local"
        return base

    picked = pick_allowlisted_url(host.get("links") or {})
    if picked is None:
        base["suggestion_reason"] = "skip_no_link"
        return base

    link_key, url = picked
    base["url"] = url

    status, body = fetch(url)
    base["http_status"] = status
    if status != 200 or not body:
        base["status"] = "error"
        base["suggestion_reason"] = "fetch_error"
        return base

    text = strip_html(body)
    evidence = extract_evidence(text)
    price_in, price_out = parse_usd_prices(text)
    suggested, reason = suggest_capacity(link_key, text, evidence)

    base["status"] = "ok"
    base["evidence"] = evidence
    base["price_in"] = price_in
    base["price_out"] = price_out
    base["suggested_capacity"] = suggested
    base["suggestion_reason"] = reason
    return base


def retrieve(
    catalog: dict,
    *,
    fetch: Callable[[str], tuple[int, str]] | None = None,
    now: datetime | None = None,
    sleep: Callable[[float], None] | None = None,
) -> dict:
    if fetch is None:
        fetch = default_fetch
    if now is None:
        now = datetime.now(timezone.utc)
    if sleep is None:
        sleep = time.sleep

    rows: list[dict] = []
    cards = catalog.get("cards") or []
    for card in cards:
        for host in card.get("hosts") or []:
            row = process_host_row(card, host, fetch=fetch, now=now)
            rows.append(row)
            if row["status"] == "ok" and row["url"]:
                sleep(1)

    return {
        "review_schema": 1,
        "generated_at": now.strftime("%Y-%m-%d"),
        "rows": rows,
    }


def write_review(path: Path, payload: dict) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    tmp = path.with_suffix(path.suffix + ".tmp")
    text = json.dumps(payload, indent=2, ensure_ascii=False) + "\n"
    tmp.write_text(text, encoding="utf-8")
    os.replace(tmp, path)


def main() -> int:
    site_dir = Path(__file__).resolve().parents[1]
    repo_root = site_dir.parent
    catalog_path = repo_root / "roster" / "catalog.json"
    review_path = repo_root / "roster" / "retrieval" / "limits-review.json"

    if not catalog_path.is_file():
        print(f"catalog missing: {catalog_path}", file=sys.stderr)
        print("Run: make -C site generate-roster", file=sys.stderr)
        return 1

    with catalog_path.open(encoding="utf-8") as f:
        catalog = json.load(f)

    payload = retrieve(catalog)
    write_review(review_path, payload)
    print(f"wrote {review_path} ({len(payload['rows'])} rows)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
