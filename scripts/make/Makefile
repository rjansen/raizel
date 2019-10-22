TOOLS_DIR             ?= /usr/local/bin
OS        	          ?= darwin
ARCH                  ?= amd64
include scripts/make/go/gotestsum/Makefile
include scripts/make/go/lint/Makefile
include scripts/make//go/lib/Makefile
CODECOV_COVERAGE_FILE ?= $(GO_COVERAGE_FILE)
include scripts/make/codecov/Makefile
export FIRESTORE_EMULATOR_HOST ?= $(FIRESTORE_HOST):$(FIRESTORE_PORT)
include scripts/make/firestore/Makefile
POSTGRES_SCRIPTS_DIR ?= $(BASE_DIR)/test/scripts/integration/postgres
include scripts/make/postgres/Makefile
MIGRATE_MIGRATIONS_DIR ?= $(BASE_DIR)/test/scripts/integration/postgres
MIGRATE_SQL_DSN ?= $(SQL_DSN)
include scripts/make/go/migrate/Makefile

.PHONY: ci
ci: go.vet go.coverage.text

.PHONY: firestore.ping
firestore.ping.%:
	@echo "$(NAME) firestore.ping"
	$(BASE_DIR)/scripts/bash/ping.sh "firestore" $(FIRESTORE_HOST) $(FIRESTORE_PORT) $*

.PHONY: postgres.ping
postgres.ping.%:
	@echo "$(NAME) postgres.ping"
	sleep 2
	$(BASE_DIR)/scripts/bash/ping.sh "postgres" $(POSTGRES_HOST) $(POSTGRES_PORT) $*

.PHONY: postgres.scripts
postgres.scripts: postgres.ping.5
	@echo "$(NAME) postgres.scripts"
	cat $(POSTGRES_SCRIPTS_DIR)/* | \
		psql -h $(POSTGRES_HOST) -U $(POSTGRES_USER) -d $(POSTGRES_DATABASE) -1 -f -
