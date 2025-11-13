#!/bin/bash

echo "🔄 Fixing Weaviate authentication issues..."
echo "🛑 Stopping Weaviate container..."
docker stop esp_weaviate

echo "🗑️ Removing Weaviate container..."
docker rm esp_weaviate

echo "🗑️ Removing Weaviate volume to ensure clean config..."
docker volume rm esp-organizer_weaviate_data || true

echo "🚀 Starting Weaviate with anonymous authentication enabled..."
docker-compose -f docker-compose.yml up -d weaviate

echo "⏳ Waiting for Weaviate to initialize (30s)..."
sleep 30

echo "🔍 Checking Weaviate status..."
curl -s http://localhost:8081/v1/.well-known/ready

echo "✅ Weaviate should now be running with anonymous authentication"
echo "💡 Try running: go run cmd/setup-hsg-schema/main.go"
