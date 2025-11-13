package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func InitDefaultLogging() {

	// Ensure logs directory exists
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		fmt.Printf("Failed to create logs directory: %v\n", err)
		return
	}

	// Open log file
	logFile, err := os.OpenFile("logs/api-server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}

	// Set log output to file
	log.SetOutput(logFile) // Uncomment if using the standard log package
	fmt.Println("✅ Log output configured to write to logs/api-server.log")

}

// FindProjectRoot searches upwards from the current directory to find the project root,
// identified by the presence of a 'go.mod' file.
func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Check if go.mod exists in the current directory
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// Move up one directory
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			// Reached the filesystem root, go.mod not found
			return "", fmt.Errorf("project root with go.mod not found")
		}
		dir = parentDir
	}
}
