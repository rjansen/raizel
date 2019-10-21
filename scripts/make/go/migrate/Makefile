MIGRATE_TMP_DIR  := $(TMP_DIR)/migrate
MIGRATE_VERSION  := 4.6.2
# TODO: Create a map between os + arch filename and our scripts convention
MIGRATE_NAME     := migrate.$(OS)-$(ARCH).tar.gz
MIGRATE_FILE     := $(MIGRATE_TMP_DIR)/$(MIGRATE_NAME)
MIGRATE_URL      := https://github.com/golang-migrate/migrate/releases/download/v$(MIGRATE_VERSION)/$(MIGRATE_NAME)
MIGRATE_BIN      := migrate
MIGRATE          := $(or $(shell which $(MIGRATE_BIN)),$(MIGRATE_TMP_DIR)/$(MIGRATE_BIN))

.PHONY: migrate.clearcache
migrate.clearcache:
	@echo "$(NAME)@$(BUILD) migrate.clearcache"
	-rm -Rf $(MIGRATE_TMP_DIR)

$(MIGRATE_TMP_DIR):
	mkdir -p $(MIGRATE_TMP_DIR)

.PHONY: migrate.install
migrate.install: $(MIGRATE)
	@echo "$(NAME)@$(BUILD) migrate.install"

.PHONY: migrate.uninstall
migrate.uninstall:
	@echo "$(NAME)@$(BUILD) migrate.uninstall $(MIGRATE)"
	-rm -f $(MIGRATE)

$(MIGRATE_FILE): | $(MIGRATE_TMP_DIR)
	@echo "$(NAME)@$(BUILD) $(MIGRATE_FILE)"
	curl -o $(MIGRATE_FILE) -L $(MIGRATE_URL)

$(MIGRATE): | $(MIGRATE_FILE)
	@echo "$(NAME)@$(BUILD) $(MIGRATE)"
	cd $(MIGRATE_TMP_DIR) && tar xf $(MIGRATE_NAME) && cp -vf $(MIGRATE_BIN).$(OS)-$(ARCH) $(MIGRATE)
	$(MIGRATE) -help > /dev/null 2>&1

migrate.%:
	$(MIGRATE) -source file://$(MIGRATE_MIGRATIONS_DIR) -database $(MIGRATE_SQL_DSN) $*
