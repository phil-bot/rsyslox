# rsyslox Makefile

BINARY       := rsyslox
BUILD_DIR    := build
VERSION      ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS      := -s -w -X main.Version=$(VERSION)

FRONTEND_DIR := frontend
REDOC_JS     := docs/api-ui/redoc.standalone.js
REDOC_URL    := https://cdn.jsdelivr.net/npm/redoc/bundles/redoc.standalone.js

.PHONY: all build build-static frontend redoc dev clean test lint install uninstall help

## all: Build everything — frontend + redoc + Go binary
all: frontend redoc lint build

## build: Build Go binary (development, requires frontend/dist to exist)
build:
	mkdir -p $(BUILD_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY) .

## build-static: Build fully static Go binary (for Docker / production)
build-static:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY) .

## frontend: Install npm deps and build Vue app into frontend/dist/
frontend:
	cd $(FRONTEND_DIR) && npm install && npm run build

## redoc: Download Redoc standalone JS for offline API documentation
redoc:
	@echo "Downloading Redoc standalone..."
	mkdir -p docs/api-ui
	curl -fsSL $(REDOC_URL) -o $(REDOC_JS)
	@echo "Redoc downloaded to $(REDOC_JS)"

## dev: Run Go backend for frontend development (uses config.dev.toml)
dev:
	RSYSLOX_CONFIG=./config.dev.toml go run .

## test: Run all Go tests
test:
	go test ./...

## lint: Run Go static analysis
lint:
	go vet ./...

## install: Install rsyslox on this machine (requires sudo)
install: all
	sudo scripts/install.sh

## uninstall: Remove rsyslox from this machine (requires sudo)
uninstall:
	sudo scripts/install.sh --uninstall

## clean: Remove all build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/node_modules

## help: Show available targets
help:
	@grep -E '^## ' Makefile | sed 's/## /  /'
