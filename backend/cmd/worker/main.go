package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/config"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/worker"
)

func main() {
	// Load env vars from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Get asynq-specific Redis options
	redisOpt := config.GetAsynqRedisOpt()

	// Initialize the task processor (worker server)
	// Concurrency of 5: process up to 5 scraping jobs in parallel
	processor := worker.NewRedisTaskProcessor(redisOpt, 5)

	// Start the processor in a goroutine since server.Run() is blocking
	log.Println("Starting Asynq worker server...")
	go func() {
		if err := processor.Start(); err != nil {
			log.Fatalf("Asynq worker server failed: %v", err)
		}
	}()

	// Graceful shutdown: wait for OS termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Asynq worker server gracefully...")
	processor.Shutdown()
	log.Println("Worker server stopped.")
}
