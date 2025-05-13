# Makefile for Bishoujo-Huntress

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOCOVER=$(GOCMD) tool cover
GOLINT=golangci-lint

# Project parameters
BINARY_NAME=bishoujo-huntress
PKG=github.com/greysquirr3l/bishoujo-huntress
CMD_DIR=./cmd
EXAMPLES_DIR=$(CMD_DIR)/examples
BUILD_DIR=./build
COVERAGE_DIR=./coverage
DOCS_DIR=./docs

# Version information - can be overridden by git tags or CI/CD
VERSION ?= $(shell git describe --tags --always --dirty || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# LDFLAGS for version information
LD_FLAGS = -X '$(PKG)/pkg/huntress.Version=$(VERSION)' \
		   -X '$(PKG)/pkg/huntress.Commit=$(COMMIT)' \
		   -X '$(PKG)/pkg/huntress.BuildDate=$(BUILD_DATE)'

# Build targets
.PHONY: all build clean deps examples test test-race test-cover lint vet fmt tidy check vendor doc security-check help

# Default target
all: clean deps lint test build

# Build binary
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -ldflags "$(LD_FLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/...
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build examples
examples:
	@echo "Building examples..."
	@mkdir -p $(BUILD_DIR)/examples
	@for dir in $$(find $(EXAMPLES_DIR) -type d -mindepth 1); do \
		example_name=$$(basename $$dir); \
		$(GOBUILD) -o $(BUILD_DIR)/examples/$$example_name ./cmd/examples/$$example_name; \
		echo "Built example: $(BUILD_DIR)/examples/$$example_name"; \
	done

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(COVERAGE_DIR)
	@echo "Cleaned"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@$(GOMOD) tidy
	@$(GOMOD) download
	@echo "Checking for golangci-lint..."
	@if ! command -v $(GOLINT) >/dev/null 2>&1; then \
	  echo "golangci-lint not found, installing v1.56.2..."; \
	  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.56.2; \
	  if ! command -v $(GOLINT) >/dev/null 2>&1; then \
		echo "golangci-lint still not found. Please ensure $(GOLINT) is in your PATH."; \
		exit 1; \
	  fi; \
	else \
	  echo "golangci-lint found: $$($(GOLINT) --version)"; \
	fi
	@echo "Dependencies installed"

# Run all Go fuzz targets in pkg/huntress for a short duration
fuzz:
	@echo "Running all Go fuzz targets in pkg/huntress..."
	go test -fuzz=FuzzIncidentListOptionsValidate -fuzztime=10s ./pkg/huntress || exit 1
	go test -fuzz=FuzzEncodeURLValues -fuzztime=10s ./pkg/huntress || exit 1
	go test -fuzz=FuzzAddQueryParams -fuzztime=10s ./pkg/huntress || exit 1
	go test -fuzz=FuzzExtractPagination -fuzztime=10s ./pkg/huntress || exit 1

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	$(GOTEST) -v -race ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	$(GOCOVER) -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated at $(COVERAGE_DIR)/coverage.html"
	@$(GOCOVER) -func=$(COVERAGE_DIR)/coverage.out

# Run linter
lint:
	@echo "Running linter..."
	$(GOLINT) run ./...

# Run vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Tidy modules
tidy:
	@echo "Tidying modules..."
	$(GOMOD) tidy

# Vendor dependencies
vendor:
	@echo "Vendoring dependencies..."
	$(GOMOD) vendor

# Full code quality check
check: fmt vet lint test

# Generate documentation
doc:
	@echo "Generating documentation..."
	@echo "This feature is not yet implemented"
	# Future: integrate with godoc or other doc generation tools

# Security checks
security-check:
	@echo "Running security checks..."
	@which gosec > /dev/null || go install github.com/securego/gosec/v2/cmd/gosec@latest
	@gosec -quiet ./...
	@which govulncheck > /dev/null || go install golang.org/x/vuln/cmd/govulncheck@latest
	@govulncheck ./...

# Show help
help:
	@echo "Bishoujo-Huntress Makefile"
	@echo "Version: $(VERSION)"
	@echo ""
	@echo "Available targets:"
	@echo "  all            - Clean, install dependencies, lint, test and build"
	@echo "  build          - Build the binary"
	@echo "  clean          - Remove build artifacts"
	@echo "  deps           - Install dependencies"
	@echo "  examples       - Build example programs"
	@echo "  test           - Run tests"
	@echo "  test-race      - Run tests with race detection"
	@echo "  test-cover     - Run tests with coverage"
	@echo "  lint           - Run linter"
	@echo "  vet            - Run go vet"
	@echo "  fmt            - Format code"
	@echo "  tidy           - Tidy modules"
	@echo "  vendor         - Vendor dependencies"
	@echo "  check          - Run all code quality checks"
	@echo "  doc            - Generate documentation"
	@echo "  security-check - Run security checks"
	@echo "  help           - Show this help message"

# Default target
.DEFAULT_GOAL := help
