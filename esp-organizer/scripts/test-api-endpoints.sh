#!/bin/bash
# This script tests key API endpoints to verify they're working

# Set the API URL
API_URL="http://localhost:8080"
echo "Testing API endpoints at $API_URL"

# Test root endpoint
echo -e "\nTesting API root endpoint..."
curl -s $API_URL/

# Test health endpoint
echo -e "\n\nTesting health endpoint..."
curl -s $API_URL/api/health

# Test ping endpoint
echo -e "\n\nTesting ping endpoint..."
curl -s $API_URL/api/ping

# Test diagnostics config endpoint
echo -e "\n\nTesting diagnostics config endpoint..."
curl -s $API_URL/api/diagnostics/config

echo -e "\n\nAPI testing complete!"
