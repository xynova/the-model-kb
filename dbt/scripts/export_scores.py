#!/usr/bin/env python3
"""Export analytics.fct_ranked_offers to roster/scores.json (scores_schema: 1)."""

from __future__ import annotations

import json
import os
import sys
from datetime import datetime, timezone
from pathlib import Path


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

        offer: dict = {
            "name": rec["name"],
            "title": rec["title"],
            "kind": rec["kind"],
            "status": rec["status"],
            "jobs": jobs,
            "host_id": rec["host_id"],
            "capacity": rec["capacity"],
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
