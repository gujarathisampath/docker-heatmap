package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"docker-heatmap/internal/config"
	"docker-heatmap/internal/database"
	"docker-heatmap/internal/router"
	"docker-heatmap/internal/worker"
)

func main() {
	// Load configuration
	config.Load()
	log.Println("Configuration loaded")

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	log.Println("Database connected")

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed")

	// Start background worker
	syncWorker := worker.NewSyncWorker()
	syncWorker.Start()
	defer syncWorker.Stop()

	// Setup router
	app := router.SetupRouter()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	}()

	// Start server
	port := config.AppConfig.Port
	log.Printf("Starting server on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
