#!/usr/bin/env python3
"""Copy roster Markdown into site/generated/roster and emit catalog artifacts.

Hugo reserves the frontmatter key `kind` for page kinds. Roster cards use `kind`
for sovereign/provider/service. This script validates cards, rewrites `kind` to
`roster_kind` in the generated tree, writes roster/catalog.json,
site/static/roster.json, and roster/INDEX.md.
"""

from __future__ import annotations

import json
import shutil
import sys
from datetime import datetime, timezone
from pathlib import Path

try:
    import yaml
except ImportError:
    print(
        "PyYAML is required. Install: pip install -r site/requirements-roster.txt",
        file=sys.stderr,
    )
    raise SystemExit(1) from None

# Keep in sync with roster/taxonomy.md.
ALLOWED_KINDS = frozenset({"sovereign", "provider", "service"})
ALLOWED_CATEGORIES = frozenset({"general", "coding", "multimodal", "budget"})
ALLOWED_STATUS = frozenset({"reach-for", "watch", "retired"})
ALLOWED_JOBS = frozenset(
    {
        "chat",
        "reasoning",
        "coding",
        "vision",
        "long-context",
        "cheap",
        "high-throughput",
        "classify",
    }
)
ALLOWED_RUNTIMES = frozenset(
    {
        "cloudflare",
        "deepinfra",
        "openrouter",
        "google",
        "xinference",
        "local",
    }
)
ALLOWED_CAPACITY = frozenset({"hobby", "industrial", "local", "unknown"})
ALLOWED_HOST_IDS = frozenset(
    {
        "cloudflare",
        "deepinfra",
        "openrouter",
        "google",
        "xinference",
        "cursor",
        "local",
    }
)

SKIP_NAMES = {
    "README.md",
    "INDEX.md",
    "taxonomy.md",
    "_template.md",
    "OBJECTIVE.md",
    "catalog.json",
}

CARD_REQUIRED = (
    "name",
    "title",
    "kind",
    "category",
    "jobs",
    "tags",
    "status",
    "lab",
    "source",
    "refreshed",
    "hosts",
)


def split_frontmatter(text: str) -> tuple[dict | None, str]:
    if not text.startswith("---"):
        return None, text
    end = text.find("\n---", 3)
    if end < 0:
        return None, text
    fm_raw = text[3:end]
    body = text[end + 4 :]
    if body.startswith("\n"):
        body = body[1:]
    data = yaml.safe_load(fm_raw)
    if not isinstance(data, dict):
        return None, text
    return data, body


def yaml_unknown_to_json(value):
    if value == "unknown":
        return None
    if isinstance(value, dict):
        return {k: yaml_unknown_to_json(v) for k, v in value.items()}
    if isinstance(value, list):
        return [yaml_unknown_to_json(v) for v in value]
    return value


def validate_card(fm: dict, rel_path: Path) -> None:
    missing = [k for k in CARD_REQUIRED if k not in fm]
    if missing:
        raise ValueError(f"{rel_path}: missing required fields: {missing}")

    kind = fm["kind"]
    if kind not in ALLOWED_KINDS:
        raise ValueError(f"{rel_path}: invalid kind {kind!r}")

    if fm["category"] not in ALLOWED_CATEGORIES:
        raise ValueError(f"{rel_path}: invalid category {fm['category']!r}")

    if fm["status"] not in ALLOWED_STATUS:
        raise ValueError(f"{rel_path}: invalid status {fm['status']!r}")

    jobs = fm["jobs"]
    if not isinstance(jobs, list) or not jobs:
        raise ValueError(f"{rel_path}: jobs must be a non-empty list")
    for job in jobs:
        if job not in ALLOWED_JOBS:
            raise ValueError(f"{rel_path}: invalid job {job!r}")

    if kind == "service" and not fm.get("service"):
        raise ValueError(f"{rel_path}: kind service requires service")

    hosts = fm["hosts"]
    if not isinstance(hosts, list) or not hosts:
        raise ValueError(f"{rel_path}: hosts must be a non-empty list")

    seen_ids: set[str] = set()
    for i, host in enumerate(hosts):
        if not isinstance(host, dict):
            raise ValueError(f"{rel_path}: hosts[{i}] must be a map")
        if "kind" in host:
            raise ValueError(f"{rel_path}: hosts[{i}] must not contain kind")
        hid = host.get("id")
        if hid not in ALLOWED_HOST_IDS:
            raise ValueError(f"{rel_path}: hosts[{i}] invalid id {hid!r}")
        if hid in seen_ids:
            raise ValueError(f"{rel_path}: duplicate host id {hid!r}")
        seen_ids.add(hid)

        cap = host.get("capacity")
        if cap not in ALLOWED_CAPACITY:
            raise ValueError(f"{rel_path}: hosts[{i}] invalid capacity {cap!r}")

        if kind in ("sovereign", "provider"):
            runtime = host.get("runtime")
            if runtime not in ALLOWED_RUNTIMES:
                raise ValueError(
                    f"{rel_path}: hosts[{i}] requires runtime for {kind} cards"
                )

    preferred = fm.get("preferred_host")
    if preferred and preferred not in seen_ids:
        raise ValueError(
            f"{rel_path}: preferred_host {preferred!r} does not match a host id"
        )


