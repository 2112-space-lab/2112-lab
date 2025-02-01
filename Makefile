SHELL := /bin/bash
# Environment Variables
VERSION_GRAPHQL_API ?= 0.0.15
VERSION_GATEWAY_SERVICE ?= 0.0.1
VERSION_PROPAGATOR_SERVICE ?= 0.0.1

COMPOSE_FOLDER=deployments/compose
DOCKER_COMPOSE_FILE = deployments/compose/docker-compose.yaml
DOCKER_PROJECT_NAME = 2112_project
DOCKER_COMPOSE_GRAFANA=-f $(COMPOSE_FOLDER)/docker-compose-grafana.yml

# Service-Specific Settings
GRAPHQL_API = src/graphql-api
GATEWAY_SERVICE = src/gateway-service
PROPAGATOR_SERVICE = src/propagator-service
REDIS_SERVICE = redis-service

ACTIVATE = $(VENV_DIR)/bin/activate
VENV_DIR = $(GRAPHQL_API)/python/venv
PYTHON_BIN = $(VENV_DIR)/bin/python
PYTHON_VBIN = venv/bin/python
PIP_BIN = $(VENV_DIR)/bin/pip
# Exported Go Variables
export GO111MODULE := on

# Default Target
.DEFAULT_GOAL := help

DOCKERNET := 2112_net
export DOCKERNET

ifndef CURDIR
    CURDIR:=`pwd`
endif

# Define a variable to check for the --disable-pull flag
DISABLE_PULL := false

################################################################################
# LINT
################################################################################

.PHONY: lint
lint:
	docker run --rm -v $(CURDIR):/app -w /app golang.org/x/lint/golint@latest:latest golangci-lint run ./src/app-service -v --timeout=4m

##############################################################################
# Docker
##############################################################################
.PHONY: docker-network
docker-network: ## spin up the local docker network so that all dockerized components can communicate
	if [ $$( docker network ls -q --filter 'name=$(DOCKERNET)' | wc -l ) -eq 0 ]; then\
			docker network create $(DOCKERNET);\
	else\
		echo "Docker Network $(DOCKERNET) already created";\
	fi

.PHONY: grafana-up -- --disable-pull
grafana-up: docker-network ## launches the docker grafana configuration in Docker
	@if [ "$(DISABLE_PULL)" = "true" ]; then \
		echo "Skipping docker compose pull (d)"; \
	else \
		docker compose $(DOCKER_COMPOSE_GRAFANA) pull; \
	fi
	docker compose \
		--project-directory . \
		$(DOCKER_COMPOSE_GRAFANA) \
		up --force-recreate --build -d \
		$(CONTAINERS)

.PHONY: grafana-down
grafana-down: ## shuts down the docker grafana configuration
	docker-compose \
		--project-directory . \
		$(DOCKER_COMPOSE_GRAFANA) \
		down --volumes \

################################################################################
# Build
################################################################################

.PHONY: build
build: ## Build all services
	docker-compose -f $(DOCKER_COMPOSE_FILE) build

.PHONY: build-gateway-service
build-gateway-service: ## Build the Gateway service
	docker-compose -f $(DOCKER_COMPOSE_FILE) build $(GATEWAY_SERVICE)

.PHONY: build-propagator-service
build-propagator-service: ## Build the Propagator service
	docker-compose -f $(DOCKER_COMPOSE_FILE) build $(PROPAGATOR_SERVICE)

.PHONY: build-redis
build-redis: ## Build the Redis service
	docker-compose -f $(DOCKER_COMPOSE_FILE) build $(REDIS_SERVICE)

################################################################################
# Run
################################################################################

.PHONY: up
up: ## Start all services
	docker-compose -f $(DOCKER_COMPOSE_FILE) up

.PHONY: down
down: ## Stop all services
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

.PHONY: restart
restart: down up ## Restart all services

################################################################################
# Logs
################################################################################

.PHONY: logs
logs: ## View logs for all services
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

.PHONY: logs-gateway-service
logs-gateway-service: ## View logs for the Gateway service
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f $(GATEWAY_SERVICE)

