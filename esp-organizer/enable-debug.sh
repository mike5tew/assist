#!/bin/bash

echo "Setting up debug environment..."

# Force environment variables
export DEBUG_MODE=true
export DEBUG_LOG_TO_FILE=true
export DEBUG_LOG_PATH=./debug.log
export DEBUG_VERBOSE=true

# Show what we're setting
echo "Environment variables set:"
echo "DEBUG_MODE=$DEBUG_MODE"
echo "DEBUG_LOG_TO_FILE=$DEBUG_LOG_TO_FILE" 
echo "DEBUG_LOG_PATH=$DEBUG_LOG_PATH"
echo "DEBUG_VERBOSE=$DEBUG_VERBOSE"

# Run the server with debug logging enabled
echo "Starting server with debug logging enabled..."
go run cmd/server/main.go
