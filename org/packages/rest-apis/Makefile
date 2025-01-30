API_NAME := rest-apis
API_VERSION := $(shell cat VERSION)
API_ROOT_FILE := api/app-rest.api.yaml
LINTER := redocly/cli
LINTER_VERSION := 1.16.0

lint:
	docker run -t --rm -v $(CURDIR):/spec $(LINTER):$(LINTER_VERSION) lint $(API_ROOT_FILE) --config=configs/redocly.yaml

generate: generate-golang

version:
	cd codegen/angular/ && npm version $(API_VERSION)
docs:
	docker run -t --rm -v $(CURDIR):/spec $(LINTER):$(LINTER_VERSION) build-docs $(API_ROOT_FILE) --output=docgen/api.html

generate-golang:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	go generate ./codegen/go/codegen.go