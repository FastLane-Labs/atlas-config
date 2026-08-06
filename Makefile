# Makefile for atlas-config project

# Variables
TS_DIR := typescript
GO_DIR := golang

# Load environment variables from .env file
-include .env
export

# Default target
all: build test

# Build all projects
build: clean build-ts build-go

# Build TypeScript project
build-ts:
	@echo "Building TypeScript project..."
	cd $(TS_DIR) && pnpm install && pnpm run build

# Build Go project
build-go:
	@echo "Building Go project..."
	cd $(GO_DIR) && go build ./...

# Test all projects
test: test-ts test-go

# Test TypeScript project
test-ts:
	@echo "Testing TypeScript project..."
	cd $(TS_DIR) && pnpm test

# Test Go project
test-go:
	@echo "Testing Go project..."
	cd $(GO_DIR) && go test -v ./...

# Clean all projects
clean: clean-ts clean-go

# Clean TypeScript project
clean-ts:
	@echo "Cleaning TypeScript project..."
	cd $(TS_DIR) && pnpm run clean

# Clean Go project
clean-go:
	@echo "Cleaning Go project..."
	cd $(GO_DIR) && go clean

# Publish TypeScript package
publish-ts: clean-ts build-ts
	@echo "Publishing TypeScript package..."
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found"; \
		exit 1; \
	fi
	@if [ -z "$(NPM_TOKEN)" ]; then \
		echo "Error: NPM_TOKEN not found in .env file"; \
		exit 1; \
	fi
	cd $(TS_DIR) && pnpm publish --access public --no-git-checks --//registry.npmjs.org/:_authToken=$(NPM_TOKEN)

.PHONY: all build build-ts build-go test test-ts test-go clean clean-ts clean-go publish-ts
