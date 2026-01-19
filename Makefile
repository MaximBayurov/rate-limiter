APP_BIN := "./bin/app"

GIT_HASH := $(shell git log --pretty="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

DOCKER_COMPOSE_INTEGRATION=docker-compose --file=deployments/integration-compose.yaml --env-file=testing.env

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.3.0
.PHONY: install-lint-deps

lint: install-lint-deps
	golangci-lint run ./... --config=.golangci.yml
.PHONY: lint

build:
	go build -v -o $(APP_BIN) -ldflags "$(LDFLAGS)" ./cmd/app
.PHONY: build

rebuild-containers:
	 docker-compose -f ./deployments/app-compose.yaml build --no-cache
.PHONY: rebuild-containers

run:
	docker-compose -f ./deployments/app-compose.yaml up -d
.PHONY: run

down:
	docker-compose -f ./deployments/app-compose.yaml down
.PHONY: run

test:
	go test -race ./internal/...
.PHONY: test

rebuild-testing:
	$(DOCKER_COMPOSE_INTEGRATION) --profile=tests build --no-cache
.PHONY: rebuild-testing

up-testing:
	$(DOCKER_COMPOSE_INTEGRATION) up -d
.PHONY: up-testing

down-testing:
	$(DOCKER_COMPOSE_INTEGRATION) down
.PHONY: down-testing

integration-tests:
	$(DOCKER_COMPOSE_INTEGRATION) up -d --wait

	$(DOCKER_COMPOSE_INTEGRATION) run --rm -P integration-tests -test.v

	TEST_EXIT_CODE=$$?

	make down-testing

	exit $$TEST_EXIT_CODE
.PHONY: integration-tests