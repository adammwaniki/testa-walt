.PHONY: help run build test clean docker-build docker-run deploy

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	go run main.go

build: ## Build the application
	go build -o testa-gava main.go


clean: ## Clean build artifacts
	rm -f testa-gava
	go clean

docker-build: ## Build Docker image
	docker build -t testa-gava:latest .

docker-run: ## Run Docker container in detatched mode
	docker compose up -d

docker-stop: ## Stop Docker container
	docker compose down

docker-logs: ## View Docker logs
	docker compose logs -f
