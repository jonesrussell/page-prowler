package prowlredis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer func(db *redis.Client) {
		err := db.Close()
		if err != nil {
			t.Error(err)
		}
	}(db)

	mock.ExpectPing().SetVal("PONG")

	client := &ClientRedis{db}
	err := client.Ping(context.Background())

	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestSAdd(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer func(db *redis.Client) {
		err := db.Close()
		if err != nil {
			t.Error(err)
		}
	}(db)

	key := "testKey"
	member := "testMember"

	mock.ExpectSAdd(key, member).SetVal(int64(1))

	client := &ClientRedis{db}
	err := client.SAdd(context.Background(), key, member)

	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestDel(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer func(db *redis.Client) {
		err := db.Close()
		if err != nil {
			t.Error(err)
		}
	}(db)

	keys := []string{"key1", "key2"}

	mock.ExpectDel(keys[0], keys[1]).SetVal(int64(2))

	client := &ClientRedis{db}
	err := client.Del(context.Background(), keys...)

	assert.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestSMembers(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer func(db *redis.Client) {
		err := db.Close()
		if err != nil {
			t.Error(err)
		}
	}(db)

	key := "testKey"
	members := []string{"member1", "member2"}

	mock.ExpectSMembers(key).SetVal(members)

	client := &ClientRedis{db}
	result, err := client.SMembers(context.Background(), key)

	assert.Nil(t, err)
	assert.Equal(t, members, result)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestSIsMember(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer func(db *redis.Client) {
		err := db.Close()
		if err != nil {
			t.Error(err)
		}
	}(db)

	key := "testKey"
	member := "testMember"

	mock.ExpectSIsMember(key, member).SetVal(true)

	client := &ClientRedis{db}
	result, err := client.SIsMember(context.Background(), key, member)

	assert.Nil(t, err)
	assert.True(t, result)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestOptions(t *testing.T) {
	db, _ := redismock.NewClientMock()
	defer func(db *redis.Client) {
		err := db.Close()
		if err != nil {
			t.Error(err)
		}
	}(db)

	client := &ClientRedis{db}
	options := client.Options()

	assert.NotNil(t, options)
	assert.Equal(t, "localhost:6379", options.Addr)
	assert.Equal(t, "", options.Password)
	assert.Equal(t, 0, options.DB)
}
