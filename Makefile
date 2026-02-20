#########################
# Run
#########################
.PHONY: run-api
run-api:
	go run ./cmd/api/

.PHONY: run-migrate
run-migrate:
	go run ./cmd/migrate/

#########################
# Docker Compose
#########################
.PHONY: dc-up
dc-up:
	docker compose -f docker-compose.yml up -d

.PHONY: dc-o11y-up
dc-o11y-up:
	docker compose -f docker-compose.o11y.yml up -d

.PHONY: dc-all-up
dc-all-up:
	docker compose -f docker-compose.yml -f docker-compose.o11y.yml up -d


#########################
# Database Migration
#########################
GOOSE_VERSION ?= v3.26.0
GOOSE_CMD = go run github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)

GOOSE_DRIVER ?= postgres
GOOSE_DBSTRING ?= postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
GOOSE_MIGRATION_DIR ?= internal/postgres/migration

GOOSE_ENV = \
	GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING=$(GOOSE_DBSTRING) \
	GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR)

.PHONY: migrate-up
migrate-up:
	$(GOOSE_ENV) $(GOOSE_CMD) up

.PHONY: migrate-down
migrate-down:
	$(GOOSE_ENV) $(GOOSE_CMD) down

.PHONY: migrate-status
migrate-status:
	$(GOOSE_ENV) $(GOOSE_CMD) status

.PHONY: migrate-create
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "invalid command. Usage: make migrate-create name=<migration_name>"; exit 1; \
	fi
	$(GOOSE_ENV) $(GOOSE_CMD) create "$(name)" sql

.PHONY: migrate-reset
migrate-reset:
	$(GOOSE_ENV) $(GOOSE_CMD) reset


#########################
# Testing
#########################
.PHONY: test
test:
	go test -v --failfast ./...

.PHONY: test-cov
test-cov:
	go test -coverprofile=bin/coverage.out ./...
	go tool cover -html=bin/coverage.out -o bin/coverage.html
	@echo "Coverage report saved to bin/coverage.html"


########################
# Lint
########################
.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0 run ./... --config .golangci.yml
