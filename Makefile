APP_BIN := "./bin/app"

GIT_HASH := $(shell git log --pretty="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.3.0
.PHONY: install-lint-deps

lint: install-lint-deps
	golangci-lint run ./... --config=.golangci.yml
.PHONY: lint

build:
	go build -v -o $(APP_BIN) -ldflags "$(LDFLAGS)" ./cmd/app
.PHONY: build

run: build
	$(APP_BIN)
.PHONY: run

test:
	go test -race -count 100 ./internal/...
.PHONY: test
