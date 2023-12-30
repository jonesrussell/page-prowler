package mocks

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type MockRedisClient struct {
	Data map[string][]string
}

func (m *MockRedisClient) Ping(_ context.Context) *redis.StatusCmd {
	return redis.NewStatusResult("PONG", nil)
}

func (m *MockRedisClient) SMembers(_ context.Context, key string) ([]string, error) {
	return m.Data[key], nil
}

func (m *MockRedisClient) Del(_ context.Context, keys ...string) (int64, error) {
	for _, key := range keys {
		delete(m.Data, key)
	}
	return int64(len(keys)), nil
}

func (m *MockRedisClient) SAdd(_ context.Context, key string, values ...interface{}) (int64, error) {
	for _, value := range values {
		m.Data[key] = append(m.Data[key], value.(string))
	}
	return int64(len(values)), nil
}

func (m *MockRedisClient) PublishHref(_ context.Context, _, _ string) error {
	return nil
}
