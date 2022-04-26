package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gocolly/colly"
	"github.com/gocolly/redisstorage"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

func main() {
	log.Println("crawler started")

	// Retrieve URL parameter
	if len(os.Args) < 2 {
		log.Fatal("url not provided. eg) ./streetcode-crawler https://www.sudbury.com")
	}

	webpage := os.Args[1]

	if godotenv.Load(".env") != nil {
		log.Fatal("error loading .env file")
	}

	redis_db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("unable to connect to redis", err)
	}

	log.Println("connected to redis")

	collector := colly.NewCollector(
		colly.Async(true),
		colly.URLFilters(
			regexp.MustCompile(`(|/police.+)$`),
		),
		// colly.Debugger(&debug.LogDebugger{}),
	)

	// Setup Redis as colly cookie storage
	storage := &redisstorage.Storage{
		Address:  fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_AUTH"),
		DB:       redis_db,
		Prefix:   os.Getenv("REDIS_STREAM"),
	}

	// add storage to the collector
	if err := collector.SetStorage(storage); err != nil {
		panic(err)
	}

	// delete previous data from storage
	if err := storage.Clear(); err != nil {
		log.Fatal(err)
	}

	// close redis client
	defer storage.Client.Close()

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*sudbury.com" glob
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*sudbury.com",
		Parallelism: 2,
		Delay:       1000 * time.Millisecond,
	})

	// Act on every link; <a href="foo">
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		foundHref := e.Request.AbsoluteURL(e.Attr("href"))

		// Determine if we will submit link to Redis
		if !strings.Contains(foundHref, "/police/") {
			// log.Printf("INFO: %s not a candidate for Streetcode", foundUrl)
		} else {
			err = publishHrefReceivedEvent(redisClient, foundHref)
			if err != nil {
				log.Fatal(err)
			}
			writeHrefCsv(foundHref)
		}

		collector.Visit(foundHref)
	})

	collector.Visit(webpage)
	collector.Wait()
}

func writeHrefCsv(href string) {
	f, err := os.OpenFile("hrefs.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := csv.NewWriter(f)
	w.Write([]string{href})
	w.Flush()
}

func publishHrefReceivedEvent(client *redis.Client, href string) error {
	log.Println("Publishing event to Redis")

	err := client.XAdd(ctx, &redis.XAddArgs{
		Stream:       os.Getenv("REDIS_STREAM"),
		MaxLen:       0,
		MaxLenApprox: 0,
		ID:           "",
		Values: map[string]interface{}{
			"whatHappened": string("href received"),
			"href":         href,
		},
	}).Err()

	return err
}
