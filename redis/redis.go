package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisInterface interface {
	Ping(ctx context.Context) *redis.StatusCmd
	SAdd(ctx context.Context, key string, values ...interface{}) (int64, error)
	SMembers(ctx context.Context, key string) ([]string, error)
	PublishHref(ctx context.Context, channel, message string) error
	Del(ctx context.Context, keys ...string) (int64, error)
	Stream(ctx context.Context, stream string, group string) error
	Entries(ctx context.Context, group string, stream string) ([]redis.XStream, error)
}

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(address string, password string) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // Use the Redis password
		DB:       0,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}
	return &RedisClient{
		Client: client,
	}, nil
}

// Implement the RedisInterface methods
func (r *RedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	return r.Client.Ping(ctx)
}

func (r *RedisClient) SAdd(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return r.Client.SAdd(ctx, key, values...).Result()
}

func (r *RedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.Client.SMembers(ctx, key).Result()
}

func (r *RedisClient) PublishHref(ctx context.Context, channel, message string) error {
	return r.Client.Publish(ctx, channel, message).Err()
}

func (r *RedisClient) Del(ctx context.Context, keys ...string) (int64, error) {
	return r.Client.Del(ctx, keys...).Result()
}

func (r *RedisClient) Stream(ctx context.Context, stream string, group string) error {
	// Implement this method based on your requirements
	return nil
}

func (r *RedisClient) Entries(ctx context.Context, group string, stream string) ([]redis.XStream, error) {
	// Implement this method based on your requirements
	return nil, nil
}
