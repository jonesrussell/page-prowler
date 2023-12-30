package redis

import (
	"context"
	"testing"

	"github.com/jonesrussell/page-prowler/cmd/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	db := &mocks.MockRedisClient{Data: make(map[string][]string)}
	client := &Client{ClientInterface: db}
	assert.NotNil(t, client)
}

func TestPing(t *testing.T) {
	db := &mocks.MockRedisClient{Data: make(map[string][]string)}
	client := &Client{ClientInterface: db}
	cmd := client.Ping(context.Background())
	assert.Equal(t, "PONG", cmd.Val())
}

func TestSAdd(t *testing.T) {
	db := &mocks.MockRedisClient{Data: make(map[string][]string)}
	client := &Client{ClientInterface: db}
	count, err := client.SAdd(context.Background(), "testKey", "testValue")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestSMembers(t *testing.T) {
	db := &mocks.MockRedisClient{Data: make(map[string][]string)}
	client := &Client{ClientInterface: db}
	_, _ = client.SAdd(context.Background(), "testKey", "testValue")
	members, err := client.SMembers(context.Background(), "testKey")
	assert.NoError(t, err)
	assert.Contains(t, members, "testValue")
}

func TestDel(t *testing.T) {
	db := &mocks.MockRedisClient{Data: make(map[string][]string)}
	client := &Client{ClientInterface: db}
	_, _ = client.SAdd(context.Background(), "testKey", "testValue")
	count, err := client.Del(context.Background(), "testKey")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestPublishHref(t *testing.T) {
	db := &mocks.MockRedisClient{Data: make(map[string][]string)}
	client := &Client{ClientInterface: db}
	err := client.PublishHref(context.Background(), "testChannel", "testMessage")
	assert.NoError(t, err)
}
