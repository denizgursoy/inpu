.PHONY: help test coverage coverage-html coverage-func coverage-report clean

# Default target
.DEFAULT_GOAL := help

# Variables
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

# Colors for output
GREEN  := \033[0;32m
YELLOW := \033[1;33m
NC     := \033[0m # No Color

## help: Show this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## test: Run all tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	go test -v ./...

## coverage: Generate coverage report and show in terminal
coverage:
	@echo "$(GREEN)Generating coverage report...$(NC)"
	go test -v -coverprofile=$(COVERAGE_FILE) ./...
	@echo "\n$(YELLOW)Coverage Summary:$(NC)"
	@go tool cover -func=$(COVERAGE_FILE) | grep total

## coverage-html: Generate coverage report and open in browser
coverage-html:
	@echo "$(GREEN)Generating coverage report...$(NC)"
	go test -v -coverprofile=$(COVERAGE_FILE) ./...
	@echo "$(GREEN)Opening coverage report in browser...$(NC)"
	go tool cover -html=$(COVERAGE_FILE)

## coverage-func: Show detailed coverage per function
coverage-func:
	@echo "$(GREEN)Generating coverage report...$(NC)"
	go test -v -coverprofile=$(COVERAGE_FILE) ./...
	@echo "\n$(YELLOW)Detailed Coverage:$(NC)"
	@go tool cover -func=$(COVERAGE_FILE)

## coverage-report: Generate HTML coverage report file
coverage-report:
	@echo "$(GREEN)Generating coverage report...$(NC)"
	go test -v -coverprofile=$(COVERAGE_FILE) ./...
	@echo "$(GREEN)Generating HTML report...$(NC)"
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "$(GREEN)Coverage report saved to $(COVERAGE_HTML)$(NC)"

## clean: Remove coverage files
clean:
	@echo "$(GREEN)Cleaning coverage files...$(NC)"
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)

## test-race: Run tests with race detector
test-race:
	@echo "$(GREEN)Running tests with race detector...$(NC)"
	go test -race -v ./...

## test-bench: Run benchmarks
test-bench:
	@echo "$(GREEN)Running benchmarks...$(NC)"
	go test -bench=. -benchmem ./...

test-bench-client:
	@echo "$(GREEN)Running benchmarks for clients...$(NC)"
	go test -bench=Benchmark_Client -benchmem -benchtime=2000x -run=^$

## test-all: Run tests, race detector, and coverage
test-all: test test-race coverage-html
	@echo "$(GREEN)All tests completed!$(NC)"