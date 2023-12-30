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

func (m *MockRedisClient) SAdd(_ context.Context, key string, members ...interface{}) *redis.IntCmd {
	for _, member := range members {
		m.Data[key] = append(m.Data[key], member.(string))
	}
	return redis.NewIntResult(int64(len(members)), nil)
}

func (m *MockRedisClient) Del(_ context.Context, keys ...string) *redis.IntCmd {
	count := 0
	for _, key := range keys {
		if _, ok := m.Data[key]; ok {
			delete(m.Data, key)
			count++
		}
	}
	return redis.NewIntResult(int64(count), nil)
}

func (m *MockRedisClient) SMembers(_ context.Context, key string) *redis.StringSliceCmd {
	return redis.NewStringSliceResult(m.Data[key], nil)
}

func (m *MockRedisClient) Publish(_ context.Context, _ string, _ interface{}) *redis.IntCmd {
	return redis.NewIntResult(1, nil)
}
