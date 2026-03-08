DIR = $(shell pwd)
SERVICES_DIR = $(DIR)/services
CONFIG_PATH = $(DIR)/config
IDL_PATH = $(DIR)/interfaces
OUTPUT_PATH = $(DIR)/output

SERVICES := api user follow interaction video chat
service = $(word 1, $@)

# mock gen
MOCKS := user_mock
mock = $(word 1, $@)

PERFIX = "[Makefile]"

.PHONY: env-up
env-up:
	sh init.sh
	docker-compose up -d

.PHONY: env-down
env-down:
	docker-compose down

.PHONY: $(SERVICES)
$(SERVICES):
	mkdir -p output
	cd $(SERVICES_DIR)/$(service) && sh build.sh
	cd $(SERVICES_DIR)/$(service)/output && cp -r . $(OUTPUT_PATH)/$(service)
	@echo "$(PERFIX) Build $(service) target completed"


.PHONY: $(MOCKS)
$(MOCKS):
	@mkdir -p mocks
	mockgen -source=./interfaces/$(mock).go -destination=./mocks/$(mock).go -package=mocks


.PHONY: clean
clean:
	@find . -type d -name "output" -exec rm -rf {} + -print


.PHONY: build-all
build-all:
	@for var in $(SERVICES); do \
		echo "$(PERFIX) building $$var service"; \
		make $$var ci=1; \
	done

.PHONY: docker
docker:
	docker build -t tiktok .
