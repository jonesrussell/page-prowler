package mocks

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type MockRedisClient struct {
	data map[string][]string
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string][]string),
	}
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
	// Implement the method...
	return nil
}

func (m *MockRedisClient) PublishHref(ctx context.Context, channel, message string) error {
	// Implement the method...
	return nil
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

func (m *MockRedisClient) Stream(ctx context.Context, stream string, group string) error {
	// Implement the method...
	return nil
}

func (m *MockRedisClient) Entries(ctx context.Context, group string, stream string) ([]redis.XStream, error) {
	// Implement the method...
	return []redis.XStream{}, nil
}
