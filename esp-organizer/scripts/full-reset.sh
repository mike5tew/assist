#!/bin/bash
# scripts/full-reset.sh
# This script provides a comprehensive restart procedure to test updates from a clean slate.

set -e # Exit immediately if a command exits with a non-zero status.

# --- Helper Functions ---
check_and_start_docker() {
    echo "🐳 Checking Docker status..."
    if ! docker info > /dev/null 2>&1; then
        echo "   - Docker daemon is not responding."
        if [ "$(uname)" = "Darwin" ]; then
            echo "   - Attempting to start Docker Desktop on macOS..."
            open -a Docker
            echo "   - Waiting for Docker to start (this may take up to 60 seconds)..."
            for i in {1..12}; do
                if docker info > /dev/null 2>&1; then
                    echo "✅ Docker is now running."
                    return 0
                fi
                sleep 5
            done
            echo "❌ Docker did not start after 60 seconds. Please start it manually and retry."
            exit 1
        else
            echo "❌ Please ensure your Docker daemon is running and try again."
            exit 1
        fi
    else
        echo "✅ Docker is already running."
    fi
}

wait_for_service() {
    local service_name=$1
    local port=$2
    echo "⏳ Waiting for $service_name to be ready on port $port..."
    while ! nc -z localhost "$port"; do
        sleep 1
    done
    echo "✅ $service_name is up and running."
}

# --- Main Workflow ---

echo "🚀 Starting Comprehensive Restart and Test Procedure..."
echo "====================================================="

# 1. Check Docker
check_and_start_docker

# 2. Clean up ports and stop all existing containers
echo ""
echo "--- Step 1: Cleaning Environment ---"
# Assuming cleanup-ports.sh is in the same directory
"$(dirname "$0")/cleanup-ports.sh"
echo "   - Stopping and removing all Docker containers, networks, and volumes..."
docker-compose down -v --remove-orphans
echo "✅ Environment cleaned."

# 3. Start fresh databases
echo ""
echo "--- Step 2: Starting Fresh Databases ---"
docker-compose up -d mongodb weaviate
wait_for_service "MongoDB" 27017
wait_for_service "Weaviate" 8081
echo "✅ Databases are ready."

# 4. Upload initial data
echo ""
echo "--- Step 3: Seeding Database with Initial Data ---"
# Use `make` to leverage existing Makefile logic for data upload
make upload
echo "✅ Data seeding complete."

# 5. Build and start all application services
echo ""
echo "--- Step 4: Building and Starting Application Services ---"
docker-compose up -d --build api frontend nginx-local
wait_for_service "API" 8080
wait_for_service "Frontend/Nginx" 3000
echo "✅ All services are running."

# 6. Run integration tests
echo ""
echo "--- Step 5: Running Integration Tests ---"
make test
echo "✅ Integration tests passed."

# 7. Final Status
echo ""
echo "--- Workflow Complete ---"
docker-compose ps
echo ""
echo "🎉 Your environment has been fully reset and is ready for testing."
echo "   - Access the application at: http://localhost:3000"
echo "   - To view logs, run: docker-compose logs -f"
