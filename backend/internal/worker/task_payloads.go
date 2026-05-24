package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

// TaskScrapeGrocery is the unique name/identifier for the scrape task.
const TaskScrapeGrocery = "task:scrape_grocery"

// ScrapeGroceryPayload holds the data required to execute a scraping job.
type ScrapeGroceryPayload struct {
	TargetURL string `json:"target_url"`
}

// NewScrapeGroceryTask creates a new asynq Task for scraping a grocery page.
func NewScrapeGroceryTask(targetURL string) (*asynq.Task, error) {
	payload, err := json.Marshal(ScrapeGroceryPayload{TargetURL: targetURL})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal scrape grocery payload: %w", err)
	}

	return asynq.NewTask(TaskScrapeGrocery, payload), nil
}

// ParseScrapeGroceryPayload parses the raw task payload into a ScrapeGroceryPayload struct.
func ParseScrapeGroceryPayload(task *asynq.Task) (ScrapeGroceryPayload, error) {
	var payload ScrapeGroceryPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return ScrapeGroceryPayload{}, fmt.Errorf("failed to unmarshal scrape grocery payload: %w", err)
	}

	if payload.TargetURL == "" {
		return ScrapeGroceryPayload{}, fmt.Errorf("target_url is required but was empty")
	}

	return payload, nil
}

// Ensure unused context import is used via handler
var _ = context.Background
