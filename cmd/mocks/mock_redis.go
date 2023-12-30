package mocks

import (
	"context"
	goredis "github.com/go-redis/redis/v8"
)

type MockRedisClient struct {
	Data map[string][]string
}

func (m *MockRedisClient) Ping(_ context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("PONG", nil)
}

func (m *MockRedisClient) SAdd(_ context.Context, key string, members ...interface{}) *goredis.IntCmd {
	for _, member := range members {
		m.Data[key] = append(m.Data[key], member.(string))
	}
	return goredis.NewIntResult(int64(len(members)), nil)
}

func (m *MockRedisClient) Del(_ context.Context, keys ...string) *goredis.IntCmd {
	count := 0
	for _, key := range keys {
		if _, ok := m.Data[key]; ok {
			delete(m.Data, key)
			count++
		}
	}
	return goredis.NewIntResult(int64(count), nil)
}

func (m *MockRedisClient) SMembers(_ context.Context, key string) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(m.Data[key], nil)
}

func (m *MockRedisClient) Publish(_ context.Context, _ string, _ interface{}) *goredis.IntCmd {
	return goredis.NewIntResult(1, nil)
}
