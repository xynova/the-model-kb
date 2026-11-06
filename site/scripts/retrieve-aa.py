#!/usr/bin/env python3
"""Fetch Artificial Analysis Intelligence Index scores into aa-intelligence-review.json."""

from __future__ import annotations

import importlib.util
import json
import os
import re
import sys
import time
from datetime import datetime, timezone
from pathlib import Path
from typing import Callable
from urllib import error, parse, request

_OC_SPEC = importlib.util.spec_from_file_location(
    "operator_config", Path(__file__).resolve().parent / "operator_config.py"
)
assert _OC_SPEC and _OC_SPEC.loader
oc = importlib.util.module_from_spec(_OC_SPEC)
_OC_SPEC.loader.exec_module(oc)

AA_FREE_URL = "https://artificialanalysis.ai/api/v2/language/models/free"
USER_AGENT = "the-model-kb-roster-mine/1.0"
MAX_BODY = 2 * 1024 * 1024
MAX_PAGES = 50
SUPPORTED_INDEX_MAJOR = 5

EFFORT_PAREN = re.compile(
    r"[\s_\-]*\((max|xhigh|high|medium|low|reasoning|non-reasoning)\)\s*$",
    re.IGNORECASE,
)
EFFORT_SLUG = re.compile(
    r"-(max|xhigh|high|medium|low|reasoning|non-reasoning)$",
    re.IGNORECASE,
)
MODEL_PAGE = re.compile(
    r"https?://(?:www\.)?artificialanalysis\.ai/models/(?:releases/)?([a-z0-9\-]+)",
    re.IGNORECASE,
)
NON_ALNUM = re.compile(r"[^a-z0-9]+")

EFFORT_RANK = {
    "max": 50,
    "xhigh": 40,
    "high": 30,
    "reasoning": 25,
    "": 20,
    "medium": 10,
    "low": 0,
}
SKIP_EFFORTS = frozenset({"non-reasoning"})


def normalize_key(value: str) -> str:
    text = value.lower().replace(".", "-")
    return NON_ALNUM.sub("-", text).strip("-")


def parse_effort(name: str, slug: str) -> tuple[str, str]:
    effort = ""
    label = name.strip()
    match = EFFORT_PAREN.search(label)
    if match:
        effort = match.group(1).lower()
        label = label[: match.start()]
    slug_text = slug.strip()
    slug_match = EFFORT_SLUG.search(slug_text)
    if slug_match:
        if not effort:
            effort = slug_match.group(1).lower()
        slug_text = slug_text[: slug_match.start()]
    key = normalize_key(slug_text) or normalize_key(label)
    return key, effort


def effort_rank(effort: str) -> int:
    return EFFORT_RANK.get(effort, 15)


def yaml_score(value: object) -> int | float | None:
    if value is None:
        return None
    if isinstance(value, bool):
        return None
    if isinstance(value, int):
        return value
    if isinstance(value, float):
        if value.is_integer():
            return int(value)
        return value
    return None


def yaml_cost_per_task(value: object) -> float | None:
    if isinstance(value, dict):
        task = value.get("cost_per_task")
        if isinstance(task, dict):
            score = yaml_score(task.get("total_cost"))
            if score is not None:
                return float(score)
        score = yaml_score(value.get("total_cost"))
        if score is not None:
            return float(score)
        return None
    score = yaml_score(value)
    if score is None:
        return None
    return float(score)


def creator_name(model: dict) -> str:
    creator = model.get("model_creator")
    if isinstance(creator, dict):
        name = creator.get("name")
        if isinstance(name, str):
            return name.strip()
    return ""


def model_page_url(slug: str) -> str:
    return f"https://artificialanalysis.ai/models/{slug}"


EVALUATION_KEYS = (
    "artificial_analysis_intelligence_index",
    "artificial_analysis_coding_index",
    "artificial_analysis_agentic_index",
    "artificial_analysis_finance_and_accounting_index",
    "artificial_analysis_strategy_and_ops_index",
    "artificial_analysis_legal_index",
    "artificial_analysis_healthcare_and_medical_index",
    "artificial_analysis_engineering_index",
    "artificial_analysis_economics_index",
)
PERFORMANCE_KEYS = (
    "median_output_tokens_per_second",
    "median_time_to_first_token_seconds",
    "median_time_to_first_answer_token_seconds",
    "median_end_to_end_response_time_seconds",
)
PRICING_KEYS = (
    "price_1m_input_tokens",
    "price_1m_output_tokens",
)


