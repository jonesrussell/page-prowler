package prowlredis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Options represents the options for a Redis client.
type Options struct {
	Addr     string
	Password string
	DB       int
}

// ClientInterface represents the interface for interacting with the Redis datastore.
type ClientInterface interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	SAdd(ctx context.Context, key string, members ...interface{}) error
	Del(ctx context.Context, keys ...string) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)
	Options() *Options
}

// Client represents the Redis client.
type Client struct {
	ClientInterface
}

// ClientRedis is a wrapper around the go-redis Client that implements the ClientInterface.
type ClientRedis struct {
	*redis.Client
}

// Set implements ClientInterface.
// Subtle: this method shadows the method (*Client).Set of ClientRedis.Client.
func (c *ClientRedis) Set(_ context.Context, _ string, _ interface{}, _ time.Duration) error {
	panic("unimplemented")
}

func (c *ClientRedis) Close() error {
	return c.Client.Close()
}

func (c *ClientRedis) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}

func (c *ClientRedis) SAdd(ctx context.Context, key string, members ...interface{}) error {
	// Check if the key is empty
	if key == "" {
		return fmt.Errorf("key is not set")
	}
	return c.Client.SAdd(ctx, key, members...).Err()
}

func (c *ClientRedis) Del(ctx context.Context, keys ...string) error {
	cmd := c.Client.Del(ctx, keys...)
	return cmd.Err()
}

func (c *ClientRedis) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.Client.SMembers(ctx, key).Result()
}

func (c *ClientRedis) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return c.Client.SIsMember(ctx, key, member).Result()
}

func (c *ClientRedis) Options() *Options {
	opts := c.Client.Options()
	return &Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	}
}

// NewClient creates a new Redis client.
func NewClient(ctx context.Context, cfg *Options) (ClientInterface, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}
	return &ClientRedis{client}, nil
}
