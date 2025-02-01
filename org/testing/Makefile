################################################################################
# TEST
################################################################################
export

DOCKER_IMAGE_POSTGRES:=postgis/postgis:latest
DOCKERNET:=2112_net

TEST_DB_HOST?=localhost
TEST_DB_PORT?=5432
TEST_DB_HOST_DOCKER?=2112-database.db
TEST_DB_PORT_DOCKER?=5432
TEST_DB_NAME?=postgres
TEST_DB_USER?=postgres
TEST_DB_PASS?=postgres
TEST_DB_CONTAINER_NAME?=2112-database.db

TEST_IT_RUNNING_ENV?=host
TEST_IT_DOCKER_NETWORK?=${DOCKERNET}
TEST_IT_TIMEOUT?=20m
TEST_IT_LOG_LEVEL?=debug
TEST_IT_GODOG_FLAGS?=WIP

TEST_IT_ARTIFACTS_PATH?=${CURDIR}/integration/_artifacts

.PHONY: test-it-db-up
test-it-db-up:
	docker pull ${DOCKER_IMAGE_POSTGRES}
	docker run -d \
		--name ${TEST_DB_CONTAINER_NAME} \
		--network=${TEST_IT_DOCKER_NETWORK} \
		-e POSTGRES_USER=${TEST_DB_USER} \
		-e POSTGRES_PASSWORD=${TEST_DB_PASS} \
		-e POSTGRES_DB=${TEST_DB_NAME} \
		-p ${TEST_DB_PORT}:${TEST_DB_PORT_DOCKER} \
		${DOCKER_IMAGE_POSTGRES}

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