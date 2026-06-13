.PHONY: help build test test-server test-gateway lint fmt vet ci run-redis run-server run-gateway docker-up docker-down

GO      ?= go
BIN_DIR ?= bin

help: ## Show available targets
	@grep -E '^[a-zA-Z0-9_-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build server and gateway binaries to bin/
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/server ./cmd/server
	$(GO) build -o $(BIN_DIR)/gateway ./cmd/gateway

test: ## Run all tests
	$(GO) test ./... -race -count=1

test-server: ## Run server package tests
	$(GO) test ./internal/server/... -race -count=1

test-gateway: ## Run gateway package tests
	$(GO) test ./internal/gateway/... -race -count=1

fmt: ## Format all Go source files
	gofmt -w $$(go list -f '{{.Dir}}' ./...)

vet: ## Run go vet
	$(GO) vet ./...

lint: vet ## Run format check and vet
	@test -z "$$(gofmt -l $$(go list -f '{{.Dir}}' ./...))" || (echo "run make fmt" && exit 1)

ci: lint test build ## Run full CI pipeline (lint + test + build)

run-redis: ## Start Redis via docker compose
	docker compose up redis -d

run-server: ## Run backend server on :8080
	$(GO) run ./cmd/server

run-gateway: ## Run gateway on :8081 (requires Redis)
	REDIS_ADDR=localhost:6379 $(GO) run ./cmd/gateway

docker-up: ## Build and start all services via docker compose
	docker compose up --build

docker-down: ## Stop all docker compose services
	docker compose down
