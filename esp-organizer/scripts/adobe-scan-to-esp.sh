#!/bin/bash

# Adobe Scan → ESP Organizer Automated Workflow
# Usage: ./adobe-scan-to-esp.sh <adobe_links_file>

LINKS_FILE="${1:-resources/immunology/adobe_scan_links.txt}"

if [ ! -f "$LINKS_FILE" ]; then
    echo "❌ Links file not found: $LINKS_FILE"
    echo "Usage: $0 <path_to_links_file>"
    exit 1
fi

echo "📱 Adobe Scan → ESP Organizer Workflow"
echo "======================================"
echo ""
echo "📂 Links file: $LINKS_FILE"
echo ""

# Count how many links we have
LINK_COUNT=$(grep -c "^https://" "$LINKS_FILE")
echo "📊 Found $LINK_COUNT PDF links"
echo ""

# Step 1: Download all PDFs
echo "📥 Step 1: Downloading PDFs from Adobe Document Cloud..."
make download-adobe-pdfs

if [ $? -ne 0 ]; then
    echo "❌ Download failed. Check your Adobe share links."
    exit 1
fi

echo "✅ Downloads complete!"
echo ""

# Step 2: Start ESP services
echo "🚀 Step 2: Starting ESP Organizer services..."
make start-dev

if [ $? -ne 0 ]; then
    echo "❌ Failed to start ESP services."
    exit 1
fi

echo "✅ Services started!"
echo ""

# Step 3: Process all PDFs
echo "⚙️  Step 3: Processing $LINK_COUNT PDFs with AWS Textract..."
echo "⏱️  Estimated time: ~$(($LINK_COUNT * 2)) minutes"
echo ""

make process-all-case-studies

if [ $? -ne 0 ]; then
    echo "⚠️  Some PDFs may have failed processing."
    echo "Check logs with: make logs-api"
fi

echo ""
echo "✅ Workflow complete!"
echo ""
echo "📊 Processing summary:"
make investigate-mongo

echo ""
echo "🔍 Test your knowledge graph:"
echo "make test-hsg-search QUERY=\"X-linked agammaglobulinemia\""
