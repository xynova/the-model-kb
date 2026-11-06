#!/usr/bin/env python3
"""Operator config discovery and AA_API_KEY secret resolve for roster mine."""

from __future__ import annotations

import argparse
import os
import stat
import subprocess
import sys
from pathlib import Path
from typing import Protocol

import yaml

APP_DIR = "the-model-kb"
OVERRIDE_ENV = "THE_MODEL_KB_CONFIG"
SECRET_AA = "AA_API_KEY"
PLACEHOLDER = "${AA_API_KEY}"

EXAMPLE_YAML = """# Live operator config. Mode 0600. No literal secrets.
secrets:
  - env: AA_API_KEY
    required: true
aa:
  api_key: "${AA_API_KEY}"
"""


class Keyring(Protocol):
    def get(self, service: str, account: str) -> str | None: ...

    def set(self, service: str, account: str, value: str) -> None: ...


class MemKeyring:
    def __init__(self) -> None:
        self.store: dict[tuple[str, str], str] = {}

    def get(self, service: str, account: str) -> str | None:
        return self.store.get((service, account))

    def set(self, service: str, account: str, value: str) -> None:
        self.store[(service, account)] = value


class MacosKeyring:
    def get(self, service: str, account: str) -> str | None:
        proc = subprocess.run(
            [
                "security",
                "find-generic-password",
                "-s",
                service,
                "-a",
                account,
                "-w",
            ],
            check=False,
            capture_output=True,
            text=True,
        )
        if proc.returncode != 0:
            return None
        value = proc.stdout.strip()
        return value or None

    def set(self, service: str, account: str, value: str) -> None:
        proc = subprocess.run(
            [
                "security",
                "add-generic-password",
                "-U",
                "-s",
                service,
                "-a",
                account,
                "-w",
                value,
            ],
            check=False,
            capture_output=True,
            text=True,
        )
        if proc.returncode != 0:
            err = (proc.stderr or proc.stdout or "keychain write failed").strip()
            raise RuntimeError(err)


def default_keyring() -> Keyring:
    if sys.platform == "darwin":
        return MacosKeyring()
    return MemKeyring()


def user_config_dir(xdg_home: str | None = None) -> Path:
    if xdg_home:
        return Path(xdg_home) / APP_DIR
    xdg = os.environ.get("XDG_CONFIG_HOME", "").strip()
    if xdg:
        return Path(xdg) / APP_DIR
    return Path.home() / ".config" / APP_DIR


def config_candidates(
    *,
    override: str | None = None,
    xdg_home: str | None = None,
    repo_root: Path | None = None,
) -> list[Path]:
    paths: list[Path] = []
    raw = override if override is not None else os.environ.get(OVERRIDE_ENV, "").strip()
    if raw:
        paths.append(Path(raw))
    paths.append(user_config_dir(xdg_home) / "config.yaml")
    if repo_root is not None:
        paths.append(repo_root / "site" / "config.yaml")
    return paths


def load_dotenv_file(path: Path, environ: dict[str, str]) -> None:
    if not path.is_file():
        return
    for raw in path.read_text(encoding="utf-8").splitlines():
        line = raw.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        key, _, value = line.partition("=")
        key = key.strip()
        value = value.strip().strip("'").strip('"')
        if key and key not in environ:
            environ[key] = value


def load_dotenv_files(repo_root: Path, environ: dict[str, str] | None = None) -> None:
    env = environ if environ is not None else os.environ
    load_dotenv_file(repo_root / ".env", env)
    load_dotenv_file(repo_root / "site" / ".env", env)


def declared_secrets(data: dict) -> list[tuple[str, bool]]:
    rows = data.get("secrets") or []
    if not isinstance(rows, list):
        return []
    out: list[tuple[str, bool]] = []
    for row in rows:
        if isinstance(row, str):
            out.append((row, False))
            continue
        if isinstance(row, dict) and row.get("env"):
            out.append((str(row["env"]), bool(row.get("required"))) )
    return out


