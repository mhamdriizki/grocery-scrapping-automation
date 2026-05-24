package worker

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
)

// TaskDistributor defines the interface for sending tasks to the background queue.
// This abstraction allows swapping implementations (e.g., for testing with a mock).
type TaskDistributor interface {
	DistributeScrapeGroceryTask(ctx context.Context, payload ScrapeGroceryPayload, opts ...asynq.Option) error
}

// RedisTaskDistributor is the production implementation of TaskDistributor,
// backed by a Redis-based asynq client.
type RedisTaskDistributor struct {
	client *asynq.Client
}

// NewRedisTaskDistributor creates a new RedisTaskDistributor.
func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) *RedisTaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{client: client}
}

// Close shuts down the underlying asynq client connection.
func (d *RedisTaskDistributor) Close() error {
	return d.client.Close()
}

// DistributeScrapeGroceryTask enqueues a scraping job into the Redis task queue.
func (d *RedisTaskDistributor) DistributeScrapeGroceryTask(ctx context.Context, payload ScrapeGroceryPayload, opts ...asynq.Option) error {
	task, err := NewScrapeGroceryTask(payload.TargetURL)
	if err != nil {
		return fmt.Errorf("failed to create scrape grocery task: %w", err)
	}

	info, err := d.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return fmt.Errorf("failed to enqueue scrape grocery task: %w", err)
	}

	// Log task info for observability
	_ = info // info.ID, info.Queue can be used for logging in the future

	return nil
}
