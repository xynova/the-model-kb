.PHONY: serve build clean generate-roster retrieve-limits apply-mine retrieve-aa apply-aa mine-hosts polypus-diff check-links ensure-htmltest test-roster init store-aa-key

HUGO ?= hugo
HTMLTEST ?= htmltest
# Build with a root baseURL so htmltest resolves site-local paths under public/.
CHECK_BASE_URL ?= /
MILL ?= go -C ../tools/roster-mill run ./cmd/roster-mill

generate-roster:
	$(MILL) generate

retrieve-limits:
	$(MILL) retrieve-limits

apply-mine:
	$(MILL) apply-mine

init:
	$(MILL) init

store-aa-key:
	$(MILL) store-aa-key

retrieve-aa:
	$(MILL) retrieve-aa

apply-aa:
	$(MILL) apply-aa

mine-hosts: generate-roster retrieve-limits apply-mine retrieve-aa apply-aa generate-roster

polypus-diff:
	$(MILL) polypus-diff

test-roster:
	cd ../tools/roster-mill && go test ./...

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
