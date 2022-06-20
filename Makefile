# Makefile
INSTALL_PATH ?= /usr/local/bin
BIN_NAME ?= schemacheck
BINDIR := $(CURDIR)/bin

.PHONY: tidy build test checks clean release

default: build

tidy:
	@go mod tidy

build:
	@go build -o dist/bin/schemacheck

release:
	@goreleaser build --rm-dist 

test:
	@go test -v

checks:
	@go fmt ./...
	@go vet ./...
	@staticcheck ./...
	@gosec ./...
	@goimports -w ./
