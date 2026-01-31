package main

import (
	"context"
	"esp-organizer/internal/domain/api" // Add this import
	"esp-organizer/internal/domain/coach"
	"esp-organizer/internal/store/db"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"esp-organizer/internal/logger"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables")
	}

	// Initialize structured logger early so all log outputs go through it
	if err := logger.Init(""); err != nil {
		log.Printf("Warning: failed to initialize structured logger: %v", err)
	}
	// route standard library log through logrus
	log.SetOutput(logger.Writer())

	// Initialize MVP Coach Service
	coachService, err := coach.NewCoachServiceMVP() // Fixed: use mvpHandler package
	if err != nil {
		log.Fatalf("FATAL: Failed to create MVP coach service: %v", err)
	}
	log.Println("MVP Coach Service initialized successfully")

	// Initialize MongoDB (commented out as in your original)
	// mongoDB, err := db.NewFromEnv()
	// if err != nil {
	// 	log.Fatalf("Failed to initialize MongoDB: %v", err)
	// }
	// log.Println("MongoDB initialized successfully")

	// Initialize Weaviate
	if err := db.InitializeWeaviateFromEnv(); err != nil {
		log.Printf("WARNING: Failed to initialize Weaviate: %v", err)
		log.Printf("Some semantic search functionality may be limited")
	} else {
		log.Println("Weaviate initialized successfully")
	}

	// Create and configure router with proper routes
	router := mux.NewRouter()

	// CRITICAL FIX: Register routes from the api package with the coach service
	log.Println("Registering API routes...")
	log.Println("🔄 ESP Organizer Backend Version: 1.1.0")
	api.RegisterRoutes(router, coachService)

	// DEBUG: Print all registered routes
	log.Println("Registered Routes:")
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}
		methods, _ := route.GetMethods()
		log.Printf("ROUTE: %v %s", methods, path)
		return nil
	})

	// Add middleware for CORS if needed
	// router.Use(corsMiddleware) // Already added in api.RegisterRoutes

	// Configure the HTTP server
	addr := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting ESP Organizer API server on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Set up graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

// CORS middleware function
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
