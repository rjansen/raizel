# Sql vars
POSTGRES_VERSION  ?= 12
POSTGRES_USER     ?= postgres
POSTGRES_PASSWORD ?=
POSTGRES_HOST     ?= localhost
POSTGRES_PORT     ?= 5432
POSTGRES_DATABASE ?= postgres
SQL_DRIVER        ?= postgres
SQL_DSN           := $(SQL_DRIVER)://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable

.PHONY: postgres.start
postgres.start:
	@echo "$(REPO) postgres.start"
	@#docker run --rm -d --name postgres-run -p $(POSTGRES_PORT):$(POSTGRES_PORT) \
	@#	-v "$(POSTGRES_SCRIPTS_DIR):/docker-entrypoint-initdb.d" postgres:$(POSTGRES_VERSION)
	@docker run --rm -d --name postgres-run -p $(POSTGRES_PORT):$(POSTGRES_PORT) postgres:$(POSTGRES_VERSION)

.PHONY: postgres.kill
postgres.kill:
	@echo "$(REPO) postgres.kill"
	docker kill postgres-run

.PHONY: postgres.psql
postgres.psql:
	@echo "$(REPO) postgres.psql"
	docker run --rm -it --name psql-run postgres:$(POSTGRES_VERSION) \
		psql -h $(POSTGRES_HOST) -U $(POSTGRES_USER) -d $(POSTGRES_DATABASE)
