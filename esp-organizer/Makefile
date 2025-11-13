# Makefile

# The logic for loading .env files is now handled directly within the Go applications,
# making the Makefile simpler and avoiding environment variable conflicts.

.PHONY: help clear upload debug test analyze server build clean setup dev-workflow start-databases stop-databases restart-databases check-databases dev-start status
.PHONY: frontend-fix frontend-reset frontend-fix-ajv restart-frontend
.PHONY: health-check check-ports restart-docker restart-docker-force check-docker-status wait-for-docker restart-docker-and-services
.PHONY: fix-duplicates fix-duplicates-enhanced fix-json-format detect-duplicates-json
.PHONY: test-queries test-full setup-validated
.PHONY: clear-educational clear-medical
.PHONY: search-immunology check-medical-content search-xlinked debug-vector-search fix-medical-vectors
.PHONY: build-revectorize revectorize clean
.PHONY: docker-hard-reset logs-weaviate
.PHONY: kill-port

# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
	include .env
	export
endif

up: check-docker
	@echo "🚀 Starting all services..."
	@mkdir -p ./data/mongodb ./data/weaviate
	docker-compose -f docker-compose.yml up -d
	@echo "⏳ Waiting for services to start..."
	@sleep 10
	@make check-databases

down:
	@echo "🛑 Stopping all services..."
	docker-compose -f docker-compose.yml down

kill-port:
	@if [ -z "$(PORT)" ]; then echo "Usage: make kill-port PORT=<port_number>"; exit 1; fi
	@echo "🔍 Attempting to free port $(PORT)..."
	@PID=$$(lsof -t -i:$(PORT) 2>/dev/null); \
	if [ -n "$$PID" ]; then \
		echo "   Found process with PID $$PID on port $(PORT). Stopping it..."; \
		kill -9 $$PID; \
		echo "✅ Port $(PORT) has been freed."; \
	else \
		echo "✅ Port $(PORT) is already free."; \
	fi

