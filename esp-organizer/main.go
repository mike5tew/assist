package main

import (
	"esp-organizer/internal/config"
	"esp-organizer/internal/server"

	"log"
	"os"
)

func main() {
	// Initialize configuration

	_ = config.InitConfig()

	// Set up and start the server
	srv := server.NewServer()

	// Log server startup
	log.Printf("ESP Organizer API starting on port %s", os.Getenv("SERVER_PORT"))

	// Start the server (this will block until the server exits)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
