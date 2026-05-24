package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/config"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/worker"
)

func main() {
	// Load env vars
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	ctx := context.Background()

	// Init Database
	dbPool, err := config.NewPostgresDB(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbPool.Close()

	log.Println("Successfully connected to PostgreSQL")

	// Init Redis Client (for direct queries if needed)
	redisClient, err := config.NewRedisClient(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}
	defer redisClient.Close()

	log.Println("Successfully connected to Redis")

	// Init Task Distributor (API server uses this to enqueue scraping jobs)
	redisOpt := config.GetAsynqRedisOpt()
	distributor := worker.NewRedisTaskDistributor(redisOpt)
	defer distributor.Close()

	// Enqueue a test scraping job to verify the queue works end-to-end
	err = distributor.DistributeScrapeGroceryTask(ctx, worker.ScrapeGroceryPayload{
		TargetURL: "https://www.tokopedia.com/superindo",
	})
	if err != nil {
		log.Printf("Warning: Failed to enqueue test task: %v", err)
	} else {
		log.Println("Successfully enqueued test task: task:scrape_grocery")
	}

	// TODO: Setup HTTP delivery layer (routes, handlers) here
}
