#!/bin/bash
# cleanup-ports.sh
echo "Cleaning up ports..."

PORTS=(8090 3000 8080 8081 27017)
for port in "${PORTS[@]}"; do
    echo "Checking port $port..."
    PIDS=$(lsof -ti:"$port")
    if [ ! -z "$PIDS" ]; then
        echo "Killing processes on port $port: $PIDS"
        kill -9 "$PIDS" 2>/dev/null
    fi
done

echo "Stopping Docker containers..."
docker-compose down --remove-orphans
echo "Cleanup complete!"