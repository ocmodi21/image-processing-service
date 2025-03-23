package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	RedisDB *redis.Client
}

// ConnectRedis establishes a connection to Redis and returns a RedisClient instance.
func ConnectRedis(ctx context.Context) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Change as needed
		Password: "",               // No password for local development
		DB:       0,                // Default DB
	})

	// Ping the Redis server to check the connection
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("error connecting to Redis: %w", err)
	}

	log.Println("Successfully connected to Redis")

	return &RedisClient{RedisDB: client}, nil
}
