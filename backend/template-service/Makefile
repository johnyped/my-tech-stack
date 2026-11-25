.PHONY: help setup config-init run run-docker test test-coverage test-integration \
	build clean lint vendor docker-build docker-run compose-up compose-down compose-logs \
	tidy deps deps-verify fmt vet \
	integration-start-mock integration-test-it integration-test-full integration-test-health \
	bruno-gen

SERVICE_NAME ?= TODO-service-name
BINARY := bin/$(SERVICE_NAME)

help: ## Show this help
	@grep -E '^[a-zA-Z0-9_.-]+:.*## ' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*## "}; {printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'

setup: ## Install local tools (golangci-lint, bru CLI) and init config
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.0
	npm install -g @usebruno/cli
	$(MAKE) config-init

config-init: ## Copy config.example.yaml to config.yaml if missing
	test -f config/config.yaml || cp config/config.example.yaml config/config.yaml

run: ## Run the application locally
	CONFIG_FILE=config/config.yaml go run cmd/main.go

run-docker: ## Run with docker-compose configuration (requires compose-up for deps)
	CONFIG_FILE=config/config.docker-compose.yaml go run cmd/main.go

test: ## Run unit tests
	go test -mod=vendor -race ./internal/... ./pkg/...

test-coverage: ## Run unit tests with coverage report
	go test -mod=vendor -race -coverprofile=coverage.out ./internal/... ./pkg/...
	go tool cover -func=coverage.out | tail -1

test-integration: ## Run integration tests against docker-compose.integration.yaml
	docker compose -f docker-compose.integration.yaml up --build --abort-on-container-exit
	docker compose -f docker-compose.integration.yaml down -v

build: ## Build binary
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags="-w -s" -o $(BINARY) cmd/main.go

clean: ## Remove build artifacts
	rm -rf bin/ coverage.out

lint: ## Run golangci-lint
	golangci-lint run

vendor: ## Vendor dependencies (required before docker build; needs local creds for private modules)
	go mod vendor -v

tidy: ## Tidy go modules
	go mod tidy

deps: ## Download dependencies
	go mod download

deps-verify: ## Verify module checksums
	go mod verify

fmt: ## Format code
	gofmt -l -w .

vet: ## Run go vet
	go vet ./...

docker-build: vendor ## Build the Docker image
	docker build -t $(SERVICE_NAME):local .

docker-run: ## Run the built image standalone
	docker run --rm -p 8080:8080 --env CONFIG_FILE=/app/config/config.docker-compose.yaml $(SERVICE_NAME):local

compose-up: ## Start the full local stack (app + whichever infra blocks are enabled)
	docker compose up --build -d

compose-down: ## Stop the local stack
	docker compose down -v

compose-logs: ## Tail local stack logs
	docker compose logs -f

## --- Bruno / integration tests (see document/MOCKING_GUIDE.md) ---

bruno-gen: ## Generate/refresh Bruno API docs from the codebase (skill::bruno-service-docs)
	@echo "Run inside Kiro: 'generate bruno docs' (loads skill::bruno-service-docs)"

integration-start-mock: ## Start Mockoon + this service only (foreground, Ctrl+C to stop)
	./script/integration/run-local.sh $(if $(SCENARIO),--scenario "$(SCENARIO)",)

integration-test-it: ## Boot service + run bruno/Integration Tests against LOCAL-IT, then tear down
	./script/integration/run-integration-tests.sh --env LOCAL-IT --folder "Integration Tests"

integration-test-full: ## Boot service + run every Bruno folder, then tear down
	./script/integration/run-integration-tests.sh

integration-test-health: ## Run only the Health folder against a running instance
	./script/integration/run-bruno-tests.sh --env LOCAL --folder "Health"