.PHONY: logs-propagator-service
logs-propagator-service: ## View logs for the Propagator service
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f $(PROPAGATOR_SERVICE)

.PHONY: logs-redis
logs-redis: ## View logs for the Redis service
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f $(REDIS_SERVICE)

################################################################################
# Clean
################################################################################

.PHONY: clean
clean: ## Clean up all Docker resources
	docker-compose -f $(DOCKER_COMPOSE_FILE) down --volumes --remove-orphans
	docker system prune -f
	docker volume prune -f

################################################################################
# GraphQL API
################################################################################

.PHONY: gqlgen-generate
gqlgen-generate: ## Generate GraphQL code with gqlgen
	cd $(GRAPHQL_API)/go && go run github.com/99designs/gqlgen generate

.PHONY: gqlgen-init
gqlgen-init: ## Initialize the GraphQL API gqlgen project
	cd $(GRAPHQL_API)/go && go run github.com/99designs/gqlgen init

.PHONY: gqlgen-clean
gqlgen-clean: ## Clean GraphQL generated files
	rm -f $(GRAPHQL_API)/go/generated.go $(GRAPHQL_API)/go/models_gen.go

.PHONY: gqlgen-vendor
gqlgen-vendor: ## Update dependencies for gqlgen
	cd $(GRAPHQL_API)/go && go mod tidy && go mod vendor

.PHONY: gqlgen-run
gqlgen-run: ## Run the GraphQL API project
	cd $(GRAPHQL_API)/go && go run .

.PHONY: gqlgen-publish
gqlgen-publish: ## Publish the GraphQL API module
	cd src/graphql-api/go && \
	git add -A && \
	git commit -m "Release version $(VERSION_GRAPHQL_API)" && \
	git tag -a src/graphql-api/go/v$(VERSION_GRAPHQL_API) -m "Version $(VERSION_GRAPHQL_API)" && \
	git push origin main && \
	git push origin src/graphql-api/go/v$(VERSION_GRAPHQL_API)
	@echo "GraphQL API module published successfully with version $(VERSION_GRAPHQL_API)"



.PHONY: create-venv
create-venv: ## Create a virtual environment for the Python project
	@echo "Creating virtual environment..."
	python3 -m venv $(VENV_DIR)

.PHONY: install-python
install-python: create-venv ## Install Python dependencies in the virtual environment
	@echo "Installing Python dependencies..."
	$(PIP_BIN) install -r $(GRAPHQL_API)/python/requirements.txt
	$(PYTHON_BIN) -m pip install twine

.PHONY: generate-python
generate-python: install-python ## Generate Python code from the GraphQL schema
	@echo "Generating Python code from the GraphQL schema..."
	$(PYTHON_BIN) $(GRAPHQL_API)/python/generate_python_from_schema.py

.PHONY: run-python
run-python: ## Run the Python GraphQL server using the virtual environment
	@echo "Running the Python GraphQL server..."
	$(PYTHON_BIN) -m uvicorn src.graphql_api.python.server:app --host 0.0.0.0 --port 8000

.PHONY: clean-venv
clean-venv: ## Remove the virtual environment
	@echo "Removing virtual environment..."
	rm -rf $(VENV_DIR)

.PHONY: build-python-package
build-python-package: generate-python ## Build the Python package
	@echo "Building the Python package..."
	cd $(GRAPHQL_API)/python && $(PYTHON_VBIN) setup.py sdist bdist_wheel

.PHONY: publish-python
publish-python: build-python-package ## Publish the Python package to PyPI
	@echo "Publishing the Python package to PyPI..."
	cd $(GRAPHQL_API)/python && $(PYTHON_VBIN) -m twine upload dist/*
	@echo "Python package published successfully!"

################################################################################
# Gateway Service
################################################################################

.PHONY: gateway-run
gateway-run: ## Run the Gateway service
	cd $(GATEWAY_SERVICE) && python -m src.main

################################################################################
# Propagator Service
################################################################################

.PHONY: propagator-run
propagator-run: ## Run the Propagator service
	cd $(PROPAGATOR_SERVICE) && python -m src.main

################################################################################
# Help
################################################################################

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
