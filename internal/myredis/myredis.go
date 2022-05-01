package myredis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	ctx    = context.Background()
	client = (*redis.Client)(nil)
)

const key = "hrefs"

func Connect() *redis.Client {
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		),
		Password: os.Getenv("REDIS_AUTH"),
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(" unbale to connect to Redis ", err)
	}

	return client
}

func SAdd(href string) (int64, error) {
	return client.SAdd(ctx, key, href).Result()
}

func SPop() (string, error) {
	return client.SPop(ctx, key).Result()
}

func PublishHref(href string) error {
	return client.XAdd(ctx, &redis.XAddArgs{
		Stream:       os.Getenv("REDIS_STREAM"),
		MaxLen:       0,
		MaxLenApprox: 0,
		ID:           "",
		Values: map[string]interface{}{
			"eventName": "receivedUrl",
			"href":      href,
		},
	}).Err()
}

func Stream() error {
	return client.XGroupCreate(
		ctx,
		os.Getenv("REDIS_STREAM"),
		os.Getenv("REDIS_GROUP"),
		"0",
	).Err()
}

func GetEntries() ([]redis.XStream, error) {
	return client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    os.Getenv("REDIS_GROUP"),
		Consumer: "*",
		Streams:  []string{os.Getenv("REDIS_STREAM"), ">"},
		Count:    1,
		Block:    0,
		NoAck:    false,
	}).Result()
}

func AckEntry(id string) {
	client.XAck(
		ctx,
		os.Getenv("REDIS_STREAM"),
		os.Getenv("REDIS_GROUP"),
		id,
	)
}
