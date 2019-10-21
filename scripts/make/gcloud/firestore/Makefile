FIRESTORE_HOST ?= 127.0.0.1
FIRESTORE_PORT ?= 8086

.PHONY: firestore.emulator.start
firestore.emulator.start:
	gcloud beta emulators firestore start --host-port=$(FIRESTORE_HOST):$(FIRESTORE_PORT)

.PHONY: firestore.createsa
firestore.createsa:
	gcloud iam service-accounts create $(SERVICE_ACCOUNT) --display-name $(SERVICE_ACCOUNT_DISPLAYNAME)
	gcloud projects add-iam-policy-binding $(PROJECT_ID) \
	   	--member serviceAccount:$(SERVICE_ACCOUNT)@$(PROJECT_ID).iam.gserviceaccount.com \
		--role roles/datastore.user

.PHONY: firestore.createkey
firestore.createkey: $(SERVICE_ACCOUNT_KEY)
	@ls -l $(SERVICE_ACCOUNT_KEY)

.PHONY: firestore.encodekeycontent
firestore.encodekeycontent: $(SERVICE_ACCOUNT_KEY)
	@cat $(SERVICE_ACCOUNT_KEY) | base64 -w 0

.PHONY: firestore.setup
firestore.setup: firestore.createsa firestore.createkey
