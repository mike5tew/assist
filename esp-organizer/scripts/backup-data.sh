#!/bin/bash

# Create backup directories
mkdir -p ./backups/mongodb
mkdir -p ./backups/weaviate
mkdir -p ./backups/exports

echo "=== ESP Organizer Data Backup ==="
echo "Creating backups for Vultr deployment..."

# Backup MongoDB
echo "📦 Backing up MongoDB..."
docker exec esp_mongodb mongodump --authenticationDatabase admin \
    -u ${MONGO_INITDB_ROOT_USERNAME} \
    -p ${MONGO_INITDB_ROOT_PASSWORD} \
    --out /backup/dump_$(date +%Y%m%d_%H%M%S)

# Backup Weaviate (copy data directory)
echo "📦 Backing up Weaviate..."
docker exec esp_weaviate cp -r /var/lib/weaviate /backup/weaviate_$(date +%Y%m%d_%H%M%S)

# Create deployment package
echo "📦 Creating deployment package..."
tar -czf ./backups/esp-organizer-data-$(date +%Y%m%d_%H%M%S).tar.gz \
    ./backups/mongodb \
    ./backups/weaviate \
    ./.env \
    ./compose.yaml

echo "✅ Backup completed!"
echo "📁 Files ready for Vultr deployment:"
echo "   - MongoDB dumps: ./backups/mongodb/"
echo "   - Weaviate data: ./backups/weaviate/"
echo "   - Deployment package: ./backups/esp-organizer-data-*.tar.gz"