def extract_overview(model: dict) -> dict:
    evaluations = model.get("evaluations") if isinstance(model.get("evaluations"), dict) else {}
    pricing = model.get("pricing") if isinstance(model.get("pricing"), dict) else {}
    performance = model.get("performance") if isinstance(model.get("performance"), dict) else {}
    return {
        "evaluations": {key: yaml_score(evaluations.get(key)) for key in EVALUATION_KEYS},
        "pricing": {key: yaml_score(pricing.get(key)) for key in PRICING_KEYS},
        "performance": {key: yaml_score(performance.get(key)) for key in PERFORMANCE_KEYS},
    }


def intelligence_per_cost(score: float | int | None, cost: float | None) -> float | None:
    if score is None or cost is None or cost <= 0:
        return None
    return float(score) / float(cost)


def compare_peer(candidate: dict) -> dict:
    return {
        "name": candidate["aa_name"],
        "creator": candidate["creator"],
        "slug": candidate["aa_slug"],
        "effort": candidate["effort"],
        "aa_intelligence": candidate["aa_intelligence"],
        "cost_per_task": candidate["cost_per_task"],
        "url": model_page_url(str(candidate["aa_slug"])),
    }


def thinking_candidates(models: list[dict]) -> list[dict]:
    out: list[dict] = []
    for model in models:
        slug = str(model.get("slug") or "")
        name = str(model.get("name") or "")
        key, effort = parse_effort(name, slug)
        if not key or effort in SKIP_EFFORTS:
            continue
        evaluations = model.get("evaluations") or {}
        score = yaml_score(
            evaluations.get("artificial_analysis_intelligence_index")
            if isinstance(evaluations, dict)
            else None
        )
        if score is None:
            continue
        out.append(
            {
                "key": key,
                "aa_slug": slug,
                "aa_name": name,
                "effort": effort,
                "creator": creator_name(model),
                "aa_intelligence": score,
                "cost_per_task": yaml_cost_per_task(
                    model.get("artificial_analysis_intelligence_index_cost")
                ),
                "rank": effort_rank(effort),
            }
        )
    return out


def pick_alternative(
    current: dict, pool: list[dict], roster_keys: set[str]
) -> dict | None:
    current_score = yaml_score(current.get("aa_intelligence"))
    current_cost = current.get("cost_per_task")
    if current_score is None or not isinstance(current_cost, (int, float)):
        return None
    current_cost = float(current_cost)
    if current_cost <= 0:
        return None
    best: dict | None = None
    best_ipc: float | None = None
    for cand in pool:
        if cand["key"] in roster_keys or cand["key"] == current.get("key"):
            continue
        score = yaml_score(cand.get("aa_intelligence"))
        cost = cand.get("cost_per_task")
        ipc = intelligence_per_cost(score, cost)
        if score is None or ipc is None or not isinstance(cost, (int, float)):
            continue
        if score < current_score or float(cost) >= current_cost:
            continue
        if best is None or best_ipc is None or ipc > best_ipc:
            best = cand
            best_ipc = ipc
            continue
        if ipc == best_ipc and float(score) > float(best["aa_intelligence"]):
            best = cand
    if best is None:
        return None
    return compare_peer(best)


def pick_similar(
    current: dict,
    pool: list[dict],
    roster_keys: set[str],
    *,
    skip_slug: str = "",
) -> dict | None:
    current_score = yaml_score(current.get("aa_intelligence"))
    if current_score is None:
        return None
    best: dict | None = None
    best_delta: float | None = None
    for cand in pool:
        if cand["key"] in roster_keys or cand["key"] == current.get("key"):
            continue
        if skip_slug and cand["aa_slug"] == skip_slug:
            continue
        score = yaml_score(cand.get("aa_intelligence"))
        if score is None:
            continue
        delta = abs(float(score) - float(current_score))
        if best is None or best_delta is None or delta < best_delta:
            best = cand
            best_delta = delta
            continue
        if delta == best_delta and float(score) > float(best["aa_intelligence"]):
            best = cand
    if best is None:
        return None
    return compare_peer(best)


