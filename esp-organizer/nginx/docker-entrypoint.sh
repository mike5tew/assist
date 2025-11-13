#!/bin/sh
# /nginx/docker-entrypoint.sh

set -e

# Wait for the frontend service to be available
echo "Waiting for frontend service (frontend:80)..."
while ! nc -z frontend 80; do
  sleep 1
done
echo "Frontend service is up."

# Wait for the API service to be available
echo "Waiting for API service (api:8080)..."
while ! nc -z api 8080; do
  sleep 1
done
echo "API service is up."

echo "Starting Nginx..."
# Execute the original Nginx entrypoint
exec /docker-entrypoint.sh nginx -g 'daemon off;'
