package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ocmodi21/image-processing-service/config"
	"github.com/ocmodi21/image-processing-service/internal/api"
	"github.com/ocmodi21/image-processing-service/internal/database"
	"github.com/ocmodi21/image-processing-service/internal/queue"
	"github.com/ocmodi21/image-processing-service/internal/service"
	"github.com/ocmodi21/image-processing-service/internal/storage"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("./config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.Database.Provider, cfg.Database.User, cfg.Database.Password, cfg.Database.Dbname, cfg.Database.Host, cfg.Database.SSLmode)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	// Create storage
	jobStorage := storage.NewJobStorage()
	storeStorage := storage.NewStoreStorage(db)

	// Create image processor
	imageProcessor := service.NewImageProcessor()

	// Create job service
	var jobService *service.JobService

	// Create job queue with processor function
	jobQueue := queue.NewJobQueue(cfg.Processing.NumWorkers, func(jobID string) {
		jobService.ProcessJob(jobID)
	})

	jobService = service.NewJobService(jobStorage, storeStorage, jobQueue, imageProcessor)

	// Create API handler and server
	handler := api.NewHandler(jobService)
	server := api.NewServer(cfg.Server.Port, handler)

	// Start job queue
	jobQueue.Start()
	defer jobQueue.Stop()

	// Start the server in a goroutine
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown the server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Println("Shutting down server...")
	if err := server.Stop(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server stopped")
}
