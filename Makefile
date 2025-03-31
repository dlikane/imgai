.DEFAULT_GOAL := help

.PHONY: deploy
deploy:  ## Install and deploy
	@echo "Installing globally..."
	go install ./...

# Helper commands
.PHONY: help
help: ## Display this help message
	@echo Available commands:
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the img binary
	@echo "Building ..."
	go build ./...

.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting Go code..."
	go fmt ./...

.PHONY: lint
lint: ## Run linter (golangci-lint)
	@echo "Running linter..."
	golangci-lint run

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	go test ./...

.PHONY: install-deps
install-deps: ## Install dependencies
	@echo "Installing dependencies..."
	go mod tidy

.PHONY: all
all: fmt lint test build ## Run format, lint, test, and build
