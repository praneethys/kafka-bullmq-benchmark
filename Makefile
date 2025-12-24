.PHONY: help build run test clean docker-up docker-down docker-restart deps fmt lint

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=benchmark
BUILD_DIR=./build
CMD_DIR=./cmd/benchmark
RESULTS_DIR=./results

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download Go dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod verify

build: ## Build the benchmark binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go
	@echo "Binary built at $(BUILD_DIR)/$(BINARY_NAME)"

run: build ## Build and run the benchmark
	@echo "Running benchmark..."
	$(BUILD_DIR)/$(BINARY_NAME)

run-quick: build ## Run quick benchmark (100K messages)
	@echo "Running quick benchmark..."
	$(BUILD_DIR)/$(BINARY_NAME) -messages 100000 -producers 10 -consumers 10

run-full: build ## Run full benchmark (1M messages)
	@echo "Running full benchmark..."
	$(BUILD_DIR)/$(BINARY_NAME) -messages 1000000 -producers 50 -consumers 50

run-kafka: build ## Run Kafka only benchmark
	@echo "Running Kafka benchmark..."
	$(BUILD_DIR)/$(BINARY_NAME) -queue kafka -messages 500000

run-redis: build ## Run Redis only benchmark
	@echo "Running Redis benchmark..."
	$(BUILD_DIR)/$(BINARY_NAME) -queue redis -messages 500000

docker-up: ## Start Docker infrastructure (Kafka, Redis)
	@echo "Starting Docker services..."
	docker-compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 10
	docker-compose ps

docker-down: ## Stop Docker infrastructure
	@echo "Stopping Docker services..."
	docker-compose down

docker-restart: docker-down docker-up ## Restart Docker infrastructure

docker-clean: ## Stop Docker and remove volumes
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker system prune -f

docker-logs: ## Show Docker logs
	docker-compose logs -f

test: ## Run tests
	@echo "Running tests..."
	go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

lint: ## Run linter
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

clean: ## Clean build artifacts and results
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -rf $(RESULTS_DIR)/*.json
	rm -rf $(RESULTS_DIR)/*.csv
	rm -f coverage.out coverage.html

install: build ## Install binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	go install $(CMD_DIR)/main.go

mod-tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	go mod tidy

mod-update: ## Update all dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

results-clean: ## Clean only results directory
	@echo "Cleaning results..."
	rm -rf $(RESULTS_DIR)/*.json
	rm -rf $(RESULTS_DIR)/*.csv

results-view: ## View latest results
	@echo "Latest results:"
	@ls -lt $(RESULTS_DIR) | head -5

# Development targets
dev: docker-up ## Start development environment
	@echo "Development environment ready!"
	@echo "Kafka UI: http://localhost:8080"
	@echo "Redis Commander: http://localhost:8081"

dev-down: docker-down ## Stop development environment

# Benchmark suite
benchmark-suite: docker-up ## Run complete benchmark suite
	@echo "Running complete benchmark suite..."
	@mkdir -p $(RESULTS_DIR)
	@echo "Test 1: Small messages (1KB)..."
	$(BUILD_DIR)/$(BINARY_NAME) -messages 500000 -size 1024 -producers 50 -consumers 50
	@sleep 5
	@echo "Test 2: Medium messages (5KB)..."
	$(BUILD_DIR)/$(BINARY_NAME) -messages 200000 -size 5120 -producers 30 -consumers 30
	@sleep 5
	@echo "Test 3: Large messages (10KB)..."
	$(BUILD_DIR)/$(BINARY_NAME) -messages 100000 -size 10240 -producers 20 -consumers 20
	@echo "Benchmark suite completed!"
