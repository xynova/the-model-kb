.PHONY: serve build clean generate-roster

HUGO ?= hugo
PYTHON ?= python3

generate-roster:
	$(PYTHON) scripts/generate-roster.py

serve: generate-roster
	$(HUGO) server --buildDrafts --disableFastRender --baseURL http://127.0.0.1:1313/

build: generate-roster
	$(HUGO) --gc --minify --printPathWarnings

clean:
	rm -rf public resources generated
