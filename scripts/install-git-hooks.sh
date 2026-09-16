#!/usr/bin/env bash
# Point this clone at the versioned hooks under scripts/githooks/.
set -euo pipefail

repo_root="$(cd "$(dirname "$0")/.." && pwd)"
hooks_path="scripts/githooks"

cd "$repo_root"
chmod +x "${hooks_path}/pre-push"
git config core.hooksPath "$hooks_path"
echo "Configured core.hooksPath=${hooks_path}"
echo "pre-push will run site link checks when site/library/roster/models change."
