.PHONY: build test clean

# Build the package
build:
	go build ./...

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	go clean
	rm -f atlasconfig

# Run tests with coverage
cover:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Check for lint issues
lint:
	golangci-lint run

# Run all checks (lint and test)
check: lint test

# Default target
all: build test