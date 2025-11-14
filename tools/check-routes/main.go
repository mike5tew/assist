package main

import (
	"esp-organizer/internal/api"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// This utility helps check if routes are properly registered
func main() {
	// Create a new router
	router := gin.New()

	// Register diagnostics routes
	api.RegisterDiagnosticsRoutes(router)

	// Print all registered routes
	routes := router.Routes()
	fmt.Println("Registered Routes:")
	fmt.Println("=================")

	for _, route := range routes {
		fmt.Printf("%s %s\n", route.Method, route.Path)
	}

	// Check specifically for diagnostics routes
	diagnosticsRoutes := []string{
		"GET /api/diagnostics/config",
		"GET /api/diagnostics/document",
		"GET /api/diagnostics/collections",
	}

	fmt.Println("\nChecking for Diagnostics Routes:")
	fmt.Println("=============================")

	for _, path := range diagnosticsRoutes {
		found := false
		for _, route := range routes {
			if fmt.Sprintf("%s %s", route.Method, route.Path) == path {
				found = true
				fmt.Printf("✅ %s is registered\n", path)
				break
			}
		}
		if !found {
			fmt.Printf("❌ %s is NOT registered\n", path)
		}
	}

	os.Exit(0)
}
