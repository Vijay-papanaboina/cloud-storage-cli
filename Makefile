.PHONY: build build-dev build-staging build-prod test test-verbose test-race test-cover help

# Load .env file if it exists
-include .env
export

# Default API URL if not set in .env
API_URL ?= http://localhost:8000

# Output binary name
OUTPUT = cloud-storage-api-cli

help:
	@echo "Available targets:"
	@echo "  make build        - Build with API_URL from .env (default: http://localhost:8000)"
	@echo "  make build-dev    - Build for development"
	@echo "  make build-staging - Build for staging"
	@echo "  make build-prod   - Build for production"
	@echo "  make test         - Run all tests"
	@echo "  make test-verbose - Run all tests with verbose output"
	@echo "  make test-race    - Run all tests with race detector"
	@echo "  make test-cover   - Run all tests with coverage report"
	@echo ""
	@echo "Create a .env file with:"
	@echo "  API_URL=http://your-api-url.com"

build:
	@echo "Building with API_URL: $(API_URL)"
	go build -ldflags "-X github.com/vijay-papanaboina/cloud-storage-api-cli/internal/config.BuildTimeAPIURL=$(API_URL)" -o $(OUTPUT) .

build-dev:
	@echo "Building for development..."
	@$(MAKE) build API_URL=http://localhost:8000

build-staging:
	@echo "Building for staging..."
	@$(MAKE) build API_URL=https://api.staging.com

build-prod:
	@echo "Building for production..."
	@$(MAKE) build API_URL=https://api.production.com

test:
	@echo "Running tests..."
	go test ./...

test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./...

test-race:
	@echo "Running tests with race detector..."
	go test -race ./...

test-cover:
	@echo "Running tests with coverage..."
	go test -cover ./...
