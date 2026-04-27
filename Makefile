.DEFAULT_GOAL := help

COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

.PHONY: help check fmt lint vet test integration coverage clean compose-up compose-down

help:
	@echo "Targets:"
	@echo "  check        Verify formatting + vet + lint without modifying"
	@echo "  fmt          Format all Go sources"
	@echo "  lint         Run golangci-lint --fix"
	@echo "  vet          Run go vet"
	@echo "  test         Run unit tests (in-memory SQLite, no Docker)"
	@echo "  integration  Spin up compose, run -tags=integration, tear down"
	@echo "  coverage     Run tests with coverage and HTML report"
	@echo "  compose-up   docker compose up -d (postgres + oracle)"
	@echo "  compose-down docker compose down"

check:
	@echo "==> gofmt"
	@UNFORMATTED=$$(gofmt -l .); \
	if [ -n "$$UNFORMATTED" ]; then \
		echo "Files need formatting:"; echo "$$UNFORMATTED"; exit 1; \
	fi
	@echo "==> go vet"
	@go vet ./...
	@echo "==> golangci-lint"
	@golangci-lint run

fmt:
	@gofmt -w .

lint:
	@golangci-lint run --fix

vet:
	@go vet ./...

test:
	@go test -short ./...

compose-up:
	@docker compose up -d
	@echo "Waiting for postgres..."
	@until docker compose exec -T postgres pg_isready -U raizel >/dev/null 2>&1; do sleep 1; done
	@echo "Waiting for oracle (this can take ~60s on first boot)..."
	@until docker compose exec -T oracle healthcheck.sh >/dev/null 2>&1; do sleep 5; done
	@echo "Both databases healthy."

compose-down:
	@docker compose down

# Default DSNs for local docker-compose; override via environment.
POSTGRES_DSN ?= postgres://raizel:raizel@localhost:5432/raizel?sslmode=disable
ORACLE_DSN   ?= oracle://raizel:raizel@localhost:1521/FREEPDB1

integration: compose-up
	@POSTGRES_DSN='$(POSTGRES_DSN)' ORACLE_DSN='$(ORACLE_DSN)' go test -tags=integration ./...
	@$(MAKE) compose-down

coverage:
	@go test -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@go tool cover -func=$(COVERAGE_FILE) | tail -1

clean:
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
