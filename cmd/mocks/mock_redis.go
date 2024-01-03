package mocks

import (
	"context"

	"github.com/jonesrussell/page-prowler/redis"
)

type MockRedisClient struct {
	Data map[string][]string
}

var _ redis.ClientInterface = (*MockRedisClient)(nil)

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		Data: make(map[string][]string),
	}
}

func (m *MockRedisClient) Options() *redis.Options {
	// Return a mock Options object
	return &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
}

func (m *MockRedisClient) Ping(_ context.Context) error {
	// Always return nil for simplicity. You could also simulate network errors by returning non-nil errors.
	return nil
}

func (m *MockRedisClient) SAdd(_ context.Context, key string, members ...interface{}) error {
	for _, member := range members {
		m.Data[key] = append(m.Data[key], member.(string))
	}
	// Always return nil for simplicity. You could also simulate network errors by returning non-nil errors.
	return nil
}

func (m *MockRedisClient) Del(_ context.Context, keys ...string) error {
	for _, key := range keys {
		delete(m.Data, key)
	}
	// Always return nil for simplicity. You could also simulate network errors by returning non-nil errors.
	return nil
}

func (m *MockRedisClient) SMembers(_ context.Context, key string) ([]string, error) {
	// Return the members of the set and nil error for simplicity. You could also simulate network errors by returning non-nil errors.
	return m.Data[key], nil
}

func (m *MockRedisClient) SIsMember(_ context.Context, key string, member interface{}) (bool, error) {
	// Check if the member exists in the set
	for _, m := range m.Data[key] {
		if m == member.(string) {
			return true, nil
		}
	}
	return false, nil
}
