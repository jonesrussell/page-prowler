package rediswrapper

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type MsgPost struct {
	Href  string `json:"href"`
	Group string `json:"group"`
}

var (
	ctx              = context.Background()
	client           = (*redis.Client)(nil)
	keySetBase       = "hrefs"
	crawlsiteID      = ""       // Store crawlsiteID as a package-level variable
	crawlsiteIDMutex sync.Mutex // Mutex to protect concurrent access to crawlsiteID
	logger           *zap.SugaredLogger
)

func SetCrawlsiteID(id string) {
	crawlsiteIDMutex.Lock()
	defer crawlsiteIDMutex.Unlock()
	crawlsiteID = id
}

func InitializeRedis(loggerInstance *zap.SugaredLogger, addr string, password string) {
	logger = loggerInstance

	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Fatal("Unable to connect to Redis: ", err)
	}
}

func SAdd(href string) (int64, error) {
	if client == nil {
		return 0, errors.New("Redis client is not initialized")
	}

	keySet := fmt.Sprintf("%s:%s", keySetBase, crawlsiteID)
	return client.SAdd(ctx, keySet, href).Result()
}

func SMembers() ([]string, error) {
	keySet := fmt.Sprintf("%s:%s", keySetBase, crawlsiteID)
	return client.SMembers(ctx, keySet).Result()
}

func Del() (int64, error) {
	keySet := fmt.Sprintf("%s:%s", keySetBase, crawlsiteID)
	return client.Del(ctx, keySet).Result()
}

func PublishHref(stream, href string) error {
	keyStream := fmt.Sprintf("%s:%s", stream, crawlsiteID)
	logger.Infof("Publishing to stream %s", keyStream)

	err := client.XAdd(ctx, &redis.XAddArgs{
		Stream:       keyStream,
		MaxLen:       0,
		MaxLenApprox: 0,
		ID:           "",
		Values: map[string]interface{}{
			"eventName": "receivedUrl",
			"href":      href,
		},
	}).Err()

	if err != nil {
		logger.Infof("Error publishing to stream %s: %v", keyStream, err)
	}

	return err
}

func Stream(stream string, group string) error {
	return client.XGroupCreate(
		ctx,
		stream,
		group,
		"0",
	).Err()
}

func Entries(group string, stream string) ([]redis.XStream, error) {
	return client.XReadGroup(ctx, &redis.XReadGroupArgs{
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

func Process(messages []redis.XMessage, stream string, group string) []MsgPost {
	var urls []MsgPost

	for i := 0; i < len(messages); i++ {
		eventName, href, group := processEntry(messages[i].Values)

		if eventName == "receivedUrl" {
			msgPost := MsgPost{
				Href:  href,
				Group: group,
			}

			urls = append(urls, msgPost)
			ackEntry(stream, group, messages[i].ID)
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

func ackEntry(stream string, group string, id string) {
	client.XAck(
		ctx,
		stream,
		group,
		id,
	)
}
