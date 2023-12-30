package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// ClientInterface represents the interface for interacting with the Redis datastore.
type ClientInterface interface {
	Ping(ctx context.Context) *redis.StatusCmd
	SAdd(ctx context.Context, key string, values ...interface{}) (int64, error)
	SMembers(ctx context.Context, key string) ([]string, error)
	PublishHref(ctx context.Context, channel, message string) error
	Del(ctx context.Context, keys ...string) (int64, error)
}

// Client represents the Redis client.
type Client struct {
	ClientInterface
}

type ClientWrapper struct {
	*redis.Client
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
	return &Client{
		&ClientWrapper{client},
	}, nil
}

// Ping sends a ping request to the Redis server.
func (r *ClientWrapper) Ping(ctx context.Context) *redis.StatusCmd {
	return r.Client.Ping(ctx)
}

// SMembers returns all the members of the set value stored at key.
func (r *ClientWrapper) SMembers(ctx context.Context, key string) ([]string, error) {
	cmd := r.Client.SMembers(ctx, key)
	result, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Del deletes one or more keys.
func (r *ClientWrapper) Del(ctx context.Context, keys ...string) (int64, error) {
	cmd := r.Client.Del(ctx, keys...)
	result, err := cmd.Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// SAdd adds one or more members to a set.
func (r *ClientWrapper) SAdd(ctx context.Context, key string, values ...interface{}) (int64, error) {
	cmd := r.Client.SAdd(ctx, key, values...)
	result, err := cmd.Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// PublishHref publishes a message to a specified channel.
func (r *ClientWrapper) PublishHref(
	ctx context.Context,
	channel string,
	message string,
) error {
	cmd := r.Client.Publish(ctx, channel, message)
	_, err := cmd.Result()
	return err
}
