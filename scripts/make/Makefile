TOOLS_DIR             ?= /usr/local/bin
OS        	          ?= darwin
ARCH                  ?= amd64
include scripts/make/go/gotestsum/Makefile
include scripts/make/go/lint/Makefile
include scripts/make//go/lib/Makefile
CODECOV_COVERAGE_FILE ?= $(GO_COVERAGE_FILE)
include scripts/make/codecov/Makefile
include scripts/make/gcloud/Makefile
include scripts/make/gcloud/firestore/Makefile
POSTGRES_SCRIPTS_DIR ?= $(BASE_DIR)/test/scripts/integration/postgres
include scripts/make/postgres/Makefile
MIGRATE_MIGRATIONS_DIR ?= $(BASE_DIR)/test/scripts/integration/postgres
MIGRATE_SQL_DSN ?= $(SQL_DSN)
include scripts/make/go/migrate/Makefile

.PHONY: ci
ci: go.vet go.coverage.text

.PHONY: postgres.scripts
postgres.scripts:
	@echo "$(NAME) postgres.scripts"
	sleep 5
	cat $(POSTGRES_SCRIPTS_DIR)/* | \
		psql -h $(POSTGRES_HOST) -U $(POSTGRES_USER) -d $(POSTGRES_DATABASE) -1 -f -