def link_slug(card: dict) -> str:
    links = card.get("links") or {}
    if not isinstance(links, dict):
        return ""
    url = links.get("artificial_analysis")
    if not isinstance(url, str):
        return ""
    match = MODEL_PAGE.search(url)
    if not match:
        return ""
    key, _ = parse_effort("", match.group(1))
    return key


def pick_best(models: list[dict]) -> dict[str, dict]:
    best: dict[str, dict] = {}
    for model in models:
        slug = str(model.get("slug") or "")
        name = str(model.get("name") or "")
        key, effort = parse_effort(name, slug)
        if not key or effort in SKIP_EFFORTS:
            continue
        evaluations = model.get("evaluations") or {}
        score = yaml_score(
            evaluations.get("artificial_analysis_intelligence_index")
            if isinstance(evaluations, dict)
            else None
        )
        candidate = {
            "key": key,
            "aa_slug": slug,
            "aa_name": name,
            "effort": effort,
            "creator": creator_name(model),
            "aa_intelligence": score,
            "cost_per_task": yaml_cost_per_task(
                model.get("artificial_analysis_intelligence_index_cost")
            ),
            "rank": effort_rank(effort),
            "overview": extract_overview(model),
        }
        current = best.get(key)
        if current is None:
            best[key] = candidate
            continue
        if candidate["rank"] > current["rank"]:
            best[key] = candidate
            continue
        if (
            candidate["rank"] == current["rank"]
            and current["aa_intelligence"] is None
            and candidate["aa_intelligence"] is not None
        ):
            best[key] = candidate
    return best


def card_keys(card: dict) -> list[str]:
    keys: list[str] = []
    name = card.get("name")
    if isinstance(name, str) and name:
        keys.append(normalize_key(name))
    title = card.get("title")
    if isinstance(title, str) and title:
        keys.append(normalize_key(title))
    linked = link_slug(card)
    if linked:
        keys.append(linked)
    seen: set[str] = set()
    out: list[str] = []
    for key in keys:
        if key and key not in seen:
            seen.add(key)
            out.append(key)
    return out


def match_cards(catalog: dict, aa_models: list[dict]) -> list[dict]:
    best = pick_best(aa_models)
    pool = thinking_candidates(aa_models)
    roster_keys: set[str] = set()
    for card in catalog.get("cards") or []:
        if card.get("kind") == "service":
            continue
        roster_keys.update(card_keys(card))
    rows: list[dict] = []
    for card in catalog.get("cards") or []:
        name = card.get("name")
        if not name:
            continue
        row = {
            "name": name,
            "status": "unmatched",
            "aa_slug": "",
            "aa_name": "",
            "effort": "",
            "aa_creator": "",
            "aa_cost_per_task": None,
            "aa_intelligence": None,
            "alternative": None,
            "similar": None,
            "overview": None,
            "previous": card.get("aa_intelligence"),
            "applied": False,
        }
        if card.get("kind") == "service":
            row["status"] = "skipped_service"
            rows.append(row)
            continue
        hit = None
        for key in card_keys(card):
            if key in best:
                hit = best[key]
                break
        if hit is None:
            rows.append(row)
            continue
        row["aa_slug"] = hit["aa_slug"]
        row["aa_name"] = hit["aa_name"]
        row["effort"] = hit["effort"]
        row["aa_creator"] = hit.get("creator") or ""
        row["aa_cost_per_task"] = hit.get("cost_per_task")
        row["aa_intelligence"] = hit["aa_intelligence"]
        row["overview"] = hit.get("overview")
        if hit["aa_intelligence"] is None:
            row["status"] = "no_score"
        else:
            row["status"] = "ok"
            alt = pick_alternative(hit, pool, roster_keys)
            row["alternative"] = alt
            skip_slug = str(alt["slug"]) if alt else ""
            row["similar"] = pick_similar(
                hit, pool, roster_keys, skip_slug=skip_slug
            )
        rows.append(row)
    return rows


