#!/bin/bash
# filepath: /Users/michaelstewart/Coding/assist/esp-organizer/scripts/dev-start.sh

# Quick development startup script
echo "🚀 ESP Organizer Development Quick Start"
echo "========================================"

# Start services
docker-compose up -d

# Wait a moment for services to start
sleep 10

# Check if we have data
if ! curl -s "http://localhost:8080/api/skills/semantic-query?q=test" | grep -q "results"; then
    echo "📊 Loading initial data..."
    go run cmd/upload-skills/main.go 2>/dev/null || echo "⚠️  Data loading skipped (Go not available)"
fi

# Open frontend
echo "🌐 Opening frontend..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    open http://localhost:3000
else
    echo "Frontend ready at: http://localhost:3000"
fi

echo "✅ Development environment ready!"