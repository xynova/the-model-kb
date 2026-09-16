#!/usr/bin/env python3
"""Copy roster Markdown into site/generated/roster, remapping reserved Hugo frontmatter.

Hugo reserves the frontmatter key `kind` for page kinds. Roster cards use `kind`
for sovereign/provider/service. Rewrite that key to `roster_kind` in a generated
tree so the catalog on disk stays unchanged.
"""

from __future__ import annotations

import shutil
import sys
from pathlib import Path

SKIP_NAMES = {"README.md", "INDEX.md", "taxonomy.md", "_template.md"}


def remap_frontmatter(text: str) -> str:
    if not text.startswith("---"):
        return text
    end = text.find("\n---", 3)
    if end < 0:
        return text
    fm = text[3:end]
    body = text[end:]
    lines = []
    replaced = False
    for line in fm.splitlines(keepends=True):
        if not replaced and line.startswith("kind:"):
            lines.append("roster_kind:" + line[len("kind:") :])
            replaced = True
        else:
            lines.append(line)
    return "---" + "".join(lines) + body


def main() -> int:
    site_dir = Path(__file__).resolve().parents[1]
    src = site_dir.parent / "roster"
    dst = site_dir / "generated" / "roster"
    if not src.is_dir():
        print(f"roster source missing: {src}", file=sys.stderr)
        return 1

    if dst.exists():
        shutil.rmtree(dst)
    dst.mkdir(parents=True)

    for path in src.rglob("*"):
        rel = path.relative_to(src)
        out = dst / rel
        if path.is_dir():
            out.mkdir(parents=True, exist_ok=True)
            continue
        if path.name in SKIP_NAMES:
            continue
        if path.suffix != ".md":
            out.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(path, out)
            continue
        out.parent.mkdir(parents=True, exist_ok=True)
        out.write_text(remap_frontmatter(path.read_text(encoding="utf-8")), encoding="utf-8")

    print(f"generated roster pages under {dst}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
