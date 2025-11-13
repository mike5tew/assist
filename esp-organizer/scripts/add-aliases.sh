#!/bin/bash

# Create aliases for common tasks
echo "# ESP Organizer aliases" >> ~/.bashrc
echo "alias esp-fix-dimensions='cd /Users/michaelstewart/Coding/assist/esp-organizer && ./scripts/reset-dimensions.sh'" >> ~/.bashrc
echo "alias esp-monitor='cd /Users/michaelstewart/Coding/assist/esp-organizer && go run cmd/pipeline-monitor/main.go'" >> ~/.bashrc
echo "alias esp-reset='cd /Users/michaelstewart/Coding/assist/esp-organizer && go run cmd/weaviate-reset/main.go'" >> ~/.bashrc

# Make the aliases available in current session
source ~/.bashrc

echo "Aliases added. You can now use:"
echo "  esp-fix-dimensions    - Fix vector dimension mismatches"
echo "  esp-monitor          - Monitor the pipeline status"
echo "  esp-reset            - Reset specific Weaviate classes"
