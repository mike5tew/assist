#!/bin/bash

echo "Vector Dimension Reset Utility"
echo "=============================="
echo "This script will reset Weaviate classes to fix dimension mismatches"

# Run the pipeline monitor with the fix flag
echo "Running pipeline-monitor with fix flag..."
go run cmd/pipeline-monitor/main.go -fix

# Check if reset was successful
if [ $? -eq 0 ]; then
  echo "✅ Successfully reset Weaviate classes with dimension mismatches"
else
  echo "⚠️ There was an issue running the reset"
  echo "Trying manual class reset..."
  
  # Attempt direct class reset
  go run cmd/weaviate-reset/main.go -class SubjectAreaContent
  go run cmd/weaviate-reset/main.go -class SemanticLinks
  
  echo "Reset completed. Please try your operation again."
fi

echo "Done!"
