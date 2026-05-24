package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/config"
	deliveryHttp "github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/delivery/http"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/repository"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/usecase/product"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/worker"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/pkg/database"
)

func main() {
	// Load env vars
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Init PostgreSQL (GORM)
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Init Redis Client & Task Distributor
	redisOpt := config.GetAsynqRedisOpt()
	distributor := worker.NewRedisTaskDistributor(redisOpt)
	defer distributor.Close()

	// Init Repositories
	productRepo := repository.NewProductRepository(db)

	// Init Usecases
	productUsecase := product.NewProductUsecase(productRepo)

	// Init Handlers
	productHandler := deliveryHttp.NewProductHandler(productUsecase)
	scrapeHandler := deliveryHttp.NewScrapeHandler(distributor)

	// Setup HTTP Router
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("/api/v1/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			productHandler.GetProducts(w, r)
		} else {
			deliveryHttp.JSON(w, http.StatusMethodNotAllowed, false, "Method not allowed", nil)
		}
	})

	mux.HandleFunc("/api/v1/scrape", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			scrapeHandler.TriggerScrape(w, r)
		} else {
			deliveryHttp.JSON(w, http.StatusMethodNotAllowed, false, "Method not allowed", nil)
		}
	})

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	// Run Server in a goroutine
	go func() {
		log.Println("Starting API server on port 8081")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("API Server exited properly")
}
