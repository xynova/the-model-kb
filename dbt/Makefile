# scores.json is produced by export_scores.py after dbt build, not by the Go CLI.
.PHONY: ensure parse build venv help

DBT ?= dbt
VENV ?= .venv
PYTHON ?= $(if $(wildcard $(VENV)/bin/python),$(VENV)/bin/python,python3)

.DEFAULT_GOAL := help

help: ## Show targets
	@grep -E '^[a-zA-Z0-9_.-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  %-12s %s\n", $$1, $$2}'

ensure: ## Print Fusion install URL if dbt missing; exit 1
	@command -v "$(DBT)" >/dev/null 2>&1 || { \
		echo 'dbt (Fusion v2) not found on PATH.' >&2; \
		echo 'Install:' >&2; \
		echo '  curl -fsSL https://public.cdn.getdbt.com/fs/install/install.sh | sh -s -- --update' >&2; \
		echo '  export PATH="$$HOME/.local/bin:$$PATH"' >&2; \
		echo '  dbt --version   # must report 2.x' >&2; \
		exit 1; \
	}

venv: ## Create .venv and install duckdb for export_scores.py
	@if [ ! -x $(VENV)/bin/python ]; then \
		python3 -m venv $(VENV); \
		$(VENV)/bin/pip install -q -r requirements.txt; \
	fi

parse: ensure ## dbt parse
	$(DBT) parse --profiles-dir .

build: ensure venv ## dbt build then export roster/scores.json
	@mkdir -p ../artifacts
	$(DBT) build --profiles-dir .
	$(PYTHON) scripts/export_scores.py
