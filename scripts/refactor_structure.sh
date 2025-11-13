#!/bin/bash

# refactor_structure.sh
# This script reorganizes the project into a cleaner, more conventional structure.
# WARNING: This script moves a lot of files. Backup your project before running.

echo "🧹 Starting project structure refactoring..."
echo "⚠️ WARNING: This script will move hundreds of files and directories."
echo "It is STRONGLY recommended that you commit your current progress to a backup branch on GitHub first."
echo "See docs/REFACTORING_PLAN.md for instructions."
echo ""
read -p "Have you backed up your project and are you ready to continue? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Refactoring cancelled."
    exit 1
fi


cd /Users/michaelstewart/Coding/assist || exit

# --- 1. Create the new top-level directories ---
echo "-> Creating new directory structure..."
mkdir -p esp-organizer-new/cmd/server
mkdir -p esp-organizer-new/internal/api
mkdir -p esp-organizer-new/internal/aws
mkdir -p esp-organizer-new/internal/config
mkdir -p esp-organizer-new/internal/domain
mkdir -p esp-organizer-new/internal/scheduler
mkdir -p esp-organizer-new/internal/store
mkdir -p frontend
mkdir -p tools
mkdir -p docs-new

# --- 2. Move the Frontend ---
echo "-> Separating frontend..."
if [ -d "esp-organizer/frontend/src" ]; then
    mv esp-organizer/frontend/* frontend/
    rm -rf esp-organizer/frontend
fi

# --- 3. Consolidate Go Backend Code ---
echo "-> Consolidating backend Go code..."

# Move the main server entrypoint
if [ -f "esp-organizer/cmd/api-server/main.go" ]; then
    mv esp-organizer/cmd/api-server/main.go esp-organizer-new/cmd/server/
fi

# Move one-off tools
echo "-> Moving utility commands to tools/ directory..."
if [ -d "esp-organizer/cmd" ]; then
    find esp-organizer/cmd/* -maxdepth 0 -type d ! -name 'api-server' -exec mv {} tools/ \;
    rm -rf esp-organizer/cmd
fi

# Refactor the internal directory
echo "-> Refactoring internal/ directory..."
if [ -d "esp-organizer/internal" ]; then
    mv esp-organizer/internal/api esp-organizer-new/internal/
    mv esp-organizer/internal/models esp-organizer-new/internal/domain/
    mv esp-organizer/internal/InfoFlow/InfoStore/db esp-organizer-new/internal/store/
    mv esp-organizer/internal/InfoFlow/InfoOut/llm esp-organizer-new/internal/aws/
    mv esp-organizer/internal/scheduler esp-organizer-new/internal/
    # Move remaining core logic
    mv esp-organizer/internal/InfoFlow/InfoIn/* esp-organizer-new/internal/domain/
    rm -rf esp-organizer/internal
fi

# Move go.mod and go.sum
mv esp-organizer/go.mod esp-organizer-new/
mv esp-organizer/go.sum esp-organizer-new/

# --- 4. Consolidate Documentation ---
echo "-> Consolidating documentation..."
mv docs/* docs-new/
if [ -d "esp-organizer/docs" ]; then
    mv esp-organizer/docs/* docs-new/
fi
rm -rf docs
mv docs-new docs

# --- 5. Clean up and Finalize ---
echo "-> Cleaning up and finalizing..."
# Remove old esp-organizer, keep the new one
rm -rf esp-organizer
mv esp-organizer-new esp-organizer

# Move other top-level clutter into the new esp-organizer
mv backend esp-organizer/ # This seems to be a duplicate
mv InnovationHopper esp-organizer/
mv Organisation esp-organizer/
mv shared esp-organizer/

echo "✅ Refactoring complete!"
echo "Next steps:"
echo "1. Review the new structure."
echo "2. Update import paths in your Go files (your IDE can help with this)."
echo "3. Run 'go mod tidy' in 'esp-organizer' and 'tools'."
echo "4. Run 'npm install' in 'frontend'."
