# Makefile
INSTALL_PATH ?= /usr/local/bin
BIN_NAME ?= schemacheck
BINDIR := $(CURDIR)/bin

.PHONY: tidy build test checks clean

default: build

tidy:
	@go mod tidy

build:
	@goreleaser build --rm-dist --skip-validate

test:
	@go test -v

checks:
	@go fmt ./...
	@go vet ./...
	@staticcheck ./...
	@gosec ./...
	@goimports -w ./
