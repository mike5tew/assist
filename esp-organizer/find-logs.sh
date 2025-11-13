#!/bin/bash
# Simple script to help find log files

echo "Searching for log files in common locations..."

# Check current directory
echo "Checking current directory..."
find . -name "*.log" -type f -not -path "*/node_modules/*" | while read -r file; do
  echo "Found: $file"
  echo "Last 5 lines:"
  tail -n 5 "$file"
  echo "-----------------------"
done

# Check temp directory
echo "Checking temp directory..."
find /tmp -name "*debug*.log" -type f 2>/dev/null | while read -r file; do
  echo "Found: $file"
  echo "Last 5 lines:"
  tail -n 5 "$file"
  echo "-----------------------"
done

# Check home directory
echo "Checking home directory..."
find ~ -name "*debug*.log" -type f -maxdepth 2 2>/dev/null | while read -r file; do
  echo "Found: $file"
  echo "Last 5 lines:"
  tail -n 5 "$file"
  echo "-----------------------"
done

# Print current working directory for reference
echo "Current working directory: $(pwd)"

echo "Done searching for log files."
