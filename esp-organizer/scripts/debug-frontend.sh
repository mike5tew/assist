#!/bin/bash

echo "🔍 ESP Organizer Frontend Debug Script"
echo "======================================"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker Desktop."
    exit 1
fi

echo "✅ Docker is running"

# Check if compose.yaml exists
if [ ! -f "compose.yaml" ]; then
    echo "❌ compose.yaml not found. Please run from project root."
    exit 1
fi

echo "✅ Found compose.yaml"

# Handle port 8080 conflict automatically
echo ""
echo "🔍 Checking for port conflicts..."

# Check what's using port 8080
if lsof -i :8080 > /dev/null 2>&1; then
    echo "❌ Port 8080 is in use - killing conflicting processes..."
    lsof -ti:8080 | xargs kill -9 2>/dev/null || true
    sleep 2
    echo "✅ Port 8080 cleared"
else
    echo "✅ Port 8080 is available"
fi

# Check port 3000
if lsof -i :3000 > /dev/null 2>&1; then
    echo "❌ Port 3000 is in use - killing conflicting processes..."
    lsof -ti:3000 | xargs kill -9 2>/dev/null || true
    sleep 2
    echo "✅ Port 3000 cleared"
else
    echo "✅ Port 3000 is available"
fi

# Stop all services first
echo ""
echo "🛑 Stopping all services to clear any issues..."
docker-compose down --remove-orphans

# Rebuild frontend to fix nginx.conf issue
echo ""
echo "🏗️  Rebuilding frontend to fix configuration issues..."
docker-compose build --no-cache frontend

# Start services in order
echo ""
echo "🚀 Starting services in order..."

# Start databases first
echo "1. Starting databases..."
docker-compose up -d mongodb weaviate
sleep 10

# Check databases
if docker-compose ps | grep -q "mongodb.*Up"; then
    echo "✅ MongoDB is running"
else
    echo "❌ MongoDB failed to start"
fi

if docker-compose ps | grep -q "weaviate.*Up"; then
    echo "✅ Weaviate is running"
else
    echo "❌ Weaviate failed to start"
fi

# Start API
echo ""
echo "2. Starting API service..."
docker-compose up -d api
sleep 10

if docker-compose ps | grep -q "api.*Up"; then
    echo "✅ API is running"
else
    echo "❌ API failed to start"
    echo "📋 API logs:"
    docker-compose logs --tail=10 api
fi

# Start frontend
echo ""
echo "3. Starting frontend service..."
docker-compose up -d frontend
sleep 15

echo ""
echo "🌐 Testing service connections..."

# Test each service
# Corrected the port for the API service to 8080
services=("mongodb:27017" "weaviate:8081" "api:8080" "frontend:3000")
for service in "${services[@]}"; do
    IFS=':' read -r name port <<< "$service"
    if [ "$name" = "mongodb" ]; then
        # Special check for MongoDB
        if nc -z localhost "$port" 2>/dev/null; then
            echo "✅ $name is responding on port $port"
        else
            echo "❌ $name is not responding on port $port"
        fi
    else
        if curl -s "http://localhost:$port" > /dev/null 2>&1; then
            echo "✅ $name is responding on port $port"
        else
            echo "❌ $name is not responding on port $port"
            if [ "$name" = "frontend" ]; then
                echo "📋 Frontend logs:"
                docker-compose logs --tail=15 frontend
            elif [ "$name" = "api" ]; then
                echo "📋 API logs:"
                docker-compose logs --tail=10 api
            fi
        fi
    fi
done

echo ""
echo "📊 Final service status:"
docker-compose ps

echo ""
if curl -s "http://localhost:3000" > /dev/null; then
    echo "✅ Frontend is responding on port 3000"
    echo "🚀 Opening browser..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        open http://localhost:3000
    else
        echo "Frontend ready at: http://localhost:3000"
    fi
else
    echo "❌ Frontend still not responding"
    echo ""
    echo "🔧 Additional troubleshooting steps:"
    echo "   1. Check frontend build logs: docker-compose logs frontend"
    echo "   2. Rebuild without cache: docker-compose build --no-cache frontend"
    echo "   3. Check Docker resources (memory/CPU/disk space)"
    echo "   4. Try manual frontend start: docker-compose up frontend"
fi

echo ""
echo "🛠️  Useful debugging commands:"
echo "   docker-compose logs frontend    # View frontend logs"
echo "   docker-compose logs api         # View API logs"
echo "   docker-compose restart frontend # Restart frontend"
echo "   docker-compose ps               # Check all services"
echo "   make restart-quick              # Full restart"