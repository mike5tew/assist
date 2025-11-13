#!/bin/bash

# Usage: ./list-structure.sh [path] [depth]
ROOT_PATH="${1:-.}"
MAX_DEPTH="${2:-5}"

tree -L "$MAX_DEPTH" \
  -I 'node_modules|.git|__pycache__|.env|dist|build|.next|.nuxt' \
  -a \
  "$ROOT_PATH"

# If tree is not installed, use alternative:
# find "$ROOT_PATH" -maxdepth "$MAX_DEPTH" -print | sed -e 's;[^/]*/;|____;g;s;____|;  |;g'
