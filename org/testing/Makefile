################################################################################
# TEST
################################################################################
export

DOCKER_IMAGE_POSTGRES:=postgis/postgis:latest
DOCKER_IMAGE_REDIS=redis:latest
DOCKER_IMAGE_RABBITMQ=rabbitmq:management
DOCKERNET:=2112_net

TEST_DB_HOST?=localhost
TEST_DB_PORT?=5432
TEST_DB_HOST_DOCKER?=2112-database.db
TEST_DB_PORT_DOCKER?=5432
TEST_DB_NAME?=postgres
TEST_DB_USER?=postgres
TEST_DB_PASS?=postgres

TEST_DB_CONTAINER_NAME?=2112-database.db
TEST_REDIS_CONTAINER_NAME=test-redis
TEST_RABBITMQ_CONTAINER_NAME=test-rabbitmq

TEST_REDIS_PORT=6379
TEST_REDIS_PORT_DOCKER=6379

TEST_RABBITMQ_USER=2112
TEST_RABBITMQ_PASS=2112
TEST_RABBITMQ_PORT=5672
TEST_RABBITMQ_PORT_DOCKER=5672
TEST_RABBITMQ_MANAGEMENT_PORT=15672
TEST_RABBITMQ_MANAGEMENT_PORT_DOCKER=15672

TEST_IT_RUNNING_ENV?=host
TEST_IT_DOCKER_NETWORK?=${DOCKERNET}
TEST_IT_TIMEOUT?=20m
TEST_IT_LOG_LEVEL?=debug
TEST_IT_GODOG_FLAGS?=~WIP

TEST_IT_ARTIFACTS_PATH?=${CURDIR}/integration/_artifacts

.PHONY: test-it-setup-up
test-it-setup-up:
	docker pull ${DOCKER_IMAGE_POSTGRES}
	docker pull ${DOCKER_IMAGE_REDIS}
	docker pull ${DOCKER_IMAGE_RABBITMQ}

	docker run -d \
		--name ${TEST_DB_CONTAINER_NAME} \
		--network=${TEST_IT_DOCKER_NETWORK} \
		-e POSTGRES_USER=${TEST_DB_USER} \
		-e POSTGRES_PASSWORD=${TEST_DB_PASS} \
		-e POSTGRES_DB=${TEST_DB_NAME} \
		-p ${TEST_DB_PORT}:${TEST_DB_PORT_DOCKER} \
		${DOCKER_IMAGE_POSTGRES}

	docker run -d \
		--name ${TEST_REDIS_CONTAINER_NAME} \
		--network=${TEST_IT_DOCKER_NETWORK} \
		-p ${TEST_REDIS_PORT}:${TEST_REDIS_PORT_DOCKER} \
		${DOCKER_IMAGE_REDIS}

	docker run -d \
		--name ${TEST_RABBITMQ_CONTAINER_NAME} \
		--network=${TEST_IT_DOCKER_NETWORK} \
		-e RABBITMQ_DEFAULT_USER=${TEST_RABBITMQ_USER} \
		-e RABBITMQ_DEFAULT_PASS=${TEST_RABBITMQ_PASS} \
		-p ${TEST_RABBITMQ_PORT}:${TEST_RABBITMQ_PORT_DOCKER} \
		-p ${TEST_RABBITMQ_MANAGEMENT_PORT}:${TEST_RABBITMQ_MANAGEMENT_PORT_DOCKER} \
		${DOCKER_IMAGE_RABBITMQ}


.PHONY: test-it
test-it:
	mkdir -p ${TEST_IT_ARTIFACTS_PATH}
	go test --race ./integration/ \
		-v \
		--test.timeout=${TEST_IT_TIMEOUT} \
		--tapp.options.app-workspace-path=${CURDIR} \
		--tapp.options.artifacts-path=${TEST_IT_ARTIFACTS_PATH} \
		--tapp.db.ssl-mode=disable \
		--tapp.options.preserve-databases=false \
		--godog.strict \
		--godog.tags=${TEST_IT_GODOG_FLAGS}

.PHONY: test-it-usage
test-it-usage:
	mkdir -p ${CURDIR}/integration/_artifacts
	go test --race ./integration/ \
		-v \
		--test.timeout=60s \
		--t2112.show-help

.PHONY: test-it-report
test-it-report:
	docker run --rm \
		-v ${CURDIR}/integration/_artifacts/results:/results \
		-v ${CURDIR}/integration/_artifacts/reports:/reports \
		-v ${CURDIR}/integration/_artifacts/suites:/suites \
		-v ${CURDIR}/integration/scripts:/scripts \
		node:14 bash -c "\
			find /suites/ -maxdepth 2 -exec ls -ld "{}" \; && \
			npm install cucumber-html-reporter cucumber-json-merge --save-dev;\
			node /scripts/merge-cucumber.js;\
			node /scripts/generate-report.js";