docker-hard-reset:
	@echo "🔥 Performing a hard reset of the Docker environment..."
	@echo "   This will stop all containers and delete their associated volumes (data will be lost)."
	@read -p "Are you sure you want to continue? (y/N) " -n 1 -r; \
	echo; \
	if [[ ! $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "Aborted."; \
		exit 1; \
	fi
	@echo "1. Stopping and removing all services and volumes defined in docker-compose..."
	@docker-compose down -v --remove-orphans
	@echo "2. Pruning any other dangling Docker resources..."
	@docker system prune -af --volumes
	@echo "✅ Hard reset complete. Your Docker environment is now clean."
	@echo "💡 You can now restart your services, for example with 'make start-databases' or 'make start-dev'."

build-revectorize:
	@echo "Building revectorize tool for Linux/amd64 architecture..."
	GOOS=linux GOARCH=amd64 go build -o bin/revectorize cmd/revectorize/main.go

revectorize: build-revectorize
	@echo "Running content revectorization by copying binary to container..."
	@docker cp bin/revectorize esp_assist_api:/tmp/revectorize
	@docker exec esp_assist_api /tmp/revectorize
	@docker exec esp_assist_api rm /tmp/revectorize

clean:
	@echo "Cleaning binaries..."
	rm -f bin/revectorize

# Check if Docker is running
check-docker:
	@echo "🐳 Checking Docker status..."
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "❌ Docker command not found"; \
		echo "💡 Install Docker Desktop from https://www.docker.com/products/docker-desktop"; \
		exit 1; \
	fi
	@if ! docker info >/dev/null 2>&1; then \
		echo "❌ Docker daemon is not responding."; \
		if [ "$$(uname)" = "Darwin" ]; then \
			echo "💡 Attempting to start Docker Desktop on macOS..."; \
			open -a Docker; \
			echo "⏳ Waiting for Docker to start... (this may take a minute)"; \
			for i in 1 2 3 4 5 6 7 8 9 10 11 12; do \
				if docker info >/dev/null 2>&1; then \
					echo "✅ Docker is now running."; \
					break; \
				fi; \
				if [ $$i -eq 12 ]; then \
					echo "❌ Docker did not start after 60 seconds."; \
					echo "💡 Please start Docker Desktop manually and try again."; \
					exit 1; \
				fi; \
				sleep 5; \
			done; \
		else \
			echo "💡 Make sure your Docker daemon is running."; \
			exit 1; \
		fi; \
	fi
	@if ! docker ps >/dev/null 2>&1; then \
		echo "❌ Docker daemon is running but can't list containers"; \
		echo "💡 There might be permission issues or Docker needs a restart"; \
		exit 1; \
	fi
	@echo "✅ Docker is running"

# Database setup and management
start-databases: check-docker
	@echo "🚀 Starting databases (MongoDB + Weaviate)..."
	@echo "   Ensuring host data directories exist for Docker bind mounts..."
	@mkdir -p ./data/mongodb ./data/weaviate
	docker-compose -f docker-compose.yml up -d mongodb weaviate
	@echo "⏳ Waiting for databases to start..."
	@echo "   Giving containers time to initialize (30 seconds)..."
	@sleep 30
	@echo "   Now checking database health..."
	@$(MAKE) check-databases-with-retry

stop-databases:
	@echo "🛑 Stopping database services..."
	docker-compose -f docker-compose.yml stop mongodb weaviate

restart-databases: stop-databases start-databases

# Enhanced database check with retry logic
check-databases-with-retry:
	@echo "=== Checking Database Status ==="
	@echo "Checking MongoDB (port 27017)..."
	@for i in 1 2 3 4 5; do \
		if nc -z localhost 27017 2>/dev/null; then \
			echo "✅ MongoDB is running on port 27017"; \
			break; \
		else \
			if [ $$i -eq 5 ]; then \
				echo "❌ MongoDB is not running on port 27017 after 5 attempts"; \
				echo "   Run: make restart-databases"; \
				exit 1; \
			else \
				echo "⏳ Waiting for MongoDB... (attempt $$i/5)"; \
				sleep 10; \
			fi; \
		fi; \
	done
	@echo "Checking Weaviate (port 8081)..."
	@for i in 1 2 3 4 5; do \
		if curl -s -o /dev/null -w "%{http_code}" http://localhost:8081/v1/.well-known/ready | grep -q "200"; then \
			echo "✅ Weaviate is running on port 8081"; \
			break; \
		else \
			if [ $$i -eq 5 ]; then \
				echo "❌ Weaviate is not running on port 8081 after 5 attempts"; \
				echo "   Checking Weaviate container logs:"; \
				docker logs esp_weaviate | tail -n 20; \
				echo "   Run: make restart-databases"; \
				exit 1; \
			else \
				echo "⏳ Waiting for Weaviate... (attempt $$i/5)"; \
				sleep 10; \
			fi; \
		fi; \
	done

check-databases:
	@echo "Checking database services..."
	@nc -z localhost 27017 && echo "✅ MongoDB is running" || echo "❌ MongoDB is not running"
	@curl -s -o /dev/null -w "%{http_code}" http://localhost:8081/v1/.well-known/ready | grep -q "200" && echo "✅ Weaviate is running" || echo "❌ Weaviate is not running"

# Start all services in detached mode for development
start-dev: start-databases
	@echo "🚀 Starting API and nginx services..."
	@docker-compose -f docker-compose.yml up -d api nginx-local
	@echo "✅ API available at: http://localhost:8080"
	@echo "✅ Nginx serving static content at: http://localhost:8088"
	@echo ""
	@echo "⚠️  For local frontend development:"
	@echo "   Navigate to the 'frontend' directory and execute:"
	@echo "   npm install"
	@echo "   npm start"

# Stop all development services
stop-dev:
	@echo "🛑 Stopping all development services..."
	@docker-compose -f docker-compose.yml down

reset-weaviate-class:
	@if [ -z "$(CLASS)" ]; then echo "Usage: make reset-weaviate-class CLASS=<class_name>"; exit 1; fi
	@echo "🔄 Resetting Weaviate class: $(CLASS)"
	@go run cmd/weaviate-reset/main.go -class $(CLASS)

status:
	@echo "🔍 Checking Service Status..."
	@echo "--- Docker Containers (via 'docker ps') ---"
	@docker ps --format "table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}"
	@echo ""
	@echo "--- Network Port Checks ---"
	@make --no-print-directory check-databases

# Alternative: Start without Docker (native databases)
start-native-dev:
	@echo "🚀 Starting development without Docker..."
	@echo "⚠️  This requires MongoDB and Weaviate installed locally"
	@echo ""
	@echo "Installing with Homebrew:"
	@echo "  brew install mongodb/brew/mongodb-community"
	@echo "  brew install weaviate"
	@echo ""
	@echo "Starting services:"
	@brew services start mongodb/brew/mongodb-community || echo "MongoDB not installed via Homebrew"
	@brew services start weaviate || echo "Weaviate not installed via Homebrew"
	@echo ""
	@echo "If services started, run: make dev-real-processing-native"

# Development with native services (no Docker)
dev-real-processing-native:
	@echo "🚀 Starting API server with real processing (native services)..."
	@echo "✅ AWS Textract configured for real PDF processing"
	@echo "⚠️  Make sure MongoDB and Weaviate are running locally"
	@echo ""
	@echo "Checking local services:"
	@nc -z localhost 27017 && echo "✅ MongoDB running" || echo "❌ MongoDB not running"
	@nc -z localhost 8081 && echo "✅ Weaviate running" || echo "❌ Weaviate not running"
	@echo ""
	AWS_ENABLE_REAL_PROCESSING=true go run cmd/server/main.go

# Skip Docker and use cloud databases (if available)
dev-cloud-services:
	@echo "🌐 Development with cloud services..."
	@echo "This would use MongoDB Atlas and Weaviate Cloud"
	@echo "Update MONGODB_URI and WEAVIATE_URL in .env for cloud services"

# Fix the test-hsg-search target to include a waiting period after schema setup
test-hsg-search: start-databases
	@echo "🔍 Testing Hierarchical Semantic Graph search..."
	@echo "This will execute a search through the HSG architecture"
	@echo "First ensuring HSG schema is set up..."
	@if [ -f .env ]; then \
		set -a; source .env; set +a; \
	fi; \
	go run cmd/setup-hsg-schema/main.go

	@echo "⏳ Waiting for schema to propagate (10 seconds)..."
	@sleep 10
	@echo "✅ Schema setup complete, continuing with search"

	@if [ -z "$(QUERY)" ]; then \
		export QUERY="X-Linked Agammaglobulinemia"; \
		echo "Using default query: $$QUERY"; \
		if [ -f .env ]; then \
		 set -a; source .env; set +a; \
		fi; \
		go run cmd/test-hsg-search/main.go --query="X-Linked Agammaglobulinemia" --verbose || echo "❌ HSG search failed. Check the error message above."; \
	else \
		echo "Using query: $(QUERY)"; \
		if [ -f .env ]; then \
		 set -a; source .env; set +a; \
		fi; \
		go run cmd/test-hsg-search/main.go --query="$(QUERY)" --verbose || echo "❌ HSG search failed. Check the error message above."; \
	fi
	@echo ""
	@echo "✅ HSG search test complete"
	@echo "💡 To try different queries, use: make test-hsg-search QUERY=\"your search term\""

# Add a new seed-test-data target to populate the database with sample content
seed-test-data: start-databases
	@echo "🌱 Seeding test data for HSG search..."
	@if [ -f .env ]; then \
		set -a; source .env; set +a; \
	fi; \
	go run cmd/seed-test-data/main.go
	@echo "✅ Test data seeding complete"

# Simplified logs for Weaviate
logs-weaviate:
	@echo "Displaying Weaviate logs..."
	docker logs esp_weaviate

# Force restart the API service, ignoring health checks
restart-api-force:
	@echo "🔄 Force restarting API service..."
	@docker-compose -f docker-compose.yml stop api
	@sleep 2
	@docker-compose -f docker-compose.yml up -d api
	@echo "⏳ Waiting for API to initialize (10 seconds)..."
	@sleep 10
	@echo "✅ API service restarted"
	@echo ""
	@echo "API logs:"
	@docker logs --tail 10 esp_assist_api

# NEW: Run the entire PDF test workflow with a single Go command
run-full-test: start-dev
	@if [ "$(CLEANUP)" = "true" ]; then \
		make clean-test-data; \
	fi
	@echo "🚀 Running end-to-end PDF processing test..."
	@if [ -z "$(PDF)" ]; then \
		PDF_FOUND=$$(find . -name "*.pdf" -type f -not -path "*/node_modules/*" -not -path "*/\.*" | head -n 1); \
		if [ -z "$$PDF_FOUND" ]; then \
			echo "❌ No PDF specified and none found automatically. Usage: make run-full-test PDF=<path_to_file>"; \
			exit 1; \
		fi; \
		echo "⚠️ No PDF specified, using found PDF: $$PDF_FOUND"; \
		go run cmd/run-full-test/main.go --pdf="$$PDF_FOUND"; \
	else \
		go run cmd/run-full-test/main.go --pdf="$(PDF)"; \
	fi

# NEW: Clean up test data from the database
clean-test-data:
	@echo "🧹 Cleaning up test data from MongoDB..."
	@docker exec esp_mongodb mongosh "mongodb://admin:password123@localhost:27017/esp_organizer?authSource=admin" --eval 'db.immunology_content.deleteMany({ "source.title": "Test Book" })' >/dev/null
	@docker exec esp_mongodb mongosh "mongodb://admin:password123@localhost:27017/esp_organizer?authSource=admin" --eval 'db.extraction_jobs.deleteMany({ "_id": { $$regex: "^immunology-chapter-" } })' >/dev/null
	@echo "✅ Test data cleaned."

# DEPRECATED: Replaced by run-full-test
test-pdf-upload:
	@echo "DEPRECATED: This target is complex and has been replaced." >&2
	@echo "Please use 'make run-full-test PDF=<path>' instead." >&2
	@echo "Example: make run-full-test PDF=./resources/immunology_chapter_1.pdf" >&2
	@exit 1

# Add a new target to upload PDFs first
upload-pdf:
	@echo "📤 Uploading PDF document..."
	@if [ -z "$(PDF)" ]; then \
		echo "⚠️  No PDF specified. Usage: make upload-pdf PDF=path/to/your/file.pdf"; \
		exit 1; \
	fi
	@if [ ! -f "$(PDF)" ]; then \
		echo "❌ File not found: $(PDF)"; \
		exit 1; \
	fi
	@echo "📄 Uploading: $(PDF)"
	@if [ -f .env ]; then \
		set -a; source .env; set +a; \
	fi; \
	go run cmd/upload-pdf/main.go --file="$(PDF)" || echo "❌ PDF upload failed. Check the error message above."
	@echo ""
	@echo "✅ PDF upload complete. Use the source ID from above with: make test-pdf-upload SOURCE_ID=<id>"

# Check and fix API health issues
check-api-health:
	@echo "🔍 Checking API health status..."
	@if ! nc -z localhost 8080 2>/dev/null; then \
		echo "❌ API server is not running on port 8080"; \
		echo "   Starting API service..."; \
		docker-compose -f docker-compose.yml up -d api; \
		sleep 5; \
	else \
		echo "✅ API server is responding on port 8080"; \
	fi
	@if docker ps | grep esp_assist_api | grep -q "(unhealthy)"; then \
		echo "⚠️  API is marked as UNHEALTHY by Docker"; \
		echo "   Showing recent API logs:"; \
		docker logs --tail 20 esp_assist_api; \
		echo ""; \
		echo "🔄 Attempting to restart the API service..."; \
		docker-compose -f docker-compose.yml restart api; \
		echo "⏳ Waiting for API to start (30 seconds)..."; \
		sleep 30; \
		if docker ps | grep esp_assist_api | grep -q "(unhealthy)"; then \
			echo "❌ API is still unhealthy after restart"; \
			echo "   You may need to recreate the container: make stop-dev && make start-dev"; \
		else \
			echo "✅ API is now healthy"; \
		fi; \
	else \
		echo "✅ API container is marked as healthy"; \
	fi
	@echo ""
	@echo "Testing API health endpoint..."
	@curl -s http://localhost:8080/health || echo "❌ Health endpoint not responding"

# Add a new target to check API logs in real-time
logs-api:
	@echo "📋 Showing API container logs (press Ctrl+C to exit)..."
	@docker logs -f esp_assist_api

# NEW: Cleanly restart the API and show startup logs
restart-and-log-api:
	@echo "🔄 Restarting API container and showing startup logs..."
	@docker-compose stop api
	@docker-compose rm -f api
	@docker-compose up -d api
	@echo "⏳ Waiting for container to start..."
	@sleep 5
	@docker logs esp_assist_api

# Add a debug utility to check file status in API container
check-pdf-in-container:
	@if [ -z "$(PDF)" ]; then \
		echo "⚠️  No PDF specified. Usage: make check-pdf-in-container PDF=path/to/your/file.pdf"; \
		exit 1; \
	fi
	@echo "🔍 Checking if PDF exists on your machine..."
	@if [ ! -f "$(PDF)" ]; then \
		echo "❌ File not found locally: $(PDF)"; \
		exit 1; \
	fi
	@echo "📄 Local file $(PDF) exists, size: $$(du -h "$(PDF)" | cut -f1)"
	@echo ""
	@echo "🔍 Testing file upload to API container..."
	@FILENAME=$$(basename "$(PDF)")
	@echo "   Copying file to container temp directory..."
	@docker cp "$(PDF)" esp_assist_api:/tmp/
	@echo "   Checking if file exists in container..."
	@docker exec esp_assist_api ls -la /tmp/$$FILENAME || echo "❌ File not found in container"
	@echo ""
	@echo "🔍 API container upload directory contents:"
	@docker exec esp_assist_api ls -la /app/uploads/ 2>/dev/null || echo "❌ Upload directory not found or empty"
	@echo ""
	@echo "💡 If the file is present in the container but processing fails,"
	@echo "   the issue might be with PDF parsing or the extraction service."
	@echo "   Check the API logs with: make logs-api"

# DEPRECATED
check-test-pdf:
	@echo "DEPRECATED: Use 'make run-full-test' and specify the PDF you want to test." >&2
	@exit 1

# DEPRECATED
test-pdf-upload-debug:
	@echo "DEPRECATED: This target is complex and has been replaced." >&2
	@echo "Please use 'make run-full-test PDF=<path>' instead." >&2
	@exit 1

# DEPRECATED
direct-pdf-process: 
	@echo "DEPRECATED: This functionality is now part of 'make run-full-test'." >&2
	@exit 1

# Check AWS service credentials and configuration
check-aws-config:
	@echo "🔍 Checking AWS configuration in API container..."
	@echo "---"
	@echo "1️⃣ Verifying AWS identity..."
	@docker exec esp_assist_api sh -c 'aws sts get-caller-identity >/dev/null 2>&1 && echo "✅ AWS identity OK" || echo "❌ Failed to get AWS identity"'
	@echo "---"
	@echo "2️⃣ Checking AWS environment variables..."
	@docker exec esp_assist_api env | grep -E 'AWS|TEXTRACT|BEDROCK' || echo "No AWS environment variables found"
	@echo "---"
	@echo "3️⃣ Testing AWS S3 access..."
	@docker exec esp_assist_api sh -c 'aws s3 ls >/dev/null 2>&1 && echo "✅ S3 access OK" || echo "❌ S3 access failed"'
	@echo "---"
	@echo "4️⃣ Testing AWS Textract access..."
	@if docker exec esp_assist_api sh -c 'aws textract detect-document-text --document '\''{"S3Object":{"Bucket":"dummy","Name":"dummy"}}'\'' 2>&1' | grep -q "InvalidS3ObjectException"; then \
		echo "✅ Textract access OK (permission granted)"; \
	else \
		echo "❌ Textract access failed. Check IAM policy and command syntax."; \
	fi
	@echo "---"
	@echo "5️⃣ Testing AWS Bedrock model access..."
	@docker exec esp_assist_api sh -c 'aws bedrock list-foundation-models --region eu-west-2 >/dev/null 2>&1 && echo "✅ Bedrock model listing OK" || echo "❌ Bedrock model listing failed"'
	@echo "---"
	@echo "6️⃣ Testing AWS Bedrock model invocation (Claude)..."
	@if docker exec esp_assist_api sh -c ' \
		echo '\''{"anthropic_version":"bedrock-2023-05-31","max_tokens":10,"messages":[{"role":"user","content":[{"type":"text","text":"test"}]}]}'\'' > /tmp/payload.json && \
		aws bedrock-runtime invoke-model \
			--model-id "anthropic.claude-3-sonnet-20240229-v1:0" \
			--body file:///tmp/payload.json \
			--cli-binary-format raw-in-base64-out \
			--region eu-west-2 \
			/dev/null >/dev/null 2>&1 && \
		rm /tmp/payload.json'; then \
		echo "✅ Bedrock model invocation OK"; \
	else \
		echo "❌ Bedrock model invocation failed. Ensure Claude model access is enabled in AWS Bedrock console."; \
	fi
	@echo "---"

# View failed jobs and provide retry options
view-failed-jobs:
	@echo "🔍 Checking for failed extraction jobs..."
	@echo "1️⃣ Retrieving last 10 jobs from MongoDB..."
	@docker exec esp_mongodb sh -c 'mongosh "mongodb://${MONGO_INITDB_ROOT_USERNAME}:${MONGO_INITDB_ROOT_PASSWORD}@localhost:27017/${MONGO_DB_NAME}?authSource=${MONGO_AUTH_SOURCE}" --eval "db.extraction_jobs.find({\"status\": \"FAILED\"}).sort({\"created_at\": -1}).limit(10).pretty()"' || echo "Failed to query MongoDB"
	@echo ""
	@echo "2️⃣ To retry a failed job, use: make retry-job JOB_ID=<job_id>"

# Retry a failed job
retry-job:
	@if [ -z "$(JOB_ID)" ]; then \
		echo "⚠️  No job ID specified. Usage: make retry-job JOB_ID=<job_id>"; \
		exit 1; \
	fi
	@echo "🔄 Attempting to retry job $(JOB_ID)..."
	@echo "1️⃣ Updating job status in MongoDB..."
	@docker exec esp_mongodb sh -c "echo 'db.extraction_jobs.updateOne({\"job_id\": \"$(JOB_ID)\"}, {\$$set: {\"status\": \"PENDING\", \"updated_at\": new Date()}})' | mongosh esp_organizer" || echo "Failed to update MongoDB"
	@echo ""
	@echo "2️⃣ Triggering job processing..."
	@curl -s -X POST http://localhost:8080/api/extraction/retry/$(JOB_ID) || echo "Failed to trigger retry via API"
	@echo ""
	@echo "3️⃣ Restarting API to ensure workers pick up the job..."
	@docker-compose -f docker-compose.yml restart api
	@echo "⏳ Waiting for API to restart (10 seconds)..."
	@sleep 10
	@echo "✅ API restarted, job should be processing"
	@echo ""
	@echo "4️⃣ To monitor status: curl http://localhost:8080/api/extraction/status/$(JOB_ID)"

# Reprocess an existing job from MongoDB data
reprocess-job:
	@if [ -z "$(JOB_ID)" ]; then \
		echo "❌ Usage: make reprocess-job JOB_ID=<job_id>"; \
		echo "Example: make reprocess-job JOB_ID=immunology-chapter-1759699625"; \
		exit 1; \
	fi
	@echo "🔄 Reprocessing job $(JOB_ID) from MongoDB..."
	@if [ -f .env ]; then \
		set -a; . .env; set +a; \
	fi; \
	go run cmd/reprocess-job/main.go -job=$(JOB_ID)

# DEPRECATED
pdf-performance-test:
	@echo "DEPRECATED: Use 'make run-full-test' for end-to-end testing." >&2
	@exit 1

# Help section addition
help:
	@echo "ESP Organizer - Comprehensive Makefile Help"
	@echo "==========================================="
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "🚀 Development Environment:"
	@echo "  start-dev			  Start full development environment (API + frontend)"
	@echo "  start-native-dev	   Start development without Docker (requires local DBs)"
	@echo "  dev-real-processing-native Start API with real processing (native services)"
	@echo "  dev-cloud-services	 Development with cloud services (MongoDB Atlas + Weaviate Cloud)"
	@echo "  status				 Check service status"
	@echo ""
	@echo "🐳 Docker & Database Management:"
	@echo "  start-databases		Start MongoDB + Weaviate databases"
	@echo "  stop-databases		 Stop databases"
	@echo "  restart-databases	  Restart databases"
	@echo "  check-databases		Check database status"
	@echo "  docker-hard-reset	 Purge all Docker containers and volumes (DESTRUCTIVE)"
	@echo "  check-docker		   Check Docker status"
	@echo ""
	@echo "🔧 Data Operations & Testing:"
	@echo "  run-full-test PDF=<path>  Run the full, end-to-end PDF processing test"
	@echo "					  CLEANUP=true	 Clean up previous test data before running"
	@echo "  test-hsg-search		Test Hierarchical Semantic Graph search"
	@echo "  seed-test-data		 Populate database with test data for HSG"
	@echo "  reset-weaviate-class CLASS=<name> Reset a specific Weaviate class"
	@echo "  revectorize			Run content revectorization"
	@echo "  build-revectorize	  Build a new revectorize binary"
	@echo ""
	@echo "🔍 Debugging & Troubleshooting:"
	@echo "  logs-api			  Show live API container logs"
	@echo "  logs-weaviate		 Show live Weaviate container logs"
	@echo "  check-api-health	  Check API container health and fix if needed"
	@echo "  check-aws-config	  Check AWS credentials and service access from within the container"
	@echo "  view-failed-jobs	  List recent failed extraction jobs from MongoDB"
	@echo "  retry-job JOB_ID=<id>		 Retry a failed extraction job"
	@echo ""
	@echo "🔄 System Utilities:"
	@echo "  up					 Start all services"
	@echo "  down				   Stop all services"
	@echo "  kill-port PORT=<port>  Find and stop a process using a specific port"
	@echo "  clean				  Clean binaries"


	@echo "  logs-weaviate		  Show Weaviate logs"	@echo ""	@echo "For detailed help on a specific command, check the Makefile comments."
	@echo "Most commands support environment variables for customization."

# Investigate MongoDB collection history
investigate-mongo:
	@echo "🔍 Investigating MongoDB collections and recent operations..."
	@echo ""
	@echo "1️⃣ Listing all collections in esp_organizer database:"
	@docker exec esp_mongodb mongosh "mongodb://admin:password123@localhost:27017/esp_organizer?authSource=admin" --eval "db.getCollectionNames()" --quiet
	@echo ""
	@echo "2️⃣ Counting documents in extraction_jobs:"
	@docker exec esp_mongodb mongosh "mongodb://admin:password123@localhost:27017/esp_organizer?authSource=admin" --eval "db.extraction_jobs.countDocuments()" --quiet
	@echo ""
	@echo "3️⃣ Counting documents in immunology_content:"
	@docker exec esp_mongodb mongosh "mongodb://admin:password123@localhost:27017/esp_organizer?authSource=admin" --eval "db.immunology_content.countDocuments()" --quiet
	@echo ""
	@echo "4️⃣ Checking MongoDB oplog for recent delete operations (if replica set is enabled):"
	@docker exec esp_mongodb mongosh "mongodb://admin:password123@localhost:27017/local?authSource=admin" --eval "db.oplog.rs.find({op: 'd', ns: 'esp_organizer.extraction_jobs'}).sort({ts: -1}).limit(5).pretty()" --quiet 2>/dev/null || echo "⚠️ Oplog not available (not a replica set)"
	@echo ""
	@echo "5️⃣ Checking for any jobs matching your pattern:"
	@docker exec esp_mongodb mongosh "mongodb://admin:password123@localhost:27017/esp_organizer?authSource=admin" --eval "db.extraction_jobs.find({_id: /immunology-chapter/}).pretty()" --quiet
	@echo ""
	@echo "6️⃣ Listing ALL extraction jobs:"
	@docker exec esp_mongodb mongosh "mongodb://admin:password123@localhost:27017/esp_organizer?authSource=admin" --eval "db.extraction_jobs.find({}, {_id: 1, status: 1, created_at: 1}).sort({created_at: -1}).limit(10).pretty()" --quiet
	@echo ""
	@echo "7️⃣ Checking if CLEANUP=true was used recently (checking make command history):"
	@echo "Recent make commands from shell history:"
	@history | grep "make.*run-full-test" | tail -5 || echo "No recent make commands found"
	@echo ""
	@echo "💡 Most likely causes:"
	@echo "   • 'make run-full-test CLEANUP=true' was run (deletes test data)"
	@echo "   • 'make clean-test-data' was run manually"
	@echo "   • Manual deletion via MongoDB shell"
	@echo "   • Application code deleted the job after completion"
extract-chapter:
	@if [ -z "$(CHAPTER)" ]; then echo "Usage: make extract-chapter CHAPTER=<file>"; exit 1; fi
	@echo "📚 Extracting semantic links from: $(CHAPTER)"
	@go run cmd/extract-semantic-links/main.go \
		-chapter=$(CHAPTER) \
		-output=extracted_links_$(shell date +%Y%m%d_%H%M%S).json \
		-domain=immunology

# Validate your manually curated links against LLM extractions
validate-curated-links:
	@echo "🔍 Validating curated semantic links..."
	@go run cmd/validate-curated-links/main.go

# Quick extraction from clipboard text (for quick experiments)
extract-from-clipboard:
	@echo "📋 Extracting from clipboard..."
	@pbpaste > /tmp/clipboard_extract.txt
	@go run cmd/extract-semantic-links/main.go \
		-chapter=/tmp/clipboard_extract.txt \
		-output=clipboard_links.json \
		-domain=immunology
	@rm /tmp/clipboard_extract.txt
	@echo "✅ Results in clipboard_links.json"

# Show extraction statistics
show-stats:
	@echo "📊 Semantic Link Statistics:"
	@echo "Human-curated links:"
	@cat docs/training_data/human_extracted_links.md | grep -c "^[A-Za-z*]" || echo "0"
	@echo "LLM extractions this week:"
	@find . -name "extracted_links_*.json" -mtime -7 -exec wc -l {} \; | awk '{sum += $$1} END {print sum/2}'

# Download PDFs from Adobe Document Cloud
download-adobe-pdfs:
	@echo "📥 Downloading immunology PDFs from Adobe Cloud..."
	@go run cmd/download-adobe-pdfs/main.go \
		-links=resources/immunology/adobe_links.txt \
		-output=resources/immunology/case_studies
	@echo "✅ Downloads complete!"

# Download and process in one step
download-and-process-adobe:
	@echo "📥 Downloading and processing Adobe PDFs..."
	@make download-adobe-pdfs
	@echo "⏳ Waiting 5 seconds before starting processing..."
	@sleep 5
	@make process-all-case-studies

# Process all downloaded case studies
process-all-case-studies: start-dev
	@echo "🚀 Processing all case studies..."
	@for pdf in resources/immunology/case_studies/*.pdf; do \
		echo "📄 Processing: $$pdf"; \
		make run-full-test PDF="$$pdf"; \
		sleep 10; \
	done
	@echo "✅ All case studies processed!"