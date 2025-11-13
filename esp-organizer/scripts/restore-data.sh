#!/bin/bash

# For use on Vultr deployment
echo "=== ESP Organizer Data Restoration ==="

if [ -z "$1" ]; then
    echo "Usage: $0 <backup-file.tar.gz>"
    echo "Example: $0 esp-organizer-data-20241202_1430.tar.gz"
    exit 1
fi

BACKUP_FILE=$1

if [ ! -f "./backups/$BACKUP_FILE" ]; then
    echo "❌ Backup file not found: ./backups/$BACKUP_FILE"
    exit 1
fi

echo "📦 Extracting backup: $BACKUP_FILE"
tar -xzf "./backups/$BACKUP_FILE"

echo "🔄 Restoring MongoDB..."
# Find latest MongoDB dump
LATEST_DUMP=$(ls -t ./backups/mongodb/dump_* | head -1)
if [ -n "$LATEST_DUMP" ]; then
    docker exec esp_mongodb mongorestore --authenticationDatabase admin \
        -u ${MONGO_INITDB_ROOT_USERNAME} \
        -p ${MONGO_INITDB_ROOT_PASSWORD} \
        --drop /backup/$(basename $LATEST_DUMP)
    echo "✅ MongoDB restored from $LATEST_DUMP"
else
    echo "⚠️  No MongoDB dump found"
fi

echo "🔄 Restoring Weaviate..."
# Find latest Weaviate backup
LATEST_WEAVIATE=$(ls -t ./backups/weaviate/weaviate_* | head -1)
if [ -n "$LATEST_WEAVIATE" ]; then
    docker exec esp_weaviate cp -r /backup/$(basename $LATEST_WEAVIATE)/* /var/lib/weaviate/
    echo "✅ Weaviate restored from $LATEST_WEAVIATE"
else
    echo "⚠️  No Weaviate backup found"
fi

echo "🚀 Data restoration complete!"
echo "🔄 Restart containers to ensure clean state:"
echo "   docker-compose restart"
