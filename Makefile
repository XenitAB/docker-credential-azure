.ONESHELL:
SHELL := /bin/bash

TEST_ENV_FILE = tmp/test.env

ifneq (,$(wildcard $(TEST_ENV_FILE)))
    include $(TEST_ENV_FILE)
    export
endif

.PHONY: all
.SILENT: all
all: tidy lint fmt vet test

.PHONY: lint
.SILENT: lint
lint:
	set -e
	golangci-lint run

.PHONY: fmt
.SILENT: fmt
fmt:
	set -e
	go fmt ./...

.PHONY: tidy
.SILENT: tidy
tidy:
	set -e
	go mod tidy

.PHONY: vet
.SILENT: vet
vet:
	set -e
	go vet ./...

.SILENT: test
.PHONY: test
test:
	set -e
	go test -timeout 30s -cover ./...

.PHONY: cover
.SILENT: cover
cover:
	set -e
	mkdir -p tmp/
	go test -timeout 1m ./... -coverprofile=tmp/coverage.out
	if [[ "$${CI}" == "" ]]; then
		go tool cover -html=tmp/coverage.out
	fi
	
.SILENT: go-update
.PHONY: go-update
go-update:
	set -e
	go mod tidy
	go get -u ./...
	go mod tidy
	