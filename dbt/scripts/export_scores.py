#!/usr/bin/env python3
"""Export analytics.fct_ranked_offers to roster/scores.json (scores_schema: 1)."""

from __future__ import annotations

import json
import os
import sys
from datetime import datetime, timezone
from pathlib import Path

# Keep in sync with roster/taxonomy.md Quality leaves and tools/roster-pick/internal/scores/leaves.go
QUALITY_LEAF_IDS = [
    "quality.general_problem_solving.benchmark_composite",
    "quality.human_preference.pairwise_rating",
    "quality.coding.agent_benchmark_composite",
    "quality.coding.code_edit_first_attempt",
    "quality.multimodal_understanding.benchmark_accuracy",
    "quality.typed_decision.classification_accuracy",
    "quality.typed_decision.precision",
    "quality.typed_decision.recall",
    "quality.typed_decision.brier_loss",
    "quality.typed_decision.ordinal_mae",
    "quality.speech.human_preference",
]
PERFORMANCE_LEAF_IDS = [
    "performance.end_to_end_latency_ms",
    "performance.time_to_first_token_ms",
    "performance.request_success_rate",
]


def _parse_json_map(raw: object) -> dict:
    if raw is None:
        return {}
    if isinstance(raw, str):
        return json.loads(raw)
    if isinstance(raw, dict):
        return raw
    raise TypeError(f"expected JSON map, got {type(raw)!r}")


def _normalize_leaf(raw: object) -> dict:
    leaf = _parse_json_map(raw) if not isinstance(raw, dict) else raw
    out: dict = {"status": leaf["status"]}
    if leaf.get("status") == "measured":
        out["value"] = leaf["value"]
        out["attribution"] = leaf.get("attribution", "model_only")
        if "comparison_key" in leaf:
            out["comparison_key"] = leaf["comparison_key"]
    return out


def _leaf_map(raw: object, registered: list[str], label: str) -> dict:
    parsed = _parse_json_map(raw)
    keys = set(parsed)
    expected = set(registered)
    if keys != expected:
        missing = expected - keys
        extra = keys - expected
        raise ValueError(
            f"{label} keys mismatch: missing={sorted(missing)} extra={sorted(extra)}"
        )
    return {leaf_id: _normalize_leaf(parsed[leaf_id]) for leaf_id in registered}


def _string_list(raw: object) -> list[str]:
    if raw is None:
        return []
    if isinstance(raw, str):
        raw = json.loads(raw)
    return list(raw)


def main() -> int:
    dbt_dir = Path(__file__).resolve().parent.parent
    repo_root = dbt_dir.parent
    duckdb_path = Path(
        os.environ.get("ROSTER_DUCKDB", str(repo_root / "artifacts" / "roster_analytics.duckdb"))
    )
    catalog_path = Path(
        os.environ.get("ROSTER_CATALOG", str(repo_root / "roster" / "catalog.json"))
    )
    scores_path = repo_root / "roster" / "scores.json"
    scores_tmp = repo_root / "roster" / "scores.json.tmp"

    if not catalog_path.is_file():
        print(f"catalog missing: {catalog_path}", file=sys.stderr)
        return 1

    with catalog_path.open(encoding="utf-8") as f:
        catalog = json.load(f)
    expected_hosts = sum(len(c["hosts"]) for c in catalog["cards"])

    try:
        import duckdb
    except ImportError as exc:
        print(f"duckdb import failed: {exc}", file=sys.stderr)
        print("Run: python3 -m venv .venv && .venv/bin/pip install -r requirements.txt", file=sys.stderr)
        return 1

    if not duckdb_path.is_file():
        print(f"duckdb file missing: {duckdb_path}", file=sys.stderr)
        print("Run: make -C dbt build (dbt step must succeed first)", file=sys.stderr)
        return 1

    try:
        con = duckdb.connect(str(duckdb_path), read_only=True)
    except Exception as exc:
        print(f"duckdb connect failed: {exc}", file=sys.stderr)
        print(f"duckdb python version: {duckdb.__version__}", file=sys.stderr)
        return 1

    rows = con.execute(
        "select * from analytics.fct_ranked_offers order by name, host_id"
    ).fetchall()
    cols = [d[0] for d in con.description]
    con.close()

    offers: list[dict] = []
    for row in rows:
        rec = dict(zip(cols, row, strict=True))
        jobs = rec.pop("jobs")
        if isinstance(jobs, str):
            jobs = json.loads(jobs)
        elif jobs is None:
            jobs = []
        else:
            jobs = list(jobs)

        quality_raw = rec.pop("quality")
        performance_raw = rec.pop("performance")
        capability_cohort = rec.pop("capability_cohort")
        input_modalities = _string_list(rec.pop("input_modalities"))
        output_contracts = _string_list(rec.pop("output_contracts"))
        execution_modes = _string_list(rec.pop("execution_modes"))

        try:
            quality = _leaf_map(quality_raw, QUALITY_LEAF_IDS, "quality")
            performance = _leaf_map(performance_raw, PERFORMANCE_LEAF_IDS, "performance")
        except (TypeError, ValueError, KeyError, json.JSONDecodeError) as exc:
            print(
                f"leaf map failed for {rec.get('name')}|{rec.get('host_id')}: {exc}",
                file=sys.stderr,
            )
            return 1

        leaf_cohorts = sorted(
            leaf_id for leaf_id, leaf in quality.items() if leaf.get("status") == "measured"
        )

        offer: dict = {
            "name": rec["name"],
            "title": rec["title"],
            "kind": rec["kind"],
            "status": rec["status"],
            "jobs": jobs,
            "host_id": rec["host_id"],
            "capacity": rec["capacity"],
            "capability_cohort": capability_cohort,
            "input_modalities": input_modalities,
            "output_contracts": output_contracts,
            "execution_modes": execution_modes,
            "leaf_cohorts": leaf_cohorts,
            "quality": quality,
            "performance": performance,
        }
        if rec.get("runtime"):
            offer["runtime"] = rec["runtime"]
        if rec.get("runtime_model_id"):
            offer["runtime_model_id"] = rec["runtime_model_id"]
        if rec.get("price_in") is not None:
            offer["price_in"] = rec["price_in"]
        if rec.get("price_out") is not None:
            offer["price_out"] = rec["price_out"]
        if rec.get("cost_per_m") is not None:
            offer["cost_per_m"] = rec["cost_per_m"]
        if rec.get("rank_by_cost") is not None:
            offer["rank_by_cost"] = int(rec["rank_by_cost"])
        offer["industrial_ok"] = bool(rec.get("industrial_ok"))
        offers.append(offer)

    if len(offers) != expected_hosts:
        print(
            f"offer count mismatch: scores would have {len(offers)}, catalog has {expected_hosts} host rows",
            file=sys.stderr,
        )
        return 1

    generated_at = datetime.now(timezone.utc).strftime("%Y-%m-%d")
    payload = {
        "scores_schema": 1,
        "generated_at": generated_at,
        "source_catalog_schema": catalog.get("catalog_schema", 1),
        "offers": offers,
    }

    scores_tmp.parent.mkdir(parents=True, exist_ok=True)
    text = json.dumps(payload, indent=2, ensure_ascii=False) + "\n"
    scores_tmp.write_text(text, encoding="utf-8")
    os.replace(scores_tmp, scores_path)
    print(f"wrote {scores_path} ({len(offers)} offers)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