def expand_placeholders(value: str, environ: dict[str, str]) -> str:
    out = value
    for key, val in environ.items():
        out = out.replace("${" + key + "}", val)
    return out


def resolve_aa_api_key(
    *,
    environ: dict[str, str] | None = None,
    keyring: Keyring | None = None,
    override: str | None = None,
    xdg_home: str | None = None,
    repo_root: Path | None = None,
) -> str:
    env = environ if environ is not None else os.environ
    from_env = env.get(SECRET_AA, "").strip()
    if from_env:
        return from_env

    ring = keyring if keyring is not None else default_keyring()
    for path in config_candidates(
        override=override, xdg_home=xdg_home, repo_root=repo_root
    ):
        if not path.is_file():
            continue
        data = yaml.safe_load(path.read_text(encoding="utf-8")) or {}
        if not isinstance(data, dict):
            continue
        names = declared_secrets(data)
        if not names:
            continue
        required = False
        for name, is_required in names:
            if name != SECRET_AA:
                continue
            required = required or is_required
            stored = ring.get(APP_DIR, name)
            if stored and stored.strip():
                env[name] = stored.strip()
        aa = data.get("aa") if isinstance(data.get("aa"), dict) else {}
        raw = ""
        if isinstance(aa, dict):
            raw = str(aa.get("api_key") or "")
        expanded = expand_placeholders(raw, dict(env)).strip()
        if expanded and PLACEHOLDER not in expanded:
            return expanded
        if required and not env.get(SECRET_AA, "").strip():
            raise ValueError(
                f"{SECRET_AA} required in {path} but env and keyring are empty"
            )
        break
    return env.get(SECRET_AA, "").strip()


def init_user_config(
    *,
    xdg_home: str | None = None,
    force: bool = False,
    example_src: Path | None = None,
) -> tuple[Path, Path]:
    dest_dir = user_config_dir(xdg_home)
    dest_dir.mkdir(parents=True, exist_ok=True)
    live = dest_dir / "config.yaml"
    example = dest_dir / "config.yaml.example"
    body = (
        example_src.read_text(encoding="utf-8")
        if example_src and example_src.is_file()
        else EXAMPLE_YAML
    )
    example.write_text(body, encoding="utf-8")
    if force or not live.is_file():
        live.write_text(body, encoding="utf-8")
        live.chmod(stat.S_IRUSR | stat.S_IWUSR)
    return live, example


def store_aa_key(
    value: str,
    *,
    keyring: Keyring | None = None,
) -> None:
    secret = value.strip()
    if not secret:
        raise ValueError(f"{SECRET_AA} is empty")
    ring = keyring if keyring is not None else default_keyring()
    ring.set(APP_DIR, SECRET_AA, secret)


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description="the-model-kb operator config")
    sub = parser.add_subparsers(dest="cmd", required=True)
    init_p = sub.add_parser("init", help="write ~/.config/the-model-kb/config.yaml")
    init_p.add_argument("--force", action="store_true")
    store_p = sub.add_parser(
        "store-aa-key", help="copy AA_API_KEY from env into the OS keyring"
    )
    args = parser.parse_args(argv)

    site_dir = Path(__file__).resolve().parents[1]
    repo_root = site_dir.parent
    example_src = site_dir / "config.yaml.example"

    if args.cmd == "init":
        live, example = init_user_config(force=args.force, example_src=example_src)
        print(f"wrote {example}")
        print(f"live {live}")
        return 0

    if args.cmd == "store-aa-key":
        load_dotenv_files(repo_root)
        value = os.environ.get(SECRET_AA, "")
        try:
            store_aa_key(value)
        except ValueError as exc:
            print(str(exc), file=sys.stderr)
            return 1
        except RuntimeError as exc:
            print(str(exc), file=sys.stderr)
            return 1
        print(f"stored {SECRET_AA} in keyring service={APP_DIR}")
        return 0

    return 1


if __name__ == "__main__":
    raise SystemExit(main())
