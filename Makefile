.PHONY: build run test clean docker-build docker-run mcp-build mcp-run help

# Build the Slack bot application
build:
	@echo "Building banking-oncall-agent..."
	go build -o bin/banking-agent cmd/bot/main.go

# Build the MCP server
mcp-build:
	@echo "Building banking-oncall MCP server..."
	go build -o bin/banking-mcp-server cmd/mcp-server/main.go

# Run the Slack bot application locally
run: build
	@echo "Starting banking-oncall-agent..."
	./bin/banking-agent

# Run the MCP server
mcp-run: mcp-build
	@echo "Starting banking-oncall MCP server..."
	./bin/banking-mcp-server

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t banking-oncall-agent:latest .

# Run with Docker Compose
docker-run:
	@echo "Starting with Docker Compose..."
	docker-compose up --build

# Development setup
dev-setup: deps
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then \
		echo "Creating .env file from template..."; \
		cp .env.example .env; \
		echo "Please edit .env with your configuration"; \
	fi

# Help
help:
	@echo "Available targets:"
	@echo "  build         Build the application"
	@echo "  run           Build and run the application"
	@echo "  test          Run tests"
	@echo "  test-coverage Run tests with coverage report"
	@echo "  clean         Clean build artifacts"
	@echo "  deps          Download and tidy dependencies"
	@echo "  fmt           Format code"
	@echo "  lint          Run linter"
	@echo "  docker-build  Build Docker image"
	@echo "  docker-run    Run with Docker Compose"
	@echo "  dev-setup     Set up development environment"
	@echo "  help          Show this help message"