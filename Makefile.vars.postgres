# Sql vars
POSTGRES_USER        ?= postgres
POSTGRES_PASSWORD    ?=
POSTGRES_HOST        ?= localhost
POSTGRES_PORT        ?= 5432
POSTGRES_DATABASE    ?= postgres
SQL_DRIVER           ?= postgres
define SQL_DSN       =
$(SQL_DRIVER)://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable
endef
POSTGRES_SCRIPTS_DIR ?= $(BASE_DIR)/etc/test/integration/postgres/scripts


