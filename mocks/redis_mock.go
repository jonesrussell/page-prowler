package mocks

import (
	"context"

	"github.com/jonesrussell/page-prowler/internal/prowlredis"
)

type Options struct {
	Addr     string
	Password string
	DB       int
}

func (m *MockClient) Options() *prowlredis.Options {
	// Return some default options. You might want to make this configurable.
	return &prowlredis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
}

type MockClient struct {
	pingErr      error
	data         map[string][]string
	WasDelCalled bool
	DelErr       error
}

func NewMockClient() prowlredis.ClientInterface {
	return &MockClient{
		data:         make(map[string][]string),
		WasDelCalled: false,
	}
}

func (m *MockClient) Ping(_ context.Context) error {
	return m.pingErr
}

func (m *MockClient) SAdd(_ context.Context, key string, members ...interface{}) error {
	for _, member := range members {
		m.data[key] = append(m.data[key], member.(string))
	}
	return nil
}

func (m *MockClient) Del(_ context.Context, keys ...string) error {
	m.WasDelCalled = true
	if m.DelErr != nil {
		return m.DelErr
	}
	for _, key := range keys {
		delete(m.data, key)
	}
	return nil
}

func (m *MockClient) SMembers(_ context.Context, key string) ([]string, error) {
	return m.data[key], nil
}

func (m *MockClient) SIsMember(_ context.Context, key string, member interface{}) (bool, error) {
	members, ok := m.data[key]
	if !ok {
		return false, nil
	}
	memberStr := member.(string)
	for _, m := range members {
		if m == memberStr {
			return true, nil
		}
	}
	return false, nil
}