def default_fetch_page(page: int, api_key: str) -> dict:
    url = f"{AA_FREE_URL}?{parse.urlencode({'page': page})}"
    req = request.Request(
        url,
        headers={
            "User-Agent": USER_AGENT,
            "Accept": "application/json",
            "x-api-key": api_key,
        },
    )
    with request.urlopen(req, timeout=30) as resp:
        raw = resp.read(MAX_BODY + 1)
        if len(raw) > MAX_BODY:
            raw = raw[:MAX_BODY]
        charset = resp.headers.get_content_charset() or "utf-8"
        payload = json.loads(raw.decode(charset, errors="replace"))
    if not isinstance(payload, dict):
        raise ValueError("AA response is not an object")
    return payload


def collect_models(
    api_key: str,
    *,
    fetch_page: Callable[[int, str], dict] | None = None,
    sleep: Callable[[float], None] | None = None,
) -> tuple[float | int | None, list[dict]]:
    if fetch_page is None:
        fetch_page = default_fetch_page
    if sleep is None:
        sleep = time.sleep
    models: list[dict] = []
    version: float | int | None = None
    page = 1
    while page <= MAX_PAGES:
        payload = fetch_page(page, api_key)
        if version is None:
            version = payload.get("intelligence_index_version")
        chunk = payload.get("data") or []
        if isinstance(chunk, list):
            models.extend(item for item in chunk if isinstance(item, dict))
        pagination = payload.get("pagination") or {}
        if not pagination.get("has_more"):
            break
        page += 1
        sleep(0.2)
    return version, models


def retrieve(
    catalog: dict,
    *,
    api_key: str | None,
    fetch_page: Callable[[int, str], dict] | None = None,
    now: datetime | None = None,
    sleep: Callable[[float], None] | None = None,
) -> dict:
    if now is None:
        now = datetime.now(timezone.utc)
    generated = now.strftime("%Y-%m-%d")
    if not api_key:
        return {
            "review_schema": 1,
            "generated_at": generated,
            "source": AA_FREE_URL,
            "status": "skipped_no_key",
            "intelligence_index_version": None,
            "rows": [],
        }
    try:
        version, models = collect_models(
            api_key, fetch_page=fetch_page, sleep=sleep
        )
    except (error.URLError, error.HTTPError, TimeoutError, ValueError, json.JSONDecodeError) as exc:
        return {
            "review_schema": 1,
            "generated_at": generated,
            "source": AA_FREE_URL,
            "status": "error",
            "error": str(exc),
            "intelligence_index_version": None,
            "rows": [],
        }
    major = int(version) if isinstance(version, (int, float)) else None
    if major != SUPPORTED_INDEX_MAJOR:
        return {
            "review_schema": 1,
            "generated_at": generated,
            "source": AA_FREE_URL,
            "status": "skipped_version",
            "intelligence_index_version": version,
            "rows": [],
        }
    return {
        "review_schema": 1,
        "generated_at": generated,
        "source": AA_FREE_URL,
        "status": "ok",
        "intelligence_index_version": version,
        "rows": match_cards(catalog, models),
    }


def write_review(path: Path, payload: dict) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    tmp = path.with_suffix(path.suffix + ".tmp")
    text = json.dumps(payload, indent=2, ensure_ascii=False) + "\n"
    tmp.write_text(text, encoding="utf-8")
    os.replace(tmp, path)


def compact_num(value: float) -> int | float:
    if float(value).is_integer():
        return int(value)
    return round(float(value), 4)


def metric_summary(values: list[float]) -> dict | None:
    nums = sorted(float(v) for v in values)
    if not nums:
        return None
    n = len(nums)
    mid = n // 2
    if n % 2 == 1:
        median = nums[mid]
    else:
        median = (nums[mid - 1] + nums[mid]) / 2
    return {
        "min": compact_num(nums[0]),
        "max": compact_num(nums[-1]),
        "median": compact_num(median),
        "mean": compact_num(sum(nums) / n),
        "n": n,
    }


