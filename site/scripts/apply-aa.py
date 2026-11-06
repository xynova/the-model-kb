#!/usr/bin/env python3
"""Apply aa-intelligence-review.json scores onto roster card frontmatter."""

from __future__ import annotations

import importlib.util
import json
import sys
from datetime import datetime, timezone
from pathlib import Path

SCRIPTS = Path(__file__).resolve().parent
APPLY_SPEC = importlib.util.spec_from_file_location(
    "apply_mine", SCRIPTS / "apply-mine.py"
)
assert APPLY_SPEC and APPLY_SPEC.loader
am = importlib.util.module_from_spec(APPLY_SPEC)
APPLY_SPEC.loader.exec_module(am)

GR_SPEC = importlib.util.spec_from_file_location(
    "generate_roster", SCRIPTS / "generate-roster.py"
)
assert GR_SPEC and GR_SPEC.loader
gr = importlib.util.module_from_spec(GR_SPEC)
GR_SPEC.loader.exec_module(gr)

AA_SPEC = importlib.util.spec_from_file_location(
    "retrieve_aa", SCRIPTS / "retrieve-aa.py"
)
assert AA_SPEC and AA_SPEC.loader
ra = importlib.util.module_from_spec(AA_SPEC)
AA_SPEC.loader.exec_module(ra)


CLEAR_STATUSES = frozenset({"no_score", "unmatched"})


def same_score(left: object, right: object) -> bool:
    a = ra.yaml_score(left)
    b = ra.yaml_score(right)
    if a is None or b is None:
        return False
    return float(a) == float(b)


def has_filed_score(value: object) -> bool:
    return ra.yaml_score(value) is not None


def model_page_url(slug: str) -> str:
    return f"https://artificialanalysis.ai/models/{slug}"


def write_card(card_path: Path, fm: dict, body: str) -> None:
    sep = "\n\n" if body and not body.startswith("\n") else "\n"
    card_path.write_text(
        "---\n" + gr.dump_source_frontmatter(fm) + "\n---" + sep + body,
        encoding="utf-8",
    )


def load_card(
    cards_root: Path, name: str, errors: list[str]
) -> tuple[Path, dict, str] | None:
    card_path = am.find_card_path(cards_root, name)
    if card_path is None:
        errors.append(f"no card for name {name!r}")
        return None
    text = card_path.read_text(encoding="utf-8")
    fm, body = gr.split_frontmatter(text)
    if fm is None:
        errors.append(f"{card_path}: invalid frontmatter")
        return None
    return card_path, fm, body


def apply(
    review: dict,
    cards_root: Path,
    *,
    now: datetime | None = None,
) -> tuple[int, list[str]]:
    if now is None:
        now = datetime.now(timezone.utc)
    if review.get("status") != "ok":
        return 0, []
    today = now.strftime("%Y-%m-%d")
    errors: list[str] = []
    for row in review.get("rows") or []:
        name = row.get("name")
        if not name:
            continue
        status = row.get("status")
        if status == "ok":
            score = ra.yaml_score(row.get("aa_intelligence"))
            slug = row.get("aa_slug") or ""
            if score is None:
                errors.append(f"review row missing name or score: {row!r}")
                continue
            loaded = load_card(cards_root, str(name), errors)
            if loaded is None:
                continue
            card_path, fm, body = loaded
            previous = fm.get("aa_intelligence")
            row["previous"] = previous
            if same_score(previous, score):
                row["applied"] = False
                continue
            fm["aa_intelligence"] = score
            links = fm.get("links")
            if not isinstance(links, dict):
                links = {}
                fm["links"] = links
            if slug:
                links["artificial_analysis"] = model_page_url(str(slug))
            fm["refreshed"] = today
            write_card(card_path, fm, body)
            row["applied"] = True
            continue
        if status not in CLEAR_STATUSES:
            continue
        loaded = load_card(cards_root, str(name), errors)
        if loaded is None:
            continue
        card_path, fm, body = loaded
        previous = fm.get("aa_intelligence")
        row["previous"] = previous
        if not has_filed_score(previous):
            row["applied"] = False
            continue
        fm["aa_intelligence"] = "unknown"
        fm["refreshed"] = today
        write_card(card_path, fm, body)
        row["applied"] = True
    return (1 if errors else 0, errors)


def main() -> int:
    site_dir = SCRIPTS.parent
    repo_root = site_dir.parent
    review_path = repo_root / "roster" / "retrieval" / "aa-intelligence-review.json"
    cards_root = repo_root / "roster"
    if not review_path.is_file():
        print(f"AA review missing: {review_path}", file=sys.stderr)
        return 0
    with review_path.open(encoding="utf-8") as f:
        review = json.load(f)
    code, errors = apply(review, cards_root)
    am.write_review(review_path, review)
    for msg in errors:
        print(msg, file=sys.stderr)
    return code


if __name__ == "__main__":
    raise SystemExit(main())
