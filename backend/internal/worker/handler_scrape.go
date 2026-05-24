package worker

import (
	"context"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/repository"
	scraperUsecase "github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/usecase/scraper"
)

// ScrapeHandler handles the scraping background tasks.
type ScrapeHandler struct {
	productRepo repository.ProductRepository
}

// NewScrapeHandler creates a new instance of ScrapeHandler.
func NewScrapeHandler(productRepo repository.ProductRepository) *ScrapeHandler {
	return &ScrapeHandler{productRepo: productRepo}
}

// ProcessTaskScrapeGrocery is the handler function for the TaskScrapeGrocery task.
// It parses the payload, executes the scraping logic via the TipTopScraper usecase,
// and saves the results to the database using the repository.
func (h *ScrapeHandler) ProcessTaskScrapeGrocery(ctx context.Context, t *asynq.Task) error {
	payload, err := ParseScrapeGroceryPayload(t)
	if err != nil {
		// Returning a non-retryable error to avoid infinite retry loops on bad payloads
		return fmt.Errorf("%w: %w", asynq.SkipRetry, err)
	}

	log.Printf("[worker] Received scrape task | target_url=%s", payload.TargetURL)

	// Initialize the Tip Top scraper usecase
	scraper := scraperUsecase.NewTipTopScraper()

	// Execute scraping
	products, err := scraper.ScrapeKeperluanDapur(ctx)
	if err != nil {
		return fmt.Errorf("%w: tiptop scraper failed: %v", asynq.SkipRetry, err)
	}

	log.Printf("[worker] Saving %d products to database...", len(products))
	if err := h.productRepo.SaveBatch(ctx, products); err != nil {
		// If DB fails, we might want to retry, so we don't wrap with asynq.SkipRetry here
		return fmt.Errorf("failed to save products to db: %w", err)
	}

	log.Printf("[worker] Scrape task finished | products_saved=%d | target_url=%s", len(products), payload.TargetURL)
	return nil
}
