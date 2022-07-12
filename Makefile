# Makefile
INSTALL_PATH ?= /usr/local/bin
BIN_NAME ?= schemacheck
BINDIR := $(CURDIR)/bin

default: build

.PHONY:  tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@goreleaser build \
		--rm-dist \
		--skip-validate \
		--single-target \
		--output dist/$(BIN_NAME)

.PHONY: install
install: build
	@install dist/$(BIN_NAME) $(INSTALL_PATH)/$(BIN_NAME)
	@schemacheck --version

.PHONY: release
release:
	@goreleaser build --rm-dist 

.PHONY: test
test:
	@go test -v

.PHONY: checks
checks:
	@go fmt ./...
	@go vet ./...
	@staticcheck ./...
	@gosec ./...
	@goimports -w ./
