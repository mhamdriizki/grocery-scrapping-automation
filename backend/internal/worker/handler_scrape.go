package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

// ProcessTaskScrapeGrocery is the handler function for the TaskScrapeGrocery task.
// It parses the payload and executes the scraping logic.
// In this initial version, it simulates work as a placeholder for the chromedp integration.
func ProcessTaskScrapeGrocery(ctx context.Context, t *asynq.Task) error {
	payload, err := ParseScrapeGroceryPayload(t)
	if err != nil {
		// Returning a non-retryable error to avoid infinite retry loops on bad payloads
		return fmt.Errorf("%w: %w", asynq.SkipRetry, err)
	}

	log.Printf("[worker] Processing scrape task | target_url=%s", payload.TargetURL)

	// --- Placeholder for chromedp scraping logic ---
	// This simulates the time a real scraping job would take.
	select {
	case <-ctx.Done():
		// Handle context cancellation (e.g., worker shutting down)
		return ctx.Err()
	case <-time.After(2 * time.Second):
		// Simulate scraping work
	}
	// -----------------------------------------------

	log.Printf("[worker] Scrape task completed successfully | target_url=%s", payload.TargetURL)
	return nil
}
