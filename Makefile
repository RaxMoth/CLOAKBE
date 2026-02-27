.PHONY: help start stop db-up db-down migrate backend test clean setup

# Database config
DB_NAME=cloak_db
DB_USER=postgres
DB_PASS=postgres
DB_PORT=5432
DB_HOST=localhost
CONTAINER_NAME=postgres-cloak

help:
	@echo "=== CLOAK Backend - Quick Commands ==="
	@echo ""
	@echo "Environments:"
	@echo "  make env-local   - Setup local development environment"
	@echo "  make env-dev     - Setup development environment"
	@echo "  make env-prod    - Setup production environment"
	@echo ""
	@echo "Quick Start:"
	@echo "  make start       - Start everything (db + migrations + backend)"
	@echo "  make stop        - Stop database"
	@echo ""
	@echo "Local Testing with ngrok:"
	@echo "  make tunnel      - Start ngrok tunnel (for Flutter testing)"
	@echo "  make tunnel-help - Show ngrok workflow"
	@echo ""
	@echo "Database:"
	@echo "  make db-up       - Start PostgreSQL container"
	@echo "  make db-down     - Stop PostgreSQL container"
	@echo "  make db-shell    - Open PostgreSQL shell"
	@echo ""
	@echo "Setup & Migration:"
	@echo "  make setup-env   - Create .env file (uses .env.local by default)"
	@echo "  make migrate     - Run database migrations"
	@echo ""
	@echo "Backend:"
	@echo "  make backend     - Start API server"
	@echo "  make build       - Build binary"
	@echo ""
	@echo "Testing:"
	@echo "  make test        - Test health endpoint"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean       - Remove build artifacts"
	@echo ""

# ============ QUICK START ============

start: db-up migrate backend
	@echo "Starting everything..."

stop: db-down
	@echo "Stopped"

# ============ ENVIRONMENT SETUP ============

env-local:
	@echo "✓ Setting up local development environment"
	@cp .env.local .env 2>/dev/null || echo "Using existing .env.local"
	@echo "✓ Using .env.local configuration"
	@echo "  Database: localhost:5432"
	@echo "  Port: 8080"
	@echo "  JWT: dev secret (NOT for production)"

env-dev:
	@echo "✓ Setting up development environment"
	@cp .env.dev .env 2>/dev/null || echo "Using .env.dev"
	@echo "⚠️  Please update .env.dev with your dev server details:"
	@grep "CHANGE_ME" .env.dev || echo "✓ Configuration looks good"

env-prod:
	@echo "✓ Setting up production environment"
	@cp .env.prod .env 2>/dev/null || echo "Using .env.prod"
	@echo "⚠️  CRITICAL: Update .env.prod with secure production secrets:"
	@echo "   - JWT_SECRET: Use 'openssl rand -base64 32' for secure key"
	@echo "   - HMAC_SECRET: Use 'openssl rand -base64 32' for secure key"
	@echo "   - DATABASE_URL: Use production database credentials"
	@grep "CHANGE_ME" .env.prod

# ============ DATABASE ============

db-up:
	@echo "Starting PostgreSQL..."
	@docker run -d --name $(CONTAINER_NAME) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASS) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p $(DB_PORT):5432 \
		postgres:16 2>/dev/null || echo "Container already running"
	@echo "PostgreSQL started on localhost:$(DB_PORT)"
	@sleep 5

db-down:
	@echo "Stopping PostgreSQL..."
	@docker stop $(CONTAINER_NAME) 2>/dev/null
	@docker rm $(CONTAINER_NAME) 2>/dev/null
	@echo "PostgreSQL stopped"

db-clean:
	@echo "Removing PostgreSQL container and data..."
	@docker stop $(CONTAINER_NAME) 2>/dev/null
	@docker rm $(CONTAINER_NAME) 2>/dev/null
	@echo "PostgreSQL cleaned"

db-shell:
	@docker exec -it $(CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME)

db-logs:
	@docker logs -f $(CONTAINER_NAME)

# ============ SETUP ============

setup-env:
	@echo "Creating .env file..."
	@echo "PORT=8080" > .env
	@echo "ENVIRONMENT=development" >> .env
	@echo "DATABASE_URL=postgres://postgres:postgres@localhost:5432/cloak_db?sslmode=disable" >> .env
	@echo "JWT_SECRET=dev-secret-key-change-in-production" >> .env
	@echo "HMAC_SECRET=dev-hmac-secret-change-in-production" >> .env
	@echo ".env created!"

# ============ MIGRATIONS ============

migrate:
	@echo "Running migrations..."
	@migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up
	@echo "Migrations complete!"

migrate-down:
	@echo "Rolling back migrations..."
	@migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

# ============ BACKEND ============

backend:
	@echo "Starting CLOAK Backend on http://localhost:8080"
	@./bin/api

tunnel:
	@echo "Starting ngrok tunnel to localhost:8080..."
	@echo ""
	@echo "⚠️  Copy the HTTPS URL below and update:"
	@echo "    lib/core/constants/app_constants.dart in your Flutter app"
	@echo ""
	ngrok http 8080

tunnel-help:
	@echo "To use ngrok tunnel for local development:"
	@echo ""
	@echo "Terminal 1:"
	@echo "  make start    # Start backend + database"
	@echo ""
	@echo "Terminal 2:"
	@echo "  make tunnel   # Start ngrok tunnel"
	@echo ""
	@echo "Terminal 3:"
	@echo "  flutter run -d macos  # Start Flutter app"
	@echo ""
	@echo "Then update Flutter app with ngrok URL and hot reload (R)"

build:
	@echo "Building backend..."
	@go build -o bin/api ./cmd/api
	@echo "Build complete!"

rebuild: clean build

# ============ TESTING ============

test:
	@echo "Testing health endpoint..."
	@curl -s http://localhost:8080/health | jq . || echo "Backend not running"

test-register:
	@echo "Testing registration..."
	@curl -s -X POST http://localhost:8080/api/v1/auth/business/register \
		-H "Content-Type: application/json" \
		-d '{"name":"Test","email":"test@test.com","password":"pass123"}' | jq .

test-login:
	@echo "Testing login..."
	@curl -s -X POST http://localhost:8080/api/v1/auth/business/login \
		-H "Content-Type: application/json" \
		-d '{"email":"test@test.com","password":"pass123"}' | jq .

# ============ CLEANUP ============

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "Cleaned!"

all-down:
	@echo "Stopping everything..."
	@docker stop $(CONTAINER_NAME) 2>/dev/null
	@docker rm $(CONTAINER_NAME) 2>/dev/null
	@rm -rf bin/
	@echo "Everything stopped and cleaned"
