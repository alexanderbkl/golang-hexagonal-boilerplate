.PHONY: help build run test clean docker-build docker-up docker-down migrate-up migrate-down sqlc-generate grpc-generate graphql-generate

help: ## Display this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/server ./cmd/server
	@go build -o bin/migrate ./cmd/migrate
	@echo "Build complete!"

run: ## Run the application locally
	@echo "Running application..."
	@go run ./cmd/server/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@echo "Clean complete!"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker-compose build
	@echo "Docker build complete!"

docker-up: ## Start services with Docker Compose
	@echo "Starting services..."
	@docker-compose up -d
	@echo "Services started!"
	@echo "GraphQL Playground: http://localhost:8080"
	@echo "gRPC Server: localhost:9090"

docker-down: ## Stop Docker Compose services
	@echo "Stopping services..."
	@docker-compose down
	@echo "Services stopped!"

docker-logs: ## View Docker logs
	@docker-compose logs -f app

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	@go run ./cmd/migrate/main.go up
	@echo "Migrations complete!"

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	@go run ./cmd/migrate/main.go down
	@echo "Rollback complete!"

sqlc-generate: ## Generate sqlc code
	@echo "Generating sqlc code..."
	@sqlc generate
	@echo "sqlc generation complete!"

grpc-generate: ## Generate gRPC code
	@echo "Generating gRPC code..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/grpc/user.proto
	@echo "gRPC generation complete!"

graphql-generate: ## Generate GraphQL code
	@echo "Generating GraphQL code..."
	@go run github.com/99designs/gqlgen generate
	@echo "GraphQL generation complete!"

install-tools: ## Install development tools
	@echo "Installing tools..."
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Tools installed!"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies downloaded!"