def overview_stats(cards: dict[str, dict]) -> dict:
    eval_values: dict[str, list[float]] = {key: [] for key in EVALUATION_KEYS}
    price_values: dict[str, list[float]] = {key: [] for key in PRICING_KEYS}
    price_values["cost_per_task"] = []
    perf_values: dict[str, list[float]] = {key: [] for key in PERFORMANCE_KEYS}
    for card in cards.values():
        overview = card.get("overview")
        if not isinstance(overview, dict):
            overview = {}
        evaluations = overview.get("evaluations")
        if isinstance(evaluations, dict):
            for key in EVALUATION_KEYS:
                score = yaml_score(evaluations.get(key))
                if score is not None:
                    eval_values[key].append(float(score))
        pricing = overview.get("pricing")
        if isinstance(pricing, dict):
            for key in PRICING_KEYS:
                score = yaml_score(pricing.get(key))
                if score is not None:
                    price_values[key].append(float(score))
        cost = yaml_score(card.get("cost_per_task"))
        if cost is not None:
            price_values["cost_per_task"].append(float(cost))
        performance = overview.get("performance")
        if isinstance(performance, dict):
            for key in PERFORMANCE_KEYS:
                score = yaml_score(performance.get(key))
                if score is not None:
                    perf_values[key].append(float(score))
    return {
        "evaluations": {
            key: metric_summary(vals) for key, vals in eval_values.items() if vals
        },
        "pricing": {
            key: metric_summary(vals) for key, vals in price_values.items() if vals
        },
        "performance": {
            key: metric_summary(vals) for key, vals in perf_values.items() if vals
        },
    }


def card_compare_entry(row: dict) -> dict:
    slug = row.get("aa_slug") or ""
    return {
        "aa_intelligence": row.get("aa_intelligence"),
        "cost_per_task": row.get("aa_cost_per_task"),
        "creator": row.get("aa_creator") or "",
        "effort": row.get("effort") or "",
        "url": model_page_url(str(slug)) if slug else "",
        "alternative": row.get("alternative"),
        "similar": row.get("similar"),
        "overview": row.get("overview"),
    }


def compare_data(review: dict) -> dict:
    cards: dict[str, dict] = {}
    for row in review.get("rows") or []:
        name = row.get("name")
        if not isinstance(name, str) or not name:
            continue
        cards[name] = card_compare_entry(row)
    return {
        "generated_at": review.get("generated_at"),
        "cards": cards,
        "stats": overview_stats(cards),
    }


def write_compare_data(path: Path, review: dict) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    tmp = path.with_suffix(path.suffix + ".tmp")
    text = json.dumps(compare_data(review), indent=2, ensure_ascii=False) + "\n"
    tmp.write_text(text, encoding="utf-8")
    os.replace(tmp, path)


def main() -> int:
    site_dir = Path(__file__).resolve().parents[1]
    repo_root = site_dir.parent
    catalog_path = repo_root / "roster" / "catalog.json"
    review_path = repo_root / "roster" / "retrieval" / "aa-intelligence-review.json"
    oc.load_dotenv_files(repo_root)
    try:
        api_key = oc.resolve_aa_api_key(repo_root=repo_root)
    except ValueError as exc:
        print(str(exc), file=sys.stderr)
        return 1

    if not api_key:
        print("AA_API_KEY unset; skip AA retrieve", file=sys.stderr)
        return 0

    if not catalog_path.is_file():
        print(f"catalog missing: {catalog_path}", file=sys.stderr)
        print("Run: make -C site generate-roster", file=sys.stderr)
        return 1

    with catalog_path.open(encoding="utf-8") as f:
        catalog = json.load(f)

    payload = retrieve(catalog, api_key=api_key)
    if payload["status"] == "error":
        print(payload.get("error") or "AA retrieve failed", file=sys.stderr)
        return 1
    if payload["status"] == "skipped_version":
        print(
            "AA intelligence_index_version "
            f"{payload.get('intelligence_index_version')!r} is not {SUPPORTED_INDEX_MAJOR}",
            file=sys.stderr,
        )
        write_review(review_path, payload)
        return 1

    write_review(review_path, payload)
    write_compare_data(site_dir / "data" / "aa_compare.json", payload)
    print(f"wrote {review_path} ({len(payload['rows'])} rows)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
