.PHONY: help run build test test-verbose clean dev install-deps tidy

help:
	@echo "Available commands:"
	@echo "  make run          - Run the application"
	@echo "  make dev          - Run with auto-reload (requires air)"
	@echo "  make build        - Build the application"
	@echo "  make test         - Run all tests"
	@echo "  make test-verbose - Run tests with verbose output"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make tidy         - Tidy and verify dependencies"
	@echo "  make install-deps - Install development dependencies"

run:
	@echo "Starting server..."
	go run ./cmd/api

dev:
	@echo "Starting development server with auto-reload..."
	@which air > /dev/null || (echo "Air not installed. Run 'make install-deps' first" && exit 1)
	air

build:
	@echo "Building application..."
	go build -o bin/api ./cmd/api
	codesign -s - bin/api
	@echo "Build complete: bin/api"

test:
	@echo "Running tests..."
	go test ./...

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

install-deps:
	@echo "Installing development dependencies..."
	go install github.com/cosmtrek/air@latest
	@echo "Done! You can now use 'make dev' for auto-reload"
