package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

// Datastore represents the interface for interacting with the Redis datastore.
type Datastore interface {
	Ping(ctx context.Context) *redis.StatusCmd
	SAdd(ctx context.Context, key string, values ...interface{}) (int64, error)
	SMembers(ctx context.Context, key string) ([]string, error)
	PublishHref(ctx context.Context, channel, message string) error
	Del(ctx context.Context, keys ...string) (int64, error)
	Stream(ctx context.Context, stream string, group string) error
	Entries(ctx context.Context, group string, stream string) ([]redis.XStream, error)
}

// Client represents the Redis client.
type Client struct {
	Client *redis.Client
}

// NewClient creates a new Redis client.
func NewClient(address string, password string, port string) (*Client, error) {
	log.Printf("Creating new Redis client with address: %s, password: %s, port: %s", address, password, port)
	client := redis.NewClient(&redis.Options{
		Addr:     address + ":" + port,
		Password: password, // Use the Redis password
		DB:       0,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}
	return &Client{
		Client: client,
	}, nil
}

// Ping sends a ping request to the Redis server.
func (r *Client) Ping(ctx context.Context) *redis.StatusCmd {
	return r.Client.Ping(ctx)
}

// SAdd adds one or more members to a set.
func (r *Client) SAdd(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return r.Client.SAdd(ctx, key, values...).Result()
}

// SMembers returns all the members of the set value stored at key.
func (r *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.Client.SMembers(ctx, key).Result()
}

// PublishHref publishes a href to the Redis server.
func (r *Client) PublishHref(ctx context.Context, channel, message string) error {
	return r.Client.Publish(ctx, channel, message).Err()
}

// Del deletes one or more keys.
func (r *Client) Del(ctx context.Context, keys ...string) (int64, error) {
	return r.Client.Del(ctx, keys...).Result()
}

// Stream streams data from the Redis server.
func (r *Client) Stream(ctx context.Context, stream string, group string) error {
	// Implement this method based on your requirements
	return nil
}

// Entries returns the entries from the Redis server.
func (r *Client) Entries(ctx context.Context, group string, stream string) ([]redis.XStream, error) {
	// Implement this method based on your requirements
	return nil, nil
}
