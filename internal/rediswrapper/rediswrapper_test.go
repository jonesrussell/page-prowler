package rediswrapper

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

type MockLogger struct{}

func (m *MockLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (m *MockLogger) Info(msg string, keysAndValues ...interface{})  {}
func (m *MockLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (m *MockLogger) Error(msg string, keysAndValues ...interface{}) {}
func (m *MockLogger) Fatal(msg string, keysAndValues ...interface{}) {}

func TestNewRedisWrapper(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()

	// Create a mock logger
	log := &MockLogger{}

	// Expect Ping call and return nil (indicating a successful connection)
	mock.ExpectPing().SetVal("PONG")

	// Call NewRedisWrapper function with the mock logger
	rw, err := NewRedisWrapper(ctx, db, log)

	// Assert there was no error and the RedisWrapper was correctly created
	assert.NoError(t, err)
	assert.NotNil(t, rw)
	assert.Equal(t, db, rw.Client)
	assert.Equal(t, log, rw.Log)
}

func TestSAdd(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()

	// Create a mock logger
	log := &MockLogger{}

	// Create a new RedisWrapper instance with the mock logger
	rw := &RedisWrapper{
		Client: db,
		Log:    log,
	}

	key := "testKey"
	values := []interface{}{"value1", "value2"}

	// Expect SAdd call and return 2 (the number of values added)
	mock.ExpectSAdd(key, values...).SetVal(2)

	// Call SAdd method
	added, err := rw.SAdd(ctx, key, values...)

	// Assert there was no error and the return value is as expected
	assert.NoError(t, err)
	assert.Equal(t, int64(2), added)
}

func TestDel(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()

	// Create a mock logger
	log := &MockLogger{}

	// Create a new RedisWrapper instance with the mock logger
	rw := &RedisWrapper{
		Client: db,
		Log:    log,
	}

	keys := []string{"key1", "key2"}

	// Expect Del call and return 2 (the number of keys deleted)
	mock.ExpectDel(keys...).SetVal(2)

	// Call Del method
	deleted, err := rw.Del(ctx, keys...)

	// Assert there was no error and the return value is as expected
	assert.NoError(t, err)
	assert.Equal(t, int64(2), deleted)
}

func TestProcess(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()

	// Create a mock logger
	log := &MockLogger{}

	// Create a new RedisWrapper instance with the mock logger
	rw := &RedisWrapper{
		Client: db,
		Log:    log,
	}

	// Create a slice of XMessage
	messages := []redis.XMessage{
		{
			ID: "1",
			Values: map[string]interface{}{
				"event": "received",
				"href":  "http://example.com",
				"group": "testGroup",
			},
		},
	}

	stream := "testStream"
	group := "testGroup"

	// Expect XAck call and return 1 (the number of messages acknowledged)
	mock.ExpectXAck(stream, group, "1").SetVal(1)

	// Call Process method
	posts := rw.Process(ctx, messages, stream, group)

	// Assert the return value is as expected
	assert.Equal(t, 1, len(posts))
	assert.Equal(t, "http://example.com", posts[0].Href)
	assert.Equal(t, "testGroup", posts[0].Group)
}

func TestNewRedisWrapper_Error(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()

	// Create a mock logger
	log := &MockLogger{}

	// Expect Ping call and return an error
	mock.ExpectPing().SetErr(errors.New("Redis server not available"))

	// Call NewRedisWrapper function with the mock logger
	rw, err := NewRedisWrapper(ctx, db, log)

	// Assert there was an error and the RedisWrapper was not created
	assert.Error(t, err)
	assert.Nil(t, rw)
}

func TestSAdd_Error(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()

	// Create a mock logger
	log := &MockLogger{}

	// Create a new RedisWrapper instance with the mock logger
	rw := &RedisWrapper{
		Client: db,
		Log:    log,
	}

	key := "testKey"
	values := []interface{}{"value1", "value2"}

	// Expect SAdd call and return an error
	mock.ExpectSAdd(key, values...).SetErr(errors.New("Error adding to set"))

	// Call SAdd method
	added, err := rw.SAdd(ctx, key, values...)

	// Assert there was an error and no values were added
	assert.Error(t, err)
	assert.Equal(t, int64(0), added)
}

func TestDel_Error(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()

	// Create a mock logger
	log := &MockLogger{}

	// Create a new RedisWrapper instance with the mock logger
	rw := &RedisWrapper{
		Client: db,
		Log:    log,
	}

	keys := []string{"key1", "key2"}

	// Expect Del call and return an error
	mock.ExpectDel(keys...).SetErr(errors.New("Error deleting keys"))

	// Call Del method
	deleted, err := rw.Del(ctx, keys...)

	// Assert there was an error and no keys were deleted
	assert.Error(t, err)
	assert.Equal(t, int64(0), deleted)
}
