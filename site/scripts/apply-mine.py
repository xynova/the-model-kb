#!/usr/bin/env python3
"""Apply limits-review.json to roster cards; capacity only from capacity-decisions.yaml."""

from __future__ import annotations

import importlib.util
import json
import os
import sys
from datetime import datetime, timezone
from pathlib import Path

import yaml

SCRIPTS = Path(__file__).resolve().parent
SPEC = importlib.util.spec_from_file_location(
    "generate_roster", SCRIPTS / "generate-roster.py"
)
assert SPEC and SPEC.loader
gr = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(gr)


def load_decisions(path: Path) -> list[dict]:
    if not path.is_file():
        return []
    with path.open(encoding="utf-8") as f:
        data = yaml.safe_load(f)
    if not data or not isinstance(data, dict):
        raise ValueError(f"{path}: expected mapping with decisions")
    decisions = data.get("decisions")
    if decisions is None:
        raise ValueError(f"{path}: missing decisions key")
    if not isinstance(decisions, list):
        raise ValueError(f"{path}: decisions must be a list")

    seen: set[tuple[str, str]] = set()
    out: list[dict] = []
    for i, row in enumerate(decisions):
        if not isinstance(row, dict):
            raise ValueError(f"{path}: decisions[{i}] must be a map")
        name = row.get("name")
        host_id = row.get("host_id")
        capacity = row.get("capacity")
        if not name or not host_id:
            raise ValueError(f"{path}: decisions[{i}] needs name and host_id")
        if capacity not in gr.ALLOWED_CAPACITY:
            raise ValueError(
                f"{path}: decisions[{i}] invalid capacity {capacity!r}"
            )
        key = (str(name), str(host_id))
        if key in seen:
            raise ValueError(f"{path}: duplicate decision {key}")
        seen.add(key)
        out.append(
            {"name": str(name), "host_id": str(host_id), "capacity": capacity}
        )
    return out


def decision_map(decisions: list[dict]) -> dict[tuple[str, str], str]:
    return {(d["name"], d["host_id"]): d["capacity"] for d in decisions}


def find_card_path(cards_root: Path, name: str) -> Path | None:
    for path in cards_root.rglob("*.md"):
        if path.name.startswith("_"):
            continue
        text = path.read_text(encoding="utf-8")
        fm, _ = gr.split_frontmatter(text)
        if fm and fm.get("name") == name:
            return path
    return None


def apply(
    review: dict,
    decisions: list[dict],
    cards_root: Path,
    *,
    now: datetime | None = None,
) -> tuple[int, list[str]]:
    if now is None:
        now = datetime.now(timezone.utc)
    today = now.strftime("%Y-%m-%d")
    dmap = decision_map(decisions)
    errors: list[str] = []
    rows = review.get("rows") or []

    for row in rows:
        if row.get("status") != "ok":
            continue
        name = row.get("name")
        host_id = row.get("host_id")
        if not name or not host_id:
            errors.append(f"review row missing name or host_id: {row!r}")
            continue

        card_path = find_card_path(cards_root, str(name))
        if card_path is None:
            errors.append(f"no card for name {name!r}")
            continue

        text = card_path.read_text(encoding="utf-8")
        fm, body = gr.split_frontmatter(text)
        if fm is None:
            errors.append(f"{card_path}: invalid frontmatter")
            continue

        hosts = fm.get("hosts") or []
        target = None
        for host in hosts:
            if isinstance(host, dict) and host.get("id") == host_id:
                target = host
                break
        if target is None:
            errors.append(f"{name}: no host id {host_id!r}")
            continue

        evidence = row.get("evidence") or ""
        if evidence:
            target["limits_note"] = evidence

        price_in = row.get("price_in")
        price_out = row.get("price_out")
        if isinstance(price_in, (int, float)) and isinstance(price_out, (int, float)):
            target["price_in"] = float(price_in)
            target["price_out"] = float(price_out)
            if evidence:
                target["price_note"] = evidence

        key = (str(name), str(host_id))
        if key in dmap:
            target["capacity"] = dmap[key]
            row["capacity_applied"] = True

        fm["refreshed"] = today
        sep = "\n\n" if body and not body.startswith("\n") else "\n"
        new_text = (
            "---\n"
            + gr.dump_source_frontmatter(fm)
            + "\n---"
            + sep
            + body
        )
        card_path.write_text(new_text, encoding="utf-8")

    return (1 if errors else 0, errors)


def write_review(path: Path, payload: dict) -> None:
    tmp = path.with_suffix(path.suffix + ".tmp")
    text = json.dumps(payload, indent=2, ensure_ascii=False) + "\n"
    tmp.write_text(text, encoding="utf-8")
    os.replace(tmp, path)


def main() -> int:
    site_dir = SCRIPTS.parent
    repo_root = site_dir.parent
    review_path = repo_root / "roster" / "retrieval" / "limits-review.json"
    decisions_path = repo_root / "roster" / "retrieval" / "capacity-decisions.yaml"
    cards_root = repo_root / "roster"

    if not review_path.is_file():
        print(f"review missing: {review_path}", file=sys.stderr)
        return 1

    try:
        decisions = load_decisions(decisions_path)
    except ValueError as exc:
        print(exc, file=sys.stderr)
        return 1

    with review_path.open(encoding="utf-8") as f:
        review = json.load(f)

    code, errors = apply(review, decisions, cards_root)
    write_review(review_path, review)

    for msg in errors:
        print(msg, file=sys.stderr)
    return code


if __name__ == "__main__":
    raise SystemExit(main())
