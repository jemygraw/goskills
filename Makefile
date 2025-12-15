.PHONY: help all build clean cli runner test test-race test-coverage test-verbose lint fmt vet check deps tidy check install-tools benchmark

# Variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet
BINARY_CLI=goskills-cli
BINARY_RUNNER=goskills
BUILD_DIR=.
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Colors for terminal output
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m

# Default target
all: check test build

## help: Display this help message
help:
	@echo "$(COLOR_BOLD)GoSkills - Makefile Commands$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BOLD)Usage:$(COLOR_RESET)"
	@echo "  make $(COLOR_GREEN)<target>$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BOLD)Available targets:$(COLOR_RESET)"
	@grep -E '^## ' Makefile | sed 's/## /  $(COLOR_GREEN)/' | sed 's/:/ $(COLOR_RESET)-/'
	@echo ""

## build: Build the project
build: cli runner
	@echo "$(COLOR_GREEN)Build complete$(COLOR_RESET)"

## cli: Build CLI binary
cli:
	@echo "$(COLOR_BLUE)Building CLI...$(COLOR_RESET)"
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_CLI) ./cmd/goskills-cli

## runner: Build runner binary
runner:
	@echo "$(COLOR_BLUE)Building runner...$(COLOR_RESET)"
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_RUNNER) ./cmd/goskills

## test: Run all tests
test:
	@echo "$(COLOR_BLUE)Running tests...$(COLOR_RESET)"
	$(GOTEST) -v ./...

## test-race: Run tests with race detector
test-race:
	@echo "$(COLOR_BLUE)Running tests with race detector...$(COLOR_RESET)"
	$(GOTEST) -race ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "$(COLOR_BLUE)Running tests with coverage...$(COLOR_RESET)"
	$(GOTEST) -coverprofile=$(COVERAGE_FILE) -covermode=atomic `go list ./... | grep -v -e 'cmd/'`
	@echo "$(COLOR_GREEN)Coverage report generated: $(COVERAGE_FILE)$(COLOR_RESET)"
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "$(COLOR_GREEN)HTML coverage report: $(COVERAGE_HTML)$(COLOR_RESET)"
	go-cover-treemap -coverprofile $(COVERAGE_FILE) > coverage.svg
	@echo "$(COLOR_GREEN)SVG coverage report: coverage.svg$(COLOR_RESET)"

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "$(COLOR_BLUE)Running tests (verbose)...$(COLOR_RESET)"
	$(GOTEST) -v -count=1 ./...

## benchmark: Run benchmarks
benchmark:
	@echo "$(COLOR_BLUE)Running benchmarks...$(COLOR_RESET)"
	$(GOTEST) -bench=. -benchmem ./...

## lint: Run golangci-lint
lint:
	@echo "$(COLOR_BLUE)Running linter...$(COLOR_RESET)"
	@which golangci-lint > /dev/null || (echo "$(COLOR_YELLOW)golangci-lint not found. Run 'make install-tools'$(COLOR_RESET)" && exit 1)
	golangci-lint run ./...

## fmt: Format all Go files
fmt:
	@echo "$(COLOR_BLUE)Formatting code...$(COLOR_RESET)"
	$(GOFMT) -s -w .
	@echo "$(COLOR_GREEN)Code formatted successfully$(COLOR_RESET)"

## fmt-check: Check if code is formatted
fmt-check:
	@echo "$(COLOR_BLUE)Checking code formatting...$(COLOR_RESET)"
	@test -z "$$($(GOFMT) -l .)" || (echo "$(COLOR_YELLOW)The following files need formatting:$(COLOR_RESET)" && $(GOFMT) -l . && exit 1)
	@echo "$(COLOR_GREEN)All files are properly formatted$(COLOR_RESET)"

## vet: Run go vet
vet:
	@echo "$(COLOR_BLUE)Running go vet...$(COLOR_RESET)"
	$(GOVET) ./...

## check: Run fmt-check, vet, and lint
check: fmt-check vet lint
	@echo "$(COLOR_GREEN)All checks passed!$(COLOR_RESET)"

## clean: Clean build artifacts and test cache
clean:
	@echo "$(COLOR_BLUE)Cleaning...$(COLOR_RESET)"
	$(GOCLEAN)
	rm -f $(BUILD_DIR)/$(BINARY_CLI)
	rm -f $(BUILD_DIR)/$(BINARY_RUNNER)
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@echo "$(COLOR_GREEN)Clean complete$(COLOR_RESET)"

## deps: Download dependencies
deps:
	@echo "$(COLOR_BLUE)Downloading dependencies...$(COLOR_RESET)"
	$(GOMOD) download
	@echo "$(COLOR_GREEN)Dependencies downloaded$(COLOR_RESET)"

## tidy: Tidy and verify dependencies
tidy:
	@echo "$(COLOR_BLUE)Tidying dependencies...$(COLOR_RESET)"
	$(GOMOD) tidy
	$(GOMOD) verify
	@echo "$(COLOR_GREEN)Dependencies tidied$(COLOR_RESET)"

## install-tools: Install development tools
install-tools:
	@echo "$(COLOR_BLUE)Installing development tools...$(COLOR_RESET)"
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@echo "$(COLOR_GREEN)Tools installed$(COLOR_RESET)"

## ci: Run continuous integration checks
ci: deps check test-race test-coverage
	@echo "$(COLOR_GREEN)CI checks passed!$(COLOR_RESET)"

## pre-commit: Run pre-commit checks (fmt, vet, lint, test)
pre-commit: fmt vet lint test
	@echo "$(COLOR_GREEN)Pre-commit checks passed!$(COLOR_RESET)"

## version: Display Go version
version:
	@$(GOCMD) version

## info: Display project information
info:
	@echo "$(COLOR_BOLD)Project Information$(COLOR_RESET)"
	@echo "  Name: GoSkills"
	@echo "  Go Version: $$(go version | cut -d' ' -f3)"
	@echo "  Packages: $$(find . -name '*.go' -not -path './vendor/*' | xargs dirname | sort -u | wc -l | tr -d ' ')"
	@echo "  Lines of Code: $$(find . -name '*.go' -not -path './vendor/*' | xargs wc -l | tail -1 | awk '{print $$1}')"