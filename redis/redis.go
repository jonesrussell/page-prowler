package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
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

// RedisClient is a wrapper around the go-redis Client that implements the ClientInterface.
type RedisClient struct {
	*goredis.Client
}

func (c *RedisClient) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}

func (c *RedisClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return c.Client.SAdd(ctx, key, members...).Err()
}

func (c *RedisClient) Del(ctx context.Context, keys ...string) error {
	cmd := c.Client.Del(ctx, keys...)
	return cmd.Err()
}

func (c *RedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.Client.SMembers(ctx, key).Result()
}

func (c *RedisClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return c.Client.SIsMember(ctx, key, member).Result()
}

func (c *RedisClient) Options() *Options {
	opts := c.Client.Options()
	return &Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	}
}

// NewClient creates a new Redis client.
func NewClient(ctx context.Context, address string, password string, port string) (ClientInterface, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     address + ":" + port,
		Password: password, // Use the Redis password
		DB:       0,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}
	return &RedisClient{client}, nil
}
