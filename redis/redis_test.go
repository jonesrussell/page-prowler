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
	saddCmd := client.SAdd(context.Background(), "testKey", "testValue")
	count, err := saddCmd.Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestSMembers(t *testing.T) {
	db := &mocks.MockRedisClient{Data: make(map[string][]string)}
	client := &Client{ClientInterface: db}
	saddCmd := client.SAdd(context.Background(), "testKey", "testValue")
	_, err := saddCmd.Result()
	assert.NoError(t, err)
	smembersCmd := client.SMembers(context.Background(), "testKey")
	members, err := smembersCmd.Result()
	assert.NoError(t, err)
	assert.Contains(t, members, "testValue")
}

func TestDel(t *testing.T) {
	db := &mocks.MockRedisClient{Data: make(map[string][]string)}
	client := &Client{ClientInterface: db}
	saddCmd := client.SAdd(context.Background(), "testKey", "testValue")
	_, err := saddCmd.Result()
	assert.NoError(t, err)
	delCmd := client.Del(context.Background(), "testKey")
	count, err := delCmd.Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}
