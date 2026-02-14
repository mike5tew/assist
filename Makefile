# Makefile for Assist Project
# Handles restarts and running of services using absolute paths

# Absolute paths
PROJECT_ROOT := /Users/michaelstewart/Coding/assist
BACKEND_DIR := $(PROJECT_ROOT)/esp-organizer
FRONTEND_DIR := $(PROJECT_ROOT)/frontend

.PHONY: help services-up services-down services-restart run-api run-ui restart-api restart-ui check

help:
	@echo "================================================================"
	@echo "Assist Project Management"
	@echo "================================================================"
	@echo "commands:"
	@echo "  make services-up      - Start Docker services (Mongo, Weaviate)"
	@echo "  make services-down    - Stop Docker services"
	@echo "  make services-restart - Restart Docker services"
	@echo ""
	@echo "  make run-api          - Run Go Backend locally"
	@echo "  make restart-api      - Kill existing and run Go Backend locally"
	@echo ""
	@echo "  make run-ui           - Run React Frontend locally"
	@echo "  make restart-ui       - Run React Frontend (Auto-reloads usually)"
	@echo ""
	@echo "  make up               - Run FULL STACK via Docker Compose (Recommended)"
	@echo "  make down             - Stop FULL STACK"
	@echo "  make check            - Check running processes"
	@echo "================================================================"

# Docker Services
services-up:
	cd $(PROJECT_ROOT) && docker-compose up -d mongodb weaviate

services-down:
	cd $(PROJECT_ROOT) && docker-compose down

services-restart:
	cd $(PROJECT_ROOT) && docker-compose restart

# Full Stack Docker (Recommended)
up:
	@echo "Stopping any existing local processes..."
	@-lsof -ti:3000 | xargs kill -9 2>/dev/null || true
	@-lsof -ti:8080 | xargs kill -9 2>/dev/null || true
	@echo "Building and starting full docker stack..."
	cd $(PROJECT_ROOT) && docker-compose up -d --build --remove-orphans
	@echo "========================================================"
	@echo "App available at: http://localhost/esp-organizer/"
	@echo "========================================================"

down:
	@echo "Stopping full stack..."
	cd $(PROJECT_ROOT) && docker-compose down --remove-orphans

logs:
	cd $(PROJECT_ROOT) && docker-compose logs -f


# Local Backend (Go)
run-api:
	@echo "Starting backend from $(BACKEND_DIR)..."
	cd $(BACKEND_DIR) && go run cmd/server/main.go

restart-api:
	@echo "Attempting to stop existing backend..."
	@-pkill -f "go run cmd/server/main.go" || true
	@-pkill -f "esp-organizer" || true
	@echo "Backend stopped. Restarting..."
	cd $(BACKEND_DIR) && go run cmd/server/main.go

# Local Frontend (React)
run-ui:
	@echo "Starting frontend from $(FRONTEND_DIR)..."
	cd $(FRONTEND_DIR) && npm start

restart-ui:
	@echo "Starting frontend..."
	cd $(FRONTEND_DIR) && npm start

# Diagnostics
check:
	@echo "Checking for running backend processes..."
	@pgrep -fl "go run cmd/server/main.go" || echo "No local backend running."
	@echo "Checking for running frontend processes..."
	@pgrep -fl "react-scripts start" || echo "No local frontend running."
	@echo "Checking Docker containers..."
	cd $(PROJECT_ROOT) && docker-compose ps

# ============================================================
# Documentation & Indexing
# ============================================================

.PHONY: index backup search

# Generate comprehensive project index (shareable overview)
index:
	@echo "Generating project index..."
	@$(PROJECT_ROOT)/scripts/generate-project-index.sh $(PROJECT_ROOT)/PROJECT_INDEX.md
	@echo ""
	@echo "Generated: $(PROJECT_ROOT)/PROJECT_INDEX.md"
	@echo "Share this file with collaborators or feed to AI assistants."

# Backup Weaviate to S3
backup:
	@echo "Backing up Weaviate to S3..."
	@$(PROJECT_ROOT)/scripts/backup-weaviate-to-s3.sh http://localhost:8081 esp-weaviate-backups

# Search documentation in Weaviate
search:
	@echo "Usage: make search q='your query'"
	@echo "Searching for: $(q)"
	@curl -sS -X POST http://localhost:8081/v1/graphql \
		-H 'Content-Type: application/json' \
		-d '{"query":"{ Get { Documentation(limit: 10, nearText: {concepts: [\"$(q)\"]}) { title file_path project } } }"}' \
		| jq '.data.Get.Documentation[] | "\(.project): \(.title) -> \(.file_path)"'
