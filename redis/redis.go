package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// ClientInterface represents the interface for interacting with the Redis datastore.
type ClientInterface interface {
	Ping(ctx context.Context) *redis.StatusCmd
	SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	SMembers(ctx context.Context, key string) *redis.StringSliceCmd
}

// Client represents the Redis client.
type Client struct {
	ClientInterface
}

// NewClient creates a new Redis client.
func NewClient(address string, password string, port string) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address + ":" + port,
		Password: password, // Use the Redis password
		DB:       0,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}
	return &Client{client}, nil
}
