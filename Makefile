# Makefile for atlas-config project

# Variables
TS_DIR := typescript
GO_DIR := golang

# Default target
all: build test

# Build all projects
build: build-ts build-go

# Build TypeScript project
build-ts:
	@echo "Building TypeScript project..."
	cd $(TS_DIR) && npm install && npm run build

# Build Go project
build-go:
	@echo "Building Go project..."
	cd $(GO_DIR) && go build ./...

# Test all projects
test: test-ts test-go

# Test TypeScript project
test-ts:
	@echo "Testing TypeScript project..."
	cd $(TS_DIR) && npm test

# Test Go project
test-go:
	@echo "Testing Go project..."
	cd $(GO_DIR) && go test -v ./...

# Clean all projects
clean: clean-ts clean-go

# Clean TypeScript project
clean-ts:
	@echo "Cleaning TypeScript project..."
	cd $(TS_DIR) && npm run clean

# Clean Go project
clean-go:
	@echo "Cleaning Go project..."
	cd $(GO_DIR) && go clean

.PHONY: all build build-ts build-go test test-ts test-go clean clean-ts clean-go
