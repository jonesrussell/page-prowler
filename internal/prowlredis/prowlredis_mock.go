package prowlredis

import "context"

type MockClient struct {
	ClientInterface
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m *MockClient) Ping(ctx context.Context) error {
	return nil
}

func (m *MockClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return nil
}

func (m *MockClient) Del(ctx context.Context, keys ...string) error {
	return nil
}

func (m *MockClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return nil, nil
}

func (m *MockClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return false, nil
}

func (m *MockClient) Options() *Options {
	return &Options{}
}
