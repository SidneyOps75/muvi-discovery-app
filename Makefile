# Makefile for Muvi Discovery App

# Variables
BINARY_NAME=muvi-discovery-app
MAIN_PATH=cmd/main.go
BUILD_DIR=build

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)

# Run the application
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	@go run $(MAIN_PATH)

# Run with live reload (requires air: go install github.com/cosmtrek/air@latest)
.PHONY: dev
dev:
	@echo "Starting development server with live reload..."
	@air

# Test the application
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Linting code..."
	@golangci-lint run

# Tidy dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Build for production
.PHONY: build-prod
build-prod:
	@echo "Building for production..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download

# Setup development environment
.PHONY: setup
setup: deps
	@echo "Setting up development environment..."
	@cp .env.example .env
	@mkdir -p data
	@echo "Setup complete! Please edit .env with your API keys."

# Docker build
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME) .

# Docker run
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(BINARY_NAME)

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  run        - Run the application"
	@echo "  dev        - Run with live reload (requires air)"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  fmt        - Format code"
	@echo "  lint       - Lint code (requires golangci-lint)"
	@echo "  tidy       - Tidy dependencies"
	@echo "  build-prod - Build for production"
	@echo "  build-all  - Build for multiple platforms"
	@echo "  deps       - Install dependencies"
	@echo "  setup      - Setup development environment"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run - Run Docker container"
	@echo "  help       - Show this help message"