package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ocmodi21/image-processing-service/internal/api"
	"github.com/ocmodi21/image-processing-service/internal/queue"
	"github.com/ocmodi21/image-processing-service/internal/service"
	"github.com/ocmodi21/image-processing-service/internal/storage"
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Parse command line flags
	var (
		addr       = flag.String("addr", ":8080", "HTTP server address")
		storePath  = flag.String("store-path", "./store_master.csv", "Path to the store master CSV file")
		numWorkers = flag.Int("workers", 4, "Number of worker goroutines")
	)
	flag.Parse()

	// Create storage
	jobStorage := storage.NewJobStorage()
	storeStorage := storage.NewStoreStorage()

	// Load store data
	if err := storeStorage.LoadFromCSV(*storePath); err != nil {
		log.Fatalf("Failed to load store data: %v", err)
	}

	// Create image processor
	imageProcessor := service.NewImageProcessor()

	// Create job service (circular dependency, will be resolved below)
	var jobService *service.JobService

	// Create job queue with processor function
	jobQueue := queue.NewJobQueue(*numWorkers, func(jobID string) {
		jobService.ProcessJob(jobID)
	})

	// Resolve circular dependency
	jobService = service.NewJobService(jobStorage, storeStorage, jobQueue, imageProcessor)

	// Create API handler and server
	handler := api.NewHandler(jobService)
	server := api.NewServer(*addr, handler)

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
