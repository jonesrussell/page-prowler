package main

import (
	"log"

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
		entries, err := myredis.Entries()
		if err != nil {
			log.Fatal(err)
		}

		messages := myredis.Messages(entries)

		urls := myredis.Process(messages)

		consume(urls)
	}
}

func consume(urls []string) {
	for i := 0; i < len(urls); i++ {
		err := post.ProcessHref(urls[i])
		if err != nil {
			log.Fatal(err)
		}
	}
}
