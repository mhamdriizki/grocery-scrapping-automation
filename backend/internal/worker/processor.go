package worker

import (
	"fmt"

	"github.com/hibiken/asynq"
)

// TaskProcessor defines the interface for starting and stopping the worker server.
type TaskProcessor interface {
	Start(scrapeHandler *ScrapeHandler) error
	Shutdown()
}

// RedisTaskProcessor is the production implementation of TaskProcessor,
// backed by a Redis-based asynq server.
type RedisTaskProcessor struct {
	server *asynq.Server
}

// NewRedisTaskProcessor creates a new RedisTaskProcessor with the given Redis options.
// Concurrency defines how many tasks can be processed simultaneously.
func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, concurrency int) *RedisTaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: concurrency,
			// ErrorHandler can be added here for centralized error reporting
		},
	)

	return &RedisTaskProcessor{server: server}
}

// Start registers all task handlers and begins processing jobs from the queue.
// This is a blocking call; it runs until Shutdown() is called.
func (p *RedisTaskProcessor) Start(scrapeHandler *ScrapeHandler) error {
	mux := asynq.NewServeMux()

	// Register task handlers
	mux.HandleFunc(TaskScrapeGrocery, scrapeHandler.ProcessTaskScrapeGrocery)

	if err := p.server.Run(mux); err != nil {
		return fmt.Errorf("failed to run asynq server: %w", err)
	}

	return nil
}

// Shutdown gracefully stops the asynq server, waiting for running tasks to finish.
func (p *RedisTaskProcessor) Shutdown() {
	p.server.Shutdown()
}
