package config

import (
	"fmt"

	"github.com/hibiken/asynq"
)

// NewRedisAsynqClient initializes a Redis connection used by asynq.
func NewRedisAsynqClient() (*asynq.Client, error) {
	host := getEnvOrDefault("REDIS_HOST", "localhost")
	port := getEnvOrDefault("REDIS_PORT", "6379")
	pass := getEnvOrDefault("REDIS_PASSWORD", "")

	redisAddr := fmt.Sprintf("%s:%s", host, port)

	// In a real scenario, you can define more options
	opt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: pass,
	}

	client := asynq.NewClient(opt)

	// We ping Redis by fetching connection or using inspector if needed,
	// but asynq creates connections lazily or handles dialing on operation.
	// Asynq Client doesn't have a direct Ping method, so we can test it 
	// by checking inspector or creating a test task. 
	// For simplicity, we just return the client.
	
	return client, nil
}

// GetAsynqRedisOpt returns the Redis configuration option for Asynq workers.
func GetAsynqRedisOpt() asynq.RedisClientOpt {
	host := getEnvOrDefault("REDIS_HOST", "localhost")
	port := getEnvOrDefault("REDIS_PORT", "6379")
	pass := getEnvOrDefault("REDIS_PASSWORD", "")

	return asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: pass,
	}
}