def card_to_json(fm: dict) -> dict:
    out: dict = {
        "name": fm["name"],
        "title": fm["title"],
        "kind": fm["kind"],
        "category": fm["category"],
        "jobs": list(fm["jobs"]),
        "status": fm["status"],
        "lab": fm["lab"],
        "hosts": [],
    }
    for key in (
        "preferred_host",
        "service",
        "aa_intelligence",
        "arena_elo",
        "aider_polyglot",
        "context_window",
        "links",
    ):
        if key in fm and fm[key] not in (None, "", "unknown"):
            out[key] = yaml_unknown_to_json(fm[key])

    for host in fm["hosts"]:
        hout: dict = {"id": host["id"], "capacity": host["capacity"]}
        for key in (
            "runtime",
            "runtime_model_id",
            "price_in",
            "price_out",
            "price_note",
            "limits_note",
            "links",
        ):
            if key in host and host[key] not in (None, "", "unknown"):
                hout[key] = yaml_unknown_to_json(host[key])
        out["hosts"].append(hout)
    return out


def dump_hugo_frontmatter(fm: dict) -> str:
    hugo_fm = dict(fm)
    kind = hugo_fm.pop("kind", None)
    if kind is not None:
        hugo_fm["roster_kind"] = kind
    return yaml.dump(
        hugo_fm,
        default_flow_style=False,
        sort_keys=False,
        allow_unicode=True,
    ).rstrip()


def dump_source_frontmatter(fm: dict) -> str:
    return yaml.dump(
        fm,
        default_flow_style=False,
        sort_keys=False,
        allow_unicode=True,
    ).rstrip()


def write_index(cards: list[tuple[dict, Path]], index_path: Path) -> None:
    lines = [
        "# Roster index",
        "",
        "Regenerated by `site/scripts/generate-roster.py`.",
        "",
        "| name | title | kind | service | category | status | lab | hosts | jobs | tags | path |",
        "|------|-------|------|---------|----------|--------|-----|-------|------|------|------|",
    ]
    for fm, rel in sorted(cards, key=lambda x: x[0]["name"]):
        host_ids = ",".join(h["id"] for h in fm["hosts"])
        jobs = ",".join(fm["jobs"])
        tags = ",".join(fm["tags"]) if isinstance(fm["tags"], list) else ""
        service = fm.get("service") or ""
        lines.append(
            f"| {fm['name']} | {fm['title']} | {fm['kind']} | {service} | "
            f"{fm['category']} | {fm['status']} | {fm['lab']} | {host_ids} | "
            f"{jobs} | {tags} | {rel.as_posix()} |"
        )
    lines.append("")
    index_path.write_text("\n".join(lines), encoding="utf-8")


def main() -> int:
    site_dir = Path(__file__).resolve().parents[1]
    repo_root = site_dir.parent
    src = repo_root / "roster"
    dst = site_dir / "generated" / "roster"
    catalog_path = repo_root / "roster" / "catalog.json"
    static_json = site_dir / "static" / "roster.json"
    index_path = repo_root / "roster" / "INDEX.md"

    if not src.is_dir():
        print(f"roster source missing: {src}", file=sys.stderr)
        return 1

    cards: list[tuple[dict, str, str]] = []
    names_seen: set[str] = set()

    for path in sorted(src.rglob("*.md")):
        rel = path.relative_to(src)
        if path.name in SKIP_NAMES:
            continue
        if rel.parts and rel.parts[0].startswith("."):
            continue

        text = path.read_text(encoding="utf-8")
        fm, body = split_frontmatter(text)
        if fm is None:
            print(f"skip (no frontmatter): {rel}", file=sys.stderr)
            continue

        validate_card(fm, rel)

        name = fm["name"]
        if name in names_seen:
            raise SystemExit(f"duplicate card name: {name}")
        names_seen.add(name)

        cards.append((fm, body, str(rel)))

    if dst.exists():
        shutil.rmtree(dst)
    dst.mkdir(parents=True)

    json_cards = []
    index_rows: list[tuple[dict, Path]] = []

    for fm, body, rel_str in cards:
        rel = Path(rel_str)
        out = dst / rel
        out.parent.mkdir(parents=True, exist_ok=True)
        hugo_text = "---\n" + dump_hugo_frontmatter(fm) + "\n---\n\n" + body
        out.write_text(hugo_text, encoding="utf-8")
        json_cards.append(card_to_json(fm))
        index_rows.append((fm, rel))

    generated_at = datetime.now(timezone.utc).strftime("%Y-%m-%d")
    catalog = {
        "catalog_schema": 1,
        "generated_at": generated_at,
        "cards": sorted(json_cards, key=lambda c: c["name"]),
    }
    catalog_bytes = json.dumps(catalog, indent=2, ensure_ascii=False) + "\n"
    catalog_path.write_text(catalog_bytes, encoding="utf-8")
    static_json.parent.mkdir(parents=True, exist_ok=True)
    static_json.write_text(catalog_bytes, encoding="utf-8")
    write_index(index_rows, index_path)

    print(f"generated roster pages under {dst}")
    print(f"wrote {catalog_path} and {static_json}")
    print(f"wrote {index_path} ({len(json_cards)} cards)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
