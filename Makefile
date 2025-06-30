.PHONY: help build up down logs clean test

# Default target
help:
	@echo "ASSOS YouTube Automation Platform"
	@echo "Available commands:"
	@echo "  make build     - Build all services"
	@echo "  make up        - Start all services"
	@echo "  make down      - Stop all services"  
	@echo "  make logs      - View logs from all services"
	@echo "  make clean     - Clean up containers and volumes"
	@echo "  make test      - Run all tests"
	@echo "  make dev       - Start development environment"

# Build all services
build:
	docker-compose build

# Start all services
up:
	docker-compose up -d

# Stop all services
down:
	docker-compose down

# Start development environment (services only)
dev:
	docker-compose up -d postgres redis minio qdrant nats

# View logs
logs:
	docker-compose logs -f

# View logs for specific service
logs-%:
	docker-compose logs -f $*

# Clean up everything
clean:
	docker-compose down -v --remove-orphans
	docker system prune -f

# Run tests
test:
	@echo "Running Go tests..."
	cd services/api-gateway && go test ./...
	@echo "Running Rust tests..."
	cd services/video-processor && cargo test
	@echo "Running Python tests..."
	cd services/ai-service && python -m pytest

# Database operations
db-reset:
	docker-compose down postgres
	docker volume rm assos_postgres_data
	docker-compose up -d postgres

# Monitor services
status:
	docker-compose ps

# Service-specific commands
api-logs:
	docker-compose logs -f api-gateway

ai-logs:
	docker-compose logs -f ai-service

video-logs:
	docker-compose logs -f video-processor

frontend-logs:
	docker-compose logs -f frontend