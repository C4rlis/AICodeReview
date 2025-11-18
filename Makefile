# Makefile for Code Review AI Tool

.PHONY: help build run test clean docker-up docker-down

help:
	@echo "Available commands:"
	@echo "  make build        - Build binaries"
	@echo "  make run-webhook  - Run webhook listener"
	@echo "  make run-worker   - Run worker"
	@echo "  make test         - Run tests"
	@echo "  make docker-up    - Start Docker services"
	@echo "  make docker-down  - Stop Docker services"
	@echo "  make clean        - Clean build artifacts"

build:
	@echo "Building webhook listener..."
	go build -o bin/webhook-listener.exe ./cmd/webhook-listener
	@echo "Building worker..."
	go build -o bin/worker.exe ./cmd/worker
	@echo "Build complete!"

run-webhook:
	@echo "Starting webhook listener..."
	go run ./cmd/webhook-listener/main.go

run-worker:
	@echo "Starting worker..."
	go run ./cmd/worker/main.go

test:
	@echo "Running tests..."
	go test -v ./...

docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d
	@echo "Services started! RabbitMQ Management: http://localhost:15672"

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	@echo "Clean complete!"
