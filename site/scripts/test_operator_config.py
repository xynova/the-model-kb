"""Tests for operator_config.py (no network, no Keychain)."""

from __future__ import annotations

import importlib.util
import stat
import tempfile
import unittest
from pathlib import Path

SCRIPTS = Path(__file__).resolve().parent
SPEC = importlib.util.spec_from_file_location(
    "operator_config", SCRIPTS / "operator_config.py"
)
assert SPEC and SPEC.loader
oc = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(oc)


class TestResolveOrder(unittest.TestCase):
    def test_env_wins_over_keyring(self) -> None:
        ring = oc.MemKeyring()
        ring.set(oc.APP_DIR, oc.SECRET_AA, "from-keyring")
        got = oc.resolve_aa_api_key(
            environ={"AA_API_KEY": "from-env"},
            keyring=ring,
        )
        self.assertEqual(got, "from-env")

    def test_override_config_then_keyring(self) -> None:
        ring = oc.MemKeyring()
        ring.set(oc.APP_DIR, oc.SECRET_AA, "from-keyring")
        with tempfile.TemporaryDirectory() as tmp:
            cfg = Path(tmp) / "config.yaml"
            cfg.write_text(oc.EXAMPLE_YAML, encoding="utf-8")
            got = oc.resolve_aa_api_key(
                environ={},
                keyring=ring,
                override=str(cfg),
            )
        self.assertEqual(got, "from-keyring")

    def test_required_empty_fails_closed(self) -> None:
        ring = oc.MemKeyring()
        with tempfile.TemporaryDirectory() as tmp:
            cfg = Path(tmp) / "config.yaml"
            cfg.write_text(oc.EXAMPLE_YAML, encoding="utf-8")
            with self.assertRaises(ValueError):
                oc.resolve_aa_api_key(
                    environ={},
                    keyring=ring,
                    override=str(cfg),
                )

    def test_dotenv_fills_missing_only(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            (root / ".env").write_text("AA_API_KEY=from-dotenv\n", encoding="utf-8")
            env = {"AA_API_KEY": "already"}
            oc.load_dotenv_files(root, env)
            self.assertEqual(env["AA_API_KEY"], "already")
            env2: dict[str, str] = {}
            oc.load_dotenv_files(root, env2)
            self.assertEqual(env2["AA_API_KEY"], "from-dotenv")


class TestInit(unittest.TestCase):
    def test_writes_0600_and_does_not_overwrite(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            live, example = oc.init_user_config(xdg_home=tmp)
            self.assertTrue(example.is_file())
            self.assertTrue(live.is_file())
            mode = stat.S_IMODE(live.stat().st_mode)
            self.assertEqual(mode, 0o600)
            live.write_text("keep\n", encoding="utf-8")
            oc.init_user_config(xdg_home=tmp, force=False)
            self.assertEqual(live.read_text(encoding="utf-8"), "keep\n")


class TestStore(unittest.TestCase):
    def test_store_aa_key(self) -> None:
        ring = oc.MemKeyring()
        oc.store_aa_key("  secret-value  ", keyring=ring)
        self.assertEqual(ring.get(oc.APP_DIR, oc.SECRET_AA), "secret-value")


if __name__ == "__main__":
    unittest.main()
