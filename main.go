package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

func main() {
	// Retrieve URL to crawl
	if len(os.Args) < 2 {
		log.Fatalln("usage: ./streetcode-crawler https://www.sudbury.com")
	}
	crawlUrl := os.Args[1]

	if godotenv.Load(".env") != nil {
		log.Fatal("error loading .env file")
	}

	// Connect to Redis
	redisAddress := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: os.Getenv("REDIS_AUTH"),
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("unable to connect to redis", err)
	}

	log.Println("crawler started")

	collector := colly.NewCollector(
		colly.Async(true),
		colly.URLFilters(
			regexp.MustCompile(`(|/police.+)$`),
		),
		// colly.Debugger(&debug.LogDebugger{}),
	)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*sudbury.com" glob
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*sudbury.com",
		Parallelism: 2,
		Delay:       3000 * time.Millisecond,
	})

	// Act on every link; <a href="foo">
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		foundHref := e.Request.AbsoluteURL(e.Attr("href"))

		// Determine if we will submit link to Redis
		if strings.Contains(foundHref, "/police/") {
			if err = publishHref(redisClient, foundHref); err != nil {
				log.Fatal(err)
			}
			// TODO: put this behing a cli flag or env var
			// writeHrefCsv(foundHref)
		}

		collector.Visit(foundHref)
	})

	collector.Visit(crawlUrl)
	collector.Wait()
}

func writeHrefCsv(href string) {
	f, err := os.OpenFile(os.Getenv("CSV_FILENAME"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	w := csv.NewWriter(f)
	w.Write([]string{href})
	w.Flush()
}

func publishHref(client *redis.Client, href string) error {
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
