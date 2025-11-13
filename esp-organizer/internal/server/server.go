package server

import (
	"esp-organizer/internal/api"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server encapsulates the HTTP server for the ESP Organizer API
type Server struct {
	router *gin.Engine
	port   string
}

// NewServer creates a new server instance
func NewServer() *Server {
	router := setupRouter()
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	return &Server{
		router: router,
		port:   port,
	}
}

// setupRouter configures the Gin router with all necessary routes
func setupRouter() *gin.Engine {
	// Create a default gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Add a middleware to log the request URL
	router.Use(func(c *gin.Context) {
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	// Basic health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Add a debug endpoint directly in server.go
	router.GET("/debug-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Debug endpoint is working",
		})
	})

	// Register diagnostics routes
	api.RegisterDiagnosticsRoutes(router)

	// Log registered routes after registration
	if os.Getenv("DEBUG_MODE") == "true" {
		log.Println("Routes registered in server.go:")
		for _, route := range router.Routes() {
			log.Printf("  %s %s", route.Method, route.Path)
		}
	}

	return router
}

// Start begins listening for HTTP requests
func (s *Server) Start() error {
	log.Printf("Starting server on port %s", s.port)
	return s.router.Run(fmt.Sprintf(":%s", s.port))
}

// GetRouter returns the Gin router instance
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
