package database

import (
	"context"
	"fmt"
	"log"

	"github.com/akbarwjyy/go-commerce-api/pkg/config"
	"github.com/redis/go-redis/v9"
)

// NewRedisClient membuat koneksi baru ke Redis
func NewRedisClient(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Successfully connected to Redis")
	return client, nil
}
