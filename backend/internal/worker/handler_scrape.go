package worker

import (
	"context"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	scraperUsecase "github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/usecase/scraper"
)

// ProcessTaskScrapeGrocery is the handler function for the TaskScrapeGrocery task.
// It parses the payload and executes the scraping logic via the TipTopScraper usecase.
func ProcessTaskScrapeGrocery(ctx context.Context, t *asynq.Task) error {
	payload, err := ParseScrapeGroceryPayload(t)
	if err != nil {
		// Returning a non-retryable error to avoid infinite retry loops on bad payloads
		return fmt.Errorf("%w: %w", asynq.SkipRetry, err)
	}

	log.Printf("[worker] Received scrape task | target_url=%s", payload.TargetURL)

	// Initialize the Tip Top scraper usecase
	scraper := scraperUsecase.NewTipTopScraper()

	// Execute scraping — results are printed to terminal in Phase 1
	products, err := scraper.ScrapeKepDateperDapur(ctx)
	if err != nil {
		return fmt.Errorf("tiptop scraper failed: %w", err)
	}

	log.Printf("[worker] Scrape task finished | products_found=%d | target_url=%s", len(products), payload.TargetURL)
	return nil
}
