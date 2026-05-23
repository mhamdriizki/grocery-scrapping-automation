package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresDB initializes a new PostgreSQL connection pool using pgx.
func NewPostgresDB(ctx context.Context) (*pgxpool.Pool, error) {
	// Defaults to port 5433 as requested if not set
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := getEnvOrDefault("DB_HOST", "localhost")
		port := getEnvOrDefault("DB_PORT", "5433")
		user := getEnvOrDefault("DB_USER", "postgres")
		pass := getEnvOrDefault("DB_PASSWORD", "postgres")
		name := getEnvOrDefault("DB_NAME", "grocery_engine")
		
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres connection pool: %w", err)
	}

	// Ping the database to ensure connection is valid
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres database: %w", err)
	}

	return pool, nil
}

// getEnvOrDefault returns the value of an environment variable or a default value if not set.
func getEnvOrDefault(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
