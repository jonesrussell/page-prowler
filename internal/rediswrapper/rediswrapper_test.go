package rediswrapper

import (
	"context"
	"testing"

	redismock "github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestSMembers(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()

	rw := &RedisWrapper{Client: db}

	// Set expectations
	key := "myset"
	members := []string{"member1", "member2"}
	mock.ExpectSMembers(key).SetVal(members)

	// Call the function that uses SMembers
	gotMembers, err := rw.SMembers(ctx, key)

	assert.NoError(t, err)
	assert.Equal(t, members, gotMembers)

	// Assert that all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSAdd(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()

	rw := &RedisWrapper{Client: db}

	// Set expectations for SAdd
	mock.ExpectSAdd("mykey", "myvalue").SetVal(1)

	// Call the SAdd function
	added, err := rw.SAdd(ctx, "mykey", "myvalue")

	assert.NoError(t, err)
	assert.Equal(t, int64(1), added)

	// Assert that all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDel(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()

	rw := &RedisWrapper{Client: db}

	// Set expectations for Del
	mock.ExpectDel("mykey").SetVal(1) // Assume it deletes one item

	// Call the Del function
	deleted, err := rw.Del(ctx, "mykey")

	assert.NoError(t, err)
	assert.Equal(t, int64(1), deleted)

	// Assert that all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPublishHref(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()

	rw := &RedisWrapper{Client: db}

	// Set expectations for Publish
	channel := "mychannel"
	message := "mymessage"
	mock.ExpectPublish(channel, message).SetVal(1) // Assume it publishes one message

	// Call the PublishHref function
	err := rw.PublishHref(ctx, channel, message)

	assert.NoError(t, err)

	// Assert that all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
