.PHONY: serve build clean generate-roster check-links ensure-htmltest

HUGO ?= hugo
PYTHON ?= python3
HTMLTEST ?= htmltest
# Build with a root baseURL so htmltest resolves site-local paths under public/.
CHECK_BASE_URL ?= /

generate-roster:
	$(PYTHON) scripts/generate-roster.py

serve: generate-roster
	$(HUGO) server --buildDrafts --disableFastRender --baseURL http://127.0.0.1:1313/

build: generate-roster
	$(HUGO) --gc --minify --printPathWarnings

# Internal link gate (CI + pre-push). Skips external URLs; see .htmltest.yml.
check-links: ensure-htmltest generate-roster
	$(HUGO) --gc --minify --printPathWarnings --baseURL "$(CHECK_BASE_URL)"
	$(HTMLTEST) --conf .htmltest.yml

ensure-htmltest:
	@command -v "$(HTMLTEST)" >/dev/null 2>&1 || { \
		echo 'htmltest not found on PATH.' >&2; \
		echo 'Install: brew install htmltest' >&2; \
		echo 'Or download: https://github.com/wjdp/htmltest/releases/tag/v0.17.0' >&2; \
		exit 1; \
	}

clean:
	rm -rf public resources generated .htmltest
