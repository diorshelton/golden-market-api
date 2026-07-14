.PHONY: help run build build-mac test test-verbose test-coverage clean dev fmt vet check tidy

help:
	@echo "Available commands:"
	@echo "  make run          - Run the application"
	@echo "  make dev          - Run with auto-reload (requires air)"
	@echo "  make build        - Build the application"
	@echo "  make build-mac    - Local macOS build with code signing"
	@echo "  make fmt          - Format the code using gofmt"
	@echo "  make vet          - Vet the code for potential issues"
	@echo "  make test         - Run all tests"
	@echo "  make check        - Run fmt vet and test"
	@echo "  make test-verbose - Run tests with verbose output"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make tidy         - Tidy and verify dependencies"

run:
	@echo "Starting server..."
	go run ./cmd/api

dev:
	@echo "Starting development server with auto-reload..."
	@which air > /dev/null || (echo "Air not installed. Run 'go install github.com/air-verse/air@latest' first" && exit 1)
	air

# Generic build (works everywhere)
build:
	@echo "Building application..."
	@mkdir -p bin
	go build -o bin/api ./cmd/api
	@echo "Build complete: bin/api"

# Local macOS build with signing
build-mac: build
	@echo "Code signing for macOS..."
	codesign -s - bin/api

fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted."

vet:
	@echo "Vetting code..."
	go vet ./...
	@echo "Code vetted."

test:
	@echo "Running tests..."
	go test ./...

check: fmt vet test
	@echo "All checks passed!"

test-verbose:
	@echo "Running tests (verbose)..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

tidy:
	@echo "Tidying dependencies..."
	go mod tidy
	go mod verify

