GCLOUD_TMP_DIR  := $(TMP_DIR)/gcloud
GCLOUD_VERSION  := 267.0.0
# TODO: Create a map between os + arch filename and our scripts convention
GCLOUD_NAME     := google-cloud-sdk-$(GCLOUD_VERSION)-$(OS)-$(ARCH).tar.gz
GCLOUD_FILE     := $(GCLOUD_TMP_DIR)/$(GCLOUD_NAME)
GCLOUD_URL      := https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/$(GCLOUD_NAME)
GCLOUD_BIN      := gcloud
GCLOUD          := $(or $(shell which $(GCLOUD_BIN)),$(GCLOUD_TMP_DIR)/google-cloud-sdk/bin/$(GCLOUD_BIN))
PROJECT_ID                  ?= project-id
SERVICE_ACCOUNT             ?= service-account
SERVICE_ACCOUNT_DISPLAYNAME ?= "Service Account"
SERVICE_ACCOUNT_EMAIL       := $(SERVICE_ACCOUNT)@$(PROJECT_ID).iam.gserviceaccount.com
SERVICE_ACCOUNT_KEY         ?= $(HOME)/.$(PROJECT_ID)/$(SERVICE_ACCOUNT).json
SERVICE_ACCOUNT_KEY_DIR     := $(dir $(SERVICE_ACCOUNT_KEY))
export GOOGLE_APPLICATION_CREDENTIALS ?= $(SERVICE_ACCOUNT_KEY)

.PHONY: gcloud.clearcache
gcloud.clearcache:
	@echo "$(NAME)@$(BUILD) gcloud.clearcache"
	-rm -Rf $(GCLOUD_TMP_DIR)

$(GCLOUD_TMP_DIR):
	mkdir -p $(GCLOUD_TMP_DIR)

.PHONY: gcloud.install
gcloud.install: $(GCLOUD)
	@echo "$(NAME)@$(BUILD) gcloud.install $(GCLOUD)"

.PHONY: gcloud.uninstall
gcloud.uninstall:
	@echo "$(NAME)@$(BUILD) gcloud.uninstall $(GCLOUD)"
	-rm -f $(GCLOUD)

$(GCLOUD_FILE): | $(GCLOUD_TMP_DIR)
	@echo "$(NAM)@$(BUILD) $(GCLOUD_FILE)"
	curl -o $(GCLOUD_FILE) -L $(GCLOUD_URL)

$(GCLOUD): | $(GCLOUD_FILE)
	@echo "$(NAME)@$(BUILD) $(GCLOUD)"
	@cd $(GCLOUD_TMP_DIR) && tar xf $(GCLOUD_NAME) && google-cloud-sdk/install.sh
	$(GCLOUD) --help > /dev/null 2>&1

.PHONY: gcloud.createsa
gcloud.createsa:
	gcloud iam service-accounts create $(SERVICE_ACCOUNT) --display-name $(SERVICE_ACCOUNT_DISPLAYNAME)
	gcloud projects add-iam-policy-binding $(PROJECT_ID) \
	   	--member serviceAccount:$(SERVICE_ACCOUNT)@$(PROJECT_ID).iam.gserviceaccount.com

$(SERVICE_ACCOUNT_KEY):
	gcloud iam service-accounts keys create $(SERVICE_ACCOUNT_KEY) \
		--iam-account $(SERVICE_ACCOUNT)@$(PROJECT_ID).iam.gserviceaccount.com

.PHONY: gcloud.createkey
gcloud.createkey: $(SERVICE_ACCOUNT_KEY)
	@ls -l $(SERVICE_ACCOUNT_KEY)

.PHONY: gcloud.encodekeycontent
gcloud.encodekeycontent: $(SERVICE_ACCOUNT_KEY)
	@cat $(SERVICE_ACCOUNT_KEY) | base64 -w 0

.PHONY: gcloud.setup
gcloud.setup: gcloud.createsa gcloud.createkey

ifdef SERVICE_ACCOUNT_KEY_CONTENT
$(SERVICE_ACCOUNT_KEY_DIR):
	mkdir -p $(SERVICE_ACCOUNT_KEY_DIR)

gcloud.decodekeycontent: $(SERVICE_ACCOUNT_KEY_DIR)
	@echo $(SERVICE_ACCOUNT_KEY_CONTENT) | base64 --decode -w 0 > $(SERVICE_ACCOUNT_KEY)
	ls -l $(SERVICE_ACCOUNT_KEY)
endif
