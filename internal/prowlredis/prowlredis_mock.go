package prowlredis

import "context"

type MockClient struct {
	pingErr error
	data    map[string][]string
}

func NewMockClient() ClientInterface {
	return &MockClient{
		data: make(map[string][]string),
	}
}

func (m *MockClient) Ping(ctx context.Context) error {
	return m.pingErr
}

func (m *MockClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	for _, member := range members {
		m.data[key] = append(m.data[key], member.(string))
	}
	return nil
}

func (m *MockClient) Del(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		delete(m.data, key)
	}
	return nil
}

func (m *MockClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return m.data[key], nil
}

func (m *MockClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
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

func (m *MockClient) Options() *Options {
	// Return some default options. You might want to make this configurable.
	return &Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
}
