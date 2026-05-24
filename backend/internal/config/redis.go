package config

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

// NewRedisClient initializes a standard Redis connection using redis-go.
func NewRedisClient(ctx context.Context) (*redis.Client, error) {
	host := getEnvOrDefault("REDIS_HOST", "localhost")
	port := getEnvOrDefault("REDIS_PORT", "6379")
	pass := getEnvOrDefault("REDIS_PASSWORD", "")

	redisAddr := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: pass,
		DB:       0, // Use default DB
	})

	// Add timeout for ping to avoid hanging if redis is not running
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Ping to ensure connection
	if err := client.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return client, nil
}

// GetAsynqRedisOpt returns the Redis configuration option required by the asynq library.
// It is separate from NewRedisClient because asynq manages its own connection pool internally.
func GetAsynqRedisOpt() asynq.RedisClientOpt {
	host := getEnvOrDefault("REDIS_HOST", "localhost")
	port := getEnvOrDefault("REDIS_PORT", "6379")
	pass := getEnvOrDefault("REDIS_PASSWORD", "")

	return asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: pass,
	}
}
