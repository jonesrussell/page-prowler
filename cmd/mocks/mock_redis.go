package mocks

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type MockRedisClient struct{}

func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	// Implement the method...
	return nil
}

func (m *MockRedisClient) SAdd(ctx context.Context, key string, values ...interface{}) (int64, error) {
	// Implement the method...
	return 0, nil
}

func (m *MockRedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	// Implement the method...
	return []string{}, nil
}

func (m *MockRedisClient) PublishHref(ctx context.Context, channel, message string) error {
	// Implement the method...
	return nil
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) (int64, error) {
	// Implement the method...
	return 0, nil
}

func (m *MockRedisClient) Stream(ctx context.Context, stream string, group string) error {
	// Implement the method...
	return nil
}

func (m *MockRedisClient) Entries(ctx context.Context, group string, stream string) ([]redis.XStream, error) {
	// Implement the method...
	return []redis.XStream{}, nil
}
