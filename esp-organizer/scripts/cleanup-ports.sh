#!/bin/bash
# scripts/cleanup-ports.sh
# This script stops processes listening on specified ports.

set -e
echo "🧹 Cleaning up ports..."

PORTS=(8090 3000 8080 8081 27017)
for port in "${PORTS[@]}"; do
    echo "   - Checking port $port..."
    PIDS=$(lsof -ti:"$port" || true)
    if [ ! -z "$PIDS" ]; then
        echo "     Killing processes on port $port: $PIDS"
        kill -9 "$PIDS" 2>/dev/null || true
    else
        echo "     Port $port is clear."
    fi
done

echo "✅ Port cleanup complete."
