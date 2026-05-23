package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/config"
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

	// Init Redis Client
	redisClient, err := config.NewRedisClient(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}
	defer redisClient.Close()

	log.Println("Successfully connected to Redis")

	// TODO: Setup delivery layer (HTTP routes, CLI handlers) here
}
