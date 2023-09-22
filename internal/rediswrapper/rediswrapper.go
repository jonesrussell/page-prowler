package rediswrapper

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type MsgPost struct {
	Href  string `json:"href"`
	Group string `json:"group"`
}

var (
	ctx    = context.Background()
	client = (*redis.Client)(nil)
)

const keySet = "hrefs"

func Connect(addr string, password string) *redis.Client {
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(" unbale to connect to Redis ", err)
	}

	return client
}

func SAdd(href string) (int64, error) {
	return client.SAdd(ctx, keySet, href).Result()
}

func SMembers() ([]string, error) {
	return client.SMembers(ctx, keySet).Result()
}

func Del() (int64, error) {
	return client.Del(ctx, keySet).Result()
}

func PublishHref(stream string, href string, group string) error {
	return client.XAdd(ctx, &redis.XAddArgs{
		Stream:       stream,
		MaxLen:       0,
		MaxLenApprox: 0,
		ID:           "",
		Values: map[string]interface{}{
			"eventName": "receivedUrl",
			"href":      href,
			"group":     group,
		},
	}).Err()
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
