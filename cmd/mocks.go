package cmd

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type mockRedisClient struct{}

func (m *mockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	// Implement the method...
	return nil
}

func (m *mockRedisClient) SAdd(ctx context.Context, key string, values ...interface{}) (int64, error) {
	// Implement the method...
	return 0, nil
}

func (m *mockRedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	// Implement the method...
	return []string{}, nil
}

func (m *mockRedisClient) PublishHref(ctx context.Context, channel, message string) error {
	// Implement the method...
	return nil
}

func (m *mockRedisClient) Del(ctx context.Context, keys ...string) (int64, error) {
	// Implement the method...
	return 0, nil
}

func (m *mockRedisClient) Stream(ctx context.Context, stream string, group string) error {
	// Implement the method...
	return nil
}

func (m *mockRedisClient) Entries(ctx context.Context, group string, stream string) ([]redis.XStream, error) {
	// Implement the method...
	return []redis.XStream{}, nil
}
