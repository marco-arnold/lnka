.PHONY: help check build test fmt vet clean coverage

# Find goimports in PATH or ~/go/bin
GOIMPORTS := $(shell command -v goimports 2> /dev/null || echo $(HOME)/go/bin/goimports)

# Default target - show help
help:
	@echo "Available targets:"
	@echo "  make check      - Run fmt, vet, and test"
	@echo "  make fmt        - Format code with goimports and go fmt"
	@echo "  make vet        - Run go vet"
	@echo "  make test       - Run all tests"
	@echo "  make coverage   - Generate coverage report"
	@echo "  make build      - Build the project"
	@echo "  make clean      - Remove build artifacts"

# Run all checks (recommended before commit)
check: fmt vet test

# Format code and fix imports
fmt:
	@echo "Running goimports..."
	@$(GOIMPORTS) -w .
	@echo "Running go fmt..."
	@go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run all tests
test:
	@echo "Running tests..."
	@go test ./...

# Generate coverage report
coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build the project
build:
	@echo "Building..."
	@go build -o lnka .

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f lnka coverage.out coverage.html debug.log
	@echo "Done"
