package rediswrapper

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// MsgPost represents a message with a URL and group.
type MsgPost struct {
	Href  string `json:"href"`
	Group string `json:"group"`
}

// RedisWrapper encapsulates a Redis client and associated methods.
type RedisWrapper struct {
	Client      *redis.Client
	Logger      *zap.SugaredLogger
	crawlsiteID string
	mu          sync.Mutex
}

// NewRedisWrapper creates a new RedisWrapper with the given configuration.
func NewRedisWrapper(ctx context.Context, host, port, password string, logger *zap.SugaredLogger) (*RedisWrapper, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisWrapper{
		Client: rdb,
		Logger: logger,
	}, nil
}

// SetCrawlsiteID sets the crawlsite ID for this RedisWrapper instance.
func (rw *RedisWrapper) SetCrawlsiteID(id string) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	rw.crawlsiteID = id
}

func (rw *RedisWrapper) SAdd(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return rw.Client.SAdd(ctx, key, values...).Result()
}

func (rw *RedisWrapper) SMembers(ctx context.Context, key string) ([]string, error) {
	return rw.Client.SMembers(ctx, key).Result()
}

func (rw *RedisWrapper) PublishHref(ctx context.Context, channel, message string) error {
	return rw.Client.Publish(ctx, channel, message).Err()
}

func (rw *RedisWrapper) Del(ctx context.Context, keys ...string) (int64, error) {
	return rw.Client.Del(ctx, keys...).Result()
}

func (rw *RedisWrapper) Stream(ctx context.Context, stream string, group string) error {
	rw.Logger.Infof("Creating Redis stream: %s, Group: %s", stream, group)
	return rw.Client.XGroupCreate(ctx, stream, group, "$").Err()
}

func (rw *RedisWrapper) Entries(ctx context.Context, group string, stream string) ([]redis.XStream, error) {
	return rw.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: "*",
		Streams:  []string{stream, ">"},
		Count:    1,
		Block:    0,
		NoAck:    false,
	}).Result()
}

func Messages(entries []redis.XStream) []redis.XMessage {
	return entries[0].Messages
}

func (rw *RedisWrapper) Process(ctx context.Context, messages []redis.XMessage, stream string, group string) []MsgPost {
	var urls []MsgPost
	for _, message := range messages {
		eventName, href, group := processEntry(message.Values)
		if eventName == "receivedUrl" {
			msgPost := MsgPost{Href: href, Group: group}
			urls = append(urls, msgPost)
			rw.ackEntry(ctx, stream, group, message.ID) // Now it uses the method on the RedisWrapper
		}
	}
	return urls
}

func processEntry(values map[string]interface{}) (string, string, string) {
	eventName := fmt.Sprintf("%v", values["eventName"])
	href := fmt.Sprintf("%v", values["href"])
	group := fmt.Sprintf("%v", values["group"])

	return eventName, href, group
}

// Make sure ackEntry is correctly capitalized as it is a method of RedisWrapper
func (rw *RedisWrapper) ackEntry(ctx context.Context, stream string, group string, id string) error {
	return rw.Client.XAck(ctx, stream, group, id).Err()
}
