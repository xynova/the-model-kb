.PHONY: run build tidy clean smoke help

CMD := ./cmd/cloudrun-proxy
BINARY := bin/cloudrun-proxy

PROXY_URL := http://127.0.0.1:4318

.DEFAULT_GOAL := run

help: ## Show targets
	@grep -E '^[a-zA-Z0-9_.-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  %-12s %s\n", $$1, $$2}'

run: ## Run the proxy (go run)
	go run $(CMD)

build: ## Build bin/cloudrun-proxy
	@mkdir -p bin
	go build -o $(BINARY) $(CMD)

start: build ## Run the built binary
	$(BINARY)

tidy: ## go mod tidy
	go mod tidy

clean: ## Remove build artifacts
	rm -rf bin

smoke: ## Quick health + models check (proxy must be running)
	curl -sf $(PROXY_URL)/healthz
	@echo
	curl -sS $(PROXY_URL)/v1/models | jq .
