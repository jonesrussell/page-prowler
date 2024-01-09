package mocks

import (
	"context"
	"github.com/go-redis/redismock/v9"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
)

type MockRedisClient struct {
	redismock.ClientMock
}

func NewMockRedisClient() *MockRedisClient {
	_, mock := redismock.NewClientMock()
	return &MockRedisClient{
		ClientMock: mock,
	}
}

func (m *MockRedisClient) Ping(ctx context.Context) error {
	m.ExpectPing().SetVal("PONG")
	return nil
}

func (m *MockRedisClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	m.ExpectSAdd(key, members...).SetVal(int64(len(members)))
	return nil
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) error {
	m.ExpectDel(keys...).SetVal(int64(len(keys)))
	return nil
}

func (m *MockRedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	m.ExpectSMembers(key).SetVal([]string{})
	return []string{}, nil
}

func (m *MockRedisClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	m.ExpectSIsMember(key, member).SetVal(false)
	return false, nil
}

func (m *MockRedisClient) Options() *prowlredis.Options {
	return &prowlredis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
}
