package main

import (
	"context"
	"log"

	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/usecase/scraper"
)

func main() {
	s := scraper.NewTipTopScraper()
	products, err := s.ScrapeKepDateperDapur(context.Background())
	if err != nil {
		log.Fatalf("Scraper error: %v", err)
	}
	log.Printf("Successfully got %d products", len(products))
}
