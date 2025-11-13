#!/bin/bash

# Download PDFs from Adobe Document Cloud share links
# Usage: ./download_adobe_pdfs.sh <adobe_share_link>

SHARE_LINK="$1"

if [ -z "$SHARE_LINK" ]; then
    echo "Usage: $0 <adobe_share_link>"
    exit 1
fi

# Extract file ID from Adobe share link
# Example link: https://acrobat.adobe.com/link/review?uri=urn:aaid:scds:US:abcd1234
FILE_ID=$(echo "$SHARE_LINK" | grep -oP '(?<=uri=)[^&]+')

if [ -z "$FILE_ID" ]; then
    echo "❌ Could not extract file ID from link"
    exit 1
fi

# Create output directory
mkdir -p resources/immunology/case_studies

# Download the PDF
OUTPUT_FILE="resources/immunology/case_studies/$(date +%Y%m%d_%H%M%S).pdf"

echo "📥 Downloading from Adobe Cloud..."
curl -L "$SHARE_LINK" \
    -H "User-Agent: Mozilla/5.0" \
    --output "$OUTPUT_FILE"

if [ -f "$OUTPUT_FILE" ]; then
    echo "✅ Downloaded to: $OUTPUT_FILE"
    echo "📊 File size: $(du -h "$OUTPUT_FILE" | cut -f1)"
else
    echo "❌ Download failed"
    exit 1
fi
