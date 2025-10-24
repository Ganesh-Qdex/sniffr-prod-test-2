# User CRUD API Makefile

.PHONY: build run test clean docker-build docker-run docker-compose-up docker-compose-down

# Build the application
build:
	go build -o user-crud main.go

# Run the application
run:
	go run main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f user-crud

# Build Docker image
docker-build:
	docker build -t user-crud .

# Run with Docker
docker-run:
	docker run -p 8080:8080 --name user-crud-container user-crud

# Start services with Docker Compose
docker-compose-up:
	docker-compose up -d

# Stop services with Docker Compose
docker-compose-down:
	docker-compose down

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Run the test script
test-api:
	chmod +x examples/test_api.sh
	./examples/test_api.sh

# Help
help:
	@echo "Available commands:"
	@echo "  build              - Build the application"
	@echo "  run                - Run the application"
	@echo "  test               - Run tests"
	@echo "  clean              - Clean build artifacts"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run with Docker"
	@echo "  docker-compose-up  - Start services with Docker Compose"
	@echo "  docker-compose-down - Stop services with Docker Compose"
	@echo "  deps               - Install dependencies"
	@echo "  fmt                - Format code"
	@echo "  lint               - Lint code"
	@echo "  test-api           - Run API test script"
	@echo "  help               - Show this help message"
