#!/bin/bash

# Process all case studies
echo "📚 Processing all 55 case studies..."

CASE_DIR="resources/immunology/case_studies"
ANSWER_DIR="resources/immunology/answers"

# Process case studies
for pdf in "$CASE_DIR"/*.pdf; do
    echo "Processing case study: $pdf"
    make run-full-test PDF="$pdf"
    
    # Wait 10 seconds between uploads to avoid overwhelming the API
    sleep 10
done

# Process answers
for pdf in "$ANSWER_DIR"/*.pdf; do
    echo "Processing answers: $pdf"
    make run-full-test PDF="$pdf"
    sleep 10
done

echo "✅ All immunology content processed!"
