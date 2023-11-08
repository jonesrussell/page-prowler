// Package rediswrapper provides a wrapper for the Redis client to facilitate Redis operations.
package rediswrapper

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/jonesrussell/crawler/internal/logger"
)

// MsgPost represents a message with a URL and group, typically used for message queues.
type MsgPost struct {
	Href  string `json:"href"`
	Group string `json:"group"`
}

// RedisWrapper encapsulates a Redis client and associated methods.
type RedisWrapper struct {
	Client      redis.Cmdable
	crawlsiteID string
	mu          sync.Mutex
}

type Cmdable interface {
	redis.Cmdable
	Ping(ctx context.Context) *redis.StatusCmd
}

func NewRedisWrapper(ctx context.Context, client Cmdable) (*RedisWrapper, error) {
	// Test the connection to Redis.
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	// Return the new RedisWrapper instance.
	return &RedisWrapper{
		Client: client,
	}, nil
}

// SetCrawlsiteID sets the identifier for the current crawl site in the RedisWrapper.
func (rw *RedisWrapper) SetCrawlsiteID(id string) {
	rw.mu.Lock()
	rw.crawlsiteID = id
	rw.mu.Unlock()
}

// SAdd adds one or more values to a set at a given key.
func (rw *RedisWrapper) SAdd(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return rw.Client.SAdd(ctx, key, values...).Result()
}

// SMembers retrieves all the members of the set stored at a given key.
func (rw *RedisWrapper) SMembers(ctx context.Context, key string) ([]string, error) {
	return rw.Client.SMembers(ctx, key).Result()
}

// PublishHref publishes a message under a specific channel in Redis.
func (rw *RedisWrapper) PublishHref(ctx context.Context, channel, message string) error {
	return rw.Client.Publish(ctx, channel, message).Err()
}

// Del deletes one or more keys from Redis and returns the number of keys that were removed.
func (rw *RedisWrapper) Del(ctx context.Context, keys ...string) (int64, error) {
	return rw.Client.Del(ctx, keys...).Result()
}

// Stream creates a new group in a Redis stream with the specified stream key.
func (rw *RedisWrapper) Stream(ctx context.Context, stream string, group string) error {
	return rw.Client.XGroupCreate(ctx, stream, group, "$").Err()
}

// Entries reads entries from a stream associated with a Redis group.
func (rw *RedisWrapper) Entries(ctx context.Context, group string, stream string) ([]redis.XStream, error) {
	return rw.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: "consumer",
		Streams:  []string{stream, ">"},
		Count:    10,
		Block:    0,
		NoAck:    false,
	}).Result()
}

// Messages is a helper function that extracts messages from a slice of XStream.
func Messages(entries []redis.XStream) []redis.XMessage {
	var messages []redis.XMessage
	for _, entry := range entries {
		messages = append(messages, entry.Messages...)
	}
	return messages
}

// Process takes a slice of XMessage and processes each message accordingly.
func (rw *RedisWrapper) Process(ctx context.Context, messages []redis.XMessage, stream string, group string, log logger.Logger) []MsgPost {
	var posts []MsgPost
	for _, msg := range messages {
		values := msg.Values
		eventName := fmt.Sprintf("%v", values["event"])
		href := fmt.Sprintf("%v", values["href"])
		group := fmt.Sprintf("%v", values["group"])

		if eventName == "received" {
			post := MsgPost{
				Href:  href,
				Group: group,
			}
			posts = append(posts, post)
			// Acknowledge the message so it isn't processed again
			if err := rw.ackEntry(ctx, stream, group, msg.ID); err != nil {
				log.Error("Failed to acknowledge message", "ID", msg.ID, "stream", stream, "error", err)
			}
		}
	}
	return posts
}

// ackEntry acknowledges a processed message from a stream, so it doesn't get sent again.
func (rw *RedisWrapper) ackEntry(ctx context.Context, stream string, group string, id string) error {
	return rw.Client.XAck(ctx, stream, group, id).Err()
}
