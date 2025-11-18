package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"esp-organizer/internal/domain/infoin"
)

func main() {
	// Create a context with cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received termination signal, shutting down...")
		cancel()
	}()

	log.Println("Starting revectorization of all content...")

	// Call the function
	if err := infoin.ReVectorizeAllContent(ctx); err != nil {
		log.Fatalf("Error during revectorization: %v", err)
	}

	log.Println("✅ Revectorization completed successfully")
}
