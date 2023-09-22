package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/jonesrussell/crawler/internal/post"
	"github.com/jonesrussell/crawler/internal/rediswrapper"
)

func main() {
	if godotenv.Load(".env") != nil {
		log.Fatal("error loading .env file")
	}

	// Setup the Redis connection
	addr := fmt.Sprintf(
		"%s:%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)
	redisClient := rediswrapper.Connect(addr, os.Getenv("REDIS_AUTH"))
	defer redisClient.Close()

	// Connect to Redis Stream
	err := rediswrapper.Stream(os.Getenv("REDIS_STREAM"), os.Getenv("REDIS_GROUP"))
	if err != nil {
		log.Println(err)
	}

	log.Println("consumer started")

	for {
		entries, err := rediswrapper.Entries(
			os.Getenv("REDIS_GROUP"),
			os.Getenv("REDIS_STREAM"),
		)
		if err != nil {
			log.Fatal(err)
		}

		messages := rediswrapper.Messages(entries)

		urls := rediswrapper.Process(
			messages,
			os.Getenv("REDIS_STREAM"),
			os.Getenv("REDIS_GROUP"),
		)

		consume(urls)
	}
}

func consume(urls []rediswrapper.MsgPost) {
	for i := 0; i < len(urls); i++ {
		href := urls[i]
		post.SetUsername(os.Getenv("USERNAME"))
		post.SetPassword(os.Getenv("PASSWORD"))
		err := post.Process(href, os.Getenv("API_FILTER_URL"))
		if err != nil {
			log.Fatal(err)
		}
	}
}
