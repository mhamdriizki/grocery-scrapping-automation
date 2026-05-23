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

	// Init Redis Asynq Client
	asynqClient, err := config.NewRedisAsynqClient()
	if err != nil {
		log.Fatalf("Failed to initialize Redis Asynq client: %v", err)
	}
	defer asynqClient.Close()

	log.Println("Successfully connected to Redis Asynq")

	// TODO: Setup delivery layer (HTTP routes, CLI handlers) here
}
