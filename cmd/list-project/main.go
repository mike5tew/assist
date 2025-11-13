package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	rootPath := flag.String("path", ".", "Root path to list")
	maxDepth := flag.Int("depth", 10, "Maximum directory depth to traverse")
	excludePatterns := flag.String("exclude", ".git,.env,node_modules,__pycache__,.venv,dist,build,.next,.nuxt", "Comma-separated patterns to exclude")
	showHidden := flag.Bool("hidden", false, "Show hidden files (starting with .)")
	flag.Parse()

	excludeMap := make(map[string]bool)
	for _, pattern := range strings.Split(*excludePatterns, ",") {
		excludeMap[strings.TrimSpace(pattern)] = true
	}

	fmt.Printf("📁 Project Structure: %s\n", *rootPath)
	fmt.Println("=" + strings.Repeat("=", len(*rootPath)+20))
	listDirectory(*rootPath, "", 0, *maxDepth, excludeMap, *showHidden)
}

func listDirectory(path string, prefix string, depth int, maxDepth int, excludeMap map[string]bool, showHidden bool) error {
	if depth > maxDepth {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return err
	}

	// Sort entries
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	// Filter entries
	var filteredEntries []fs.DirEntry
	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files if not showing them
		if !showHidden && strings.HasPrefix(name, ".") && name != ".env" && name != ".gitignore" {
			continue
		}

		// Skip excluded patterns
		skip := false
		for pattern := range excludeMap {
			if name == pattern || strings.HasPrefix(name, pattern) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		filteredEntries = append(filteredEntries, entry)
	}

	for i, entry := range filteredEntries {
		isLast := i == len(filteredEntries)-1
		currentPrefix := "├── "
		nextPrefix := "│   "

		if isLast {
			currentPrefix = "└── "
			nextPrefix = "    "
		}

		if entry.IsDir() {
			fmt.Printf("%s%s📁 %s/\n", prefix, currentPrefix, entry.Name())
			subPath := filepath.Join(path, entry.Name())
			listDirectory(subPath, prefix+nextPrefix, depth+1, maxDepth, excludeMap, showHidden)
		} else {
			icon := getFileIcon(entry.Name())
			fmt.Printf("%s%s%s %s\n", prefix, currentPrefix, icon, entry.Name())
		}
	}

	return nil
}

func getFileIcon(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	icons := map[string]string{
		".go":         "🔵",
		".ts":         "🔷",
		".tsx":        "⚛️ ",
		".js":         "🟨",
		".jsx":        "⚛️ ",
		".py":         "🐍",
		".json":       "📋",
		".yaml":       "⚙️ ",
		".yml":        "⚙️ ",
		".md":         "📝",
		".env":        "🔐",
		".gitignore":  "🚫",
		".pdf":        "📄",
		".txt":        "📄",
		".csv":        "📊",
		".sql":        "💾",
		".sh":         "⚙️ ",
		".dockerfile": "🐳",
		".compose":    "🐳",
	}

	if icon, exists := icons[ext]; exists {
		return icon
	}

	// Check full filename too
	if icon, exists := icons[filename]; exists {
		return icon
	}

	return "📄"
}
