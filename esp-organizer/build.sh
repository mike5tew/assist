#!/bin/bash
# build.sh

set -e  # Exit on any error

if [ "$1" = "dev" ]; then
    echo "🚀 Building DEVELOPMENT environment..."
    docker-compose down --remove-orphans 2>/dev/null || true
    docker-compose up --build
    
elif [ "$1" = "prod" ]; then
    echo "🏗️ Building PRODUCTION environment..."
    docker-compose down --remove-orphans 2>/dev/null || true
    docker-compose up --build
    
else
    echo "❌ Usage: ./build.sh [dev|prod]"
    echo "   dev  - Development environment"
    echo "   prod - Production environment"
    exit 1
fi