package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/config"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/model"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/repository"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/worker"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/pkg/database"
)

func main() {
	// Load env vars from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Initialize Database Connection
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// AutoMigrate the database schema
	log.Println("Running AutoMigrate for database schemas...")
	if err := db.AutoMigrate(&model.Product{}); err != nil {
		log.Fatalf("Failed to run AutoMigrate: %v", err)
	}

	// Initialize Repositories
	productRepo := repository.NewProductRepository(db)

	// Initialize Handlers
	scrapeHandler := worker.NewScrapeHandler(productRepo)

	// Get asynq-specific Redis options
	redisOpt := config.GetAsynqRedisOpt()

	// Initialize the task processor (worker server)
	// Concurrency of 5: process up to 5 scraping jobs in parallel
	processor := worker.NewRedisTaskProcessor(redisOpt, 5)

	// Start the processor in a goroutine since server.Run() is blocking
	log.Println("Starting Asynq worker server...")
	go func() {
		if err := processor.Start(scrapeHandler); err != nil {
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
