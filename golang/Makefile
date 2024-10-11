.PHONY: build test clean


copy-config:
	cp ../configs/chain-config.json ./configs/chain-config.json

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
all: copy-config build test