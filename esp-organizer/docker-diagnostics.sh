#!/bin/bash

echo "Running MongoDB diagnostics via Docker..."

# Get the MongoDB container name/id
MONGO_CONTAINER=$(docker ps | grep mongo | awk '{print $1}')

if [ -z "$MONGO_CONTAINER" ]; then
    echo "No MongoDB container found running"
    exit 1
fi

# Copy the diagnostic script to the container
docker cp diagnostic.js $MONGO_CONTAINER:/tmp/diagnostic.js

# Run the script in the container
docker exec $MONGO_CONTAINER mongosh --file /tmp/diagnostic.js

echo "Diagnostics complete!"
