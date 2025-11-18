package main

import (
	"esp-organizer/internal/domain/api"
	"fmt"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	// Create a mux router, which is used by the main server
	router := mux.NewRouter()

	// Register all API routes
	api.RegisterRoutes(router) // This function should register all your app's routes

	fmt.Println("All registered routes:")
	fmt.Println("=======================")

	var registeredPaths []string
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			// If no methods are specified, it might be a subrouter
			methods = []string{"SUBROUTER"}
		}
		fmt.Printf("- %-10s %s\n", methods[0], path)
		registeredPaths = append(registeredPaths, fmt.Sprintf("%s %s", methods[0], path))
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking routes: %v\n", err)
		os.Exit(1)
	}

	// Check specifically for diagnostics routes
	diagnosticsRoutes := []string{
		"GET /api/diagnostics/config",
		"GET /api/diagnostics/document",
		"GET /api/diagnostics/collections",
		"GET /api/sources",
		"POST /api/immunology/upload-chapter",
		"GET /api/immunology/chapters",
	}

	fmt.Println("\nChecking for Key Routes:")
	fmt.Println("=============================")

	allFound := true
	for _, expectedRoute := range diagnosticsRoutes {
		found := false
		for _, registeredPath := range registeredPaths {
			if registeredPath == expectedRoute {
				found = true
				break
			}
		}
		if found {
			fmt.Printf("✅ %s is registered\n", expectedRoute)
		} else {
			fmt.Printf("❌ %s is NOT registered\n", expectedRoute)
			allFound = false
		}
	}

	if !allFound {
		fmt.Println("\nSome routes are missing. Check internal/api/routes.go or similar.")
		os.Exit(1)
	}

	fmt.Println("\nAll key routes are registered correctly.")
	os.Exit(0)
}
