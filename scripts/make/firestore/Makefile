FIRESTORE_HOST ?= 127.0.0.1
FIRESTORE_PORT ?= 8086
FIRESTORE_DOCKER_NAME ?= firestore-run
FIRESTORE_DOCKER_IMAGE ?= rjansen/firestore-emulator

.PHONY: firestore.start
firestore.start:
	@echo "$(NAME) firestore.start"
	docker run --rm -d --name $(FIRESTORE_DOCKER_NAME) -p $(FIRESTORE_PORT):$(FIRESTORE_PORT) $(FIRESTORE_DOCKER_IMAGE)

.PHONY: firestore.kill
firestore.kill:
	@echo "$(NAME) firestore.start"
	docker kill $(FIRESTORE_DOCKER_NAME)
