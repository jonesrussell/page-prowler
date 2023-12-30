package mocks

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type MockRedisClient struct {
	data map[string][]string
}

func (m *MockRedisClient) SAdd(ctx context.Context, key string, values ...interface{}) (int64, error) {
	for _, value := range values {
		m.data[key] = append(m.data[key], value.(string))
	}
	return int64(len(values)), nil
}

func (m *MockRedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return m.data[key], nil
}

func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	// Implement your test logic here
	return &redis.StatusCmd{}
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) (int64, error) {
	count := 0
	for _, key := range keys {
		if _, ok := m.data[key]; ok {
			delete(m.data, key)
			count++
		}
	}
	return int64(count), nil
}

func (m *MockRedisClient) PublishHref(ctx context.Context, channel string, message string) error {
	// Implement your test logic here
	return nil
}
