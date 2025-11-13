#!/bin/bash

echo "=== Setting up persistent data volumes ==="

# Create local data directories that will be bind-mounted
mkdir -p ./data/mongodb
mkdir -p ./data/weaviate
mkdir -p ./backups/mongodb
mkdir -p ./backups/weaviate

# On macOS with Docker Desktop, we don't need to set specific ownership
# Docker Desktop handles the permissions automatically
echo "✅ Data directories created:"
echo "   📁 ./data/mongodb (persistent MongoDB data)"
echo "   📁 ./data/weaviate (persistent Weaviate data)"
echo "   📁 ./backups/ (backup storage)"

echo ""
echo "🚀 These directories will persist data between container restarts"
echo "📦 Use ./scripts/backup-data.sh to create deployment packages"
echo "💡 Docker Desktop will handle permissions automatically"
