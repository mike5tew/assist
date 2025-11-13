#!/usr/bin/env python3
import os
import sys
from pathlib import Path

def print_tree(path, prefix="", max_depth=5, current_depth=0, exclude=None):
    if exclude is None:
        exclude = {'.git', 'node_modules', '__pycache__', '.env', 'dist', 'build', '.next'}
    
    if current_depth > max_depth:
        return
    
    try:
        entries = sorted(os.listdir(path))
    except PermissionError:
        return
    
    # Filter
    entries = [e for e in entries if e not in exclude]
    
    for i, entry in enumerate(entries):
        full_path = os.path.join(path, entry)
        is_last = i == len(entries) - 1
        
        current = "└── " if is_last else "├── "
        next_prefix = "    " if is_last else "│   "
        
        if os.path.isdir(full_path):
            print(f"{prefix}{current}📁 {entry}/")
            print_tree(full_path, prefix + next_prefix, max_depth, current_depth + 1, exclude)
        else:
            icon = get_icon(entry)
            print(f"{prefix}{current}{icon} {entry}")

def get_icon(filename):
    ext = os.path.splitext(filename)[1].lower()
    icons = {
        '.go': '🔵', '.ts': '🔷', '.tsx': '⚛️ ', '.js': '🟨', '.py': '🐍',
        '.json': '📋', '.yaml': '⚙️ ', '.yml': '⚙️ ', '.md': '📝', '.env': '🔐'
    }
    return icons.get(ext, '📄')

if __name__ == "__main__":
    path = sys.argv[1] if len(sys.argv) > 1 else "."
    depth = int(sys.argv[2]) if len(sys.argv) > 2 else 5
    print(f"📁 Project Structure: {path}\n")
    print_tree(path, max_depth=depth)
