package main

import (
	"context"
	"esp-organizer/internal/config"
	"esp-organizer/internal/domain/api"
	"esp-organizer/internal/server"
	"esp-organizer/internal/store/db"
	"esp-organizer/internal/utils"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// globalCorsHandler adds CORS headers to all responses
func globalCorsHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")

		// Special handling for OPTIONS requests
		if r.Method == "OPTIONS" {
			log.Printf("🔍 CORS: Handling OPTIONS preflight for %s", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			return
		}

		// For all other requests, pass to the next handler
		handler.ServeHTTP(w, r)
	})
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}
	// In your router setup (cmd/server/main.go or internal/routes/routes.go)

	// Configure logging
	logFile, err := os.OpenFile("api-server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		// Use a multi-writer to log to both file and stdout
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)
		log.Println("✅ Log output configured to write to api-server.log")
	} else {
		log.Printf("⚠️ Could not open log file: %v", err)
	}

	// Initialize debug system with the same log file
	utils.ForceFileLogging("api-server.log")

	// Use either the server package or direct initialization based on your preference
	useServerPackage := false

	if useServerPackage {
		// Create and start the server using the server package
		s := server.NewServer()

		// Print all registered routes for debugging
		routes := s.GetRouter().Routes()
		log.Println("=== REGISTERED ROUTES AT STARTUP ===")
		for _, route := range routes {
			log.Printf("%s %s", route.Method, route.Path)
		}
		log.Println("===================================")

		if err := s.Start(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
			os.Exit(1)
		}
	} else {
		// Direct initialization approach
		// Database connection setup
		mongoDb, err := db.NewFromEnv()
		if err != nil {
			log.Fatalf("Failed to initialize MongoDB: %v", err)
		}
		defer mongoDb.Client.Disconnect(context.Background())

		// Weaviate connection
		log.Println("Connecting to Weaviate...")
		if err := db.InitializeWeaviateFromEnv(); err != nil {
			log.Fatalf("Failed to initialize Weaviate: %v", err)
		}
		defer db.CloseWeaviate()

		// Initialize router
		r := mux.NewRouter()

		// Add logging middleware to the router
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Printf("Request: %s %s", r.Method, r.URL.Path)
				next.ServeHTTP(w, r)
			})
		})

		// Register API routes
		api.RegisterRoutes(r)

		// Create a handler chain with CORS handled at the HTTP server level
		corsHandler := globalCorsHandler(r)

		// Server configuration and startup
		port := os.Getenv("SERVER_PORT")
		if port == "" {
			port = config.InitConfig().ServerPort // fallback to config
		}
		if port == "" {
			port = "8080" // final fallback
		}

		log.Printf("🚀 ESP Organizer API Server starting on port %s", port)

		// Check if port is in use
		if isPortInUse(port) {
			log.Printf("⚠️  Port %s is already in use. Attempting to find an alternative...", port)
			alternatives := []string{"8081", "8082", "8083", "3000", "3001"}
			foundAlternative := false
			for _, altPort := range alternatives {
				if !isPortInUse(altPort) {
					port = altPort
					log.Printf("✅ Using alternative port %s", port)
					foundAlternative = true
					break
				}
			}
			if !foundAlternative {
				log.Fatalf("❌ All alternative ports are in use. Please free port 8080 or kill conflicting processes")
			}
		}

		log.Printf("📊 Available endpoints:")
		log.Printf("   GET  / - Welcome message")
		log.Printf("   GET  /api/skills/semantic-query?q=memory - Test semantic search")
		log.Printf("   POST /api/immunology/upload-chapter - Upload immunology PDFs")
		log.Printf("   GET  /api/immunology/chapters - View processed chapters")
		log.Printf("🌐 Access at: http://localhost:%s", port)

		// Start the server with our CORS handler wrapping everything
		if err := http.ListenAndServe(":"+port, corsHandler); err != nil {
			log.Fatalf("❌ Could not start server: %v", err)
		}
	}
}

// isPortInUse checks if a port is already in use
func isPortInUse(port string) bool {
	conn, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return true
	}
	conn.Close()
	return false
}
