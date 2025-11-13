#!/bin/bash

# Bulk download PDFs from Adobe Document Cloud
# Usage: ./bulk_download_adobe.sh links.txt

LINKS_FILE="$1"

if [ -z "$LINKS_FILE" ] || [ ! -f "$LINKS_FILE" ]; then
    echo "Usage: $0 <links_file>"
    echo "Create a text file with one Adobe share link per line"
    exit 1
fi

# Create output directory
mkdir -p resources/immunology/case_studies

# Read links file and download each
while IFS= read -r link; do
    # Skip empty lines and comments
    [[ -z "$link" || "$link" =~ ^# ]] && continue
    
    echo "📥 Processing: $link"
    
    # Generate filename from line number
    LINE_NUM=$((LINE_NUM + 1))
    OUTPUT_FILE="resources/immunology/case_studies/case_$(printf "%02d" $LINE_NUM).pdf"
    
    # Download with retry logic
    for attempt in 1 2 3; do
        curl -L "$link" \
            -H "User-Agent: Mozilla/5.0" \
            --connect-timeout 30 \
            --max-time 300 \
            --retry 3 \
            --output "$OUTPUT_FILE" 2>/dev/null
        
        # Check if download succeeded
        if [ -f "$OUTPUT_FILE" ] && [ -s "$OUTPUT_FILE" ]; then
            echo "✅ Downloaded: $OUTPUT_FILE ($(du -h "$OUTPUT_FILE" | cut -f1))"
            break
        else
            echo "⚠️  Attempt $attempt failed, retrying..."
            sleep 5
        fi
    done
    
    sleep 2  # Be nice to Adobe's servers
done < "$LINKS_FILE"

echo ""
echo "📊 Download Summary:"
ls -lh resources/immunology/case_studies/
