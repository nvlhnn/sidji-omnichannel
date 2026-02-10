# Makefile for sidji

# Variables
BINARY_NAME=sidji-omnichannel
MAIN_PATH=cmd/server/main.go
DOCKER_COMPOSE=docker-compose

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Test variables
TEST_DB_URL=postgres://sidji:sidji123@localhost:5433/sidji_test?sslmode=disable

.PHONY: all build clean test test-coverage test-verbose run deps docker-up docker-down migrate swagger help

all: deps build

## Build
build:
	$(GOBUILD) -o bin/$(BINARY_NAME) $(MAIN_PATH)

## Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf bin/

## Run the server
run:
	$(GOCMD) run $(MAIN_PATH)

## Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## Run all tests
test:
	TEST_DATABASE_URL=$(TEST_DB_URL) $(GOTEST) -v ./...

## Run tests with coverage
test-coverage:
	TEST_DATABASE_URL=$(TEST_DB_URL) $(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

## Run tests in verbose mode
test-verbose:
	TEST_DATABASE_URL=$(TEST_DB_URL) $(GOTEST) -v -count=1 ./...

## Run only service tests
test-services:
	TEST_DATABASE_URL=$(TEST_DB_URL) $(GOTEST) -v ./internal/services/...

## Run only handler tests
test-handlers:
	TEST_DATABASE_URL=$(TEST_DB_URL) $(GOTEST) -v ./internal/handlers/...

## Run a specific test
test-one:
	@if [ -z "$(TEST)" ]; then echo "Usage: make test-one TEST=TestFunctionName"; exit 1; fi
	TEST_DATABASE_URL=$(TEST_DB_URL) $(GOTEST) -v -run $(TEST) ./...

## Setup test database
setup-test-db:
	@echo "Creating test database..."
	docker exec sidji-omnichannel-postgres psql -U sidji -c "DROP DATABASE IF EXISTS sidji_test;"
	docker exec sidji-omnichannel-postgres psql -U sidji -c "CREATE DATABASE sidji_test;"
	docker exec sidji-omnichannel-postgres psql -U sidji -d sidji_test -f /docker-entrypoint-initdb.d/001_initial_schema.up.sql || \
		docker exec -i sidji-omnichannel-postgres psql -U sidji -d sidji_test < scripts/setup_test_db.sql
	@echo "Test database ready!"

## Docker commands
docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-logs:
	$(DOCKER_COMPOSE) logs -f

## Run database migrations
migrate:
	@echo "Running migrations..."
	docker exec -i sidji-omnichannel-postgres psql -U sidji -d sidji < migrations/001_initial_schema.up.sql
	@echo "Migrations complete!"

## Generate Swagger documentation
swagger:
	swag init -g $(MAIN_PATH)
	@echo "Swagger docs generated at docs/"

## Lint the code
lint:
	golangci-lint run ./...

## Format the code
fmt:
	$(GOCMD) fmt ./...

## Check for vulnerabilities
vuln:
	$(GOCMD) install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

## Show help
help:
	@echo "sidji Makefile Commands:"
	@echo ""
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the server"
	@echo "  make test           - Run all tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make test-services  - Run only service tests"
	@echo "  make test-handlers  - Run only handler tests"
	@echo "  make test-one TEST=<name> - Run a specific test"
	@echo "  make setup-test-db  - Create and setup test database"
	@echo "  make docker-up      - Start Docker containers"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make migrate        - Run database migrations"
	@echo "  make swagger        - Generate Swagger documentation"
	@echo "  make lint           - Run linter"
	@echo "  make fmt            - Format code"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make deps           - Download dependencies"
	@echo ""
