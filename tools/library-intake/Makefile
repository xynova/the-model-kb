.PHONY: run build tidy clean test vet format help

CMD := ./cmd/library-intake
BINARY := bin/library-intake

.DEFAULT_GOAL := help

help: ## Show targets
	@grep -E '^[a-zA-Z0-9_.-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  %-12s %s\n", $$1, $$2}'

build: ## Build bin/library-intake
	@mkdir -p bin
	go build -o $(BINARY) $(CMD)

run: ## Run via go run (pass ARGS=...)
	go run $(CMD) $(ARGS)

tidy: ## go mod tidy
	go mod tidy

format: ## gofmt
	gofmt -w ./cmd ./internal

vet: ## go vet
	go vet ./...

test: ## go test
	go test ./...

clean: ## Remove build artifacts
	rm -rf bin
