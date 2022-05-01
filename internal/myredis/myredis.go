package myredis

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func Connect() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_AUTH"),
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal(" unbale to connect to Redis ", err)
	}

	return redisClient
}

func SAdd(redisClient *redis.Client, href string) (delta time.Duration) {
	key := "schref"
	t0 := time.Now()
	redisClient.SAdd(ctx, key, href)
	delta = time.Since(t0)
	redisClient.FlushDB(ctx)
	return
}

func PublishHref(client *redis.Client, href string) error {
	err := client.XAdd(ctx, &redis.XAddArgs{
		Stream:       os.Getenv("REDIS_STREAM"),
		MaxLen:       0,
		MaxLenApprox: 0,
		ID:           "",
		Values: map[string]interface{}{
			"eventName": string("href received"),
			"href":      href,
		},
	}).Err()

	return err
}
