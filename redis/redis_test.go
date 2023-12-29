package redis

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	db, _ := redismock.NewClientMock()
	client := &Client{Client: db}
	assert.NotNil(t, client)
}

func TestPing(t *testing.T) {
	db, mock := redismock.NewClientMock()
	mock.ExpectPing().SetVal("PONG")

	client := &Client{Client: db}
	cmd := client.Ping(context.Background())
	assert.Equal(t, "PONG", cmd.Val())
}

func TestSAdd(t *testing.T) {
	db, mock := redismock.NewClientMock()
	mock.ExpectSAdd("testKey", "testValue").SetVal(1)

	client := &Client{Client: db}
	count, err := client.SAdd(context.Background(), "testKey", "testValue")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestSMembers(t *testing.T) {
	db, mock := redismock.NewClientMock()
	mock.ExpectSMembers("testKey").SetVal([]string{"testValue"})

	client := &Client{Client: db}
	members, err := client.SMembers(context.Background(), "testKey")
	assert.NoError(t, err)
	assert.Contains(t, members, "testValue")
}

func TestDel(t *testing.T) {
	db, mock := redismock.NewClientMock()
	mock.ExpectDel("testKey").SetVal(1)

	client := &Client{Client: db}
	count, err := client.Del(context.Background(), "testKey")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}
