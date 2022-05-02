package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/jonesrussell/crawler/internal/myredis"
	"github.com/jonesrussell/crawler/internal/post"
)

func main() {
	if godotenv.Load(".env") != nil {
		log.Fatal("error loading .env file")
	}

	// Connect to Redis
	redisClient := myredis.Connect()
	defer redisClient.Close()

	// Connect to Redis Stream
	err := myredis.Stream()
	if err != nil {
		log.Println(err)
	}

	log.Println("consumer started")

	for {
		entries, err := myredis.GetEntries()
		if err != nil {
			log.Fatal(err)
		}

		processEntries(entries)
	}
}

func processEntries(entries []redis.XStream) {
	messages := entries[0].Messages

	for i := 0; i < len(messages); i++ {
		processEntry(messages[i].Values, messages[i].ID)
	}
}

func processEntry(values map[string]interface{}, id string) {
	eventName := fmt.Sprintf("%v", values["eventName"])
	href := fmt.Sprintf("%v", values["href"])

	if eventName == "receivedUrl" {
		err := post.ProcessHref(href)
		if err != nil {
			log.Fatal(err)
		}

		myredis.AckEntry(id)
	}
}
