package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

var (
	redis_uri    = ""
	redis_port   = ""
	redis_stream = ""
	ctx          = context.Background()
)

func main() {
	log.Println("crawler started")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}

	redis_uri = os.Getenv("REDIS_HOST")
	redis_port = os.Getenv("REDIS_PORT")
	redis_stream = os.Getenv("REDIS_STREAM")

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redis_uri, redis_port),
	})

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("unable to connect to redis", err)
	}

	log.Println("connected to redis")

	// Retrieve URL parameter
	webpage := os.Args[1]

	u, err := url.Parse(webpage)
	if err != nil {
		panic(err)
	}

	collector := colly.NewCollector(
		colly.AllowedDomains(u.Host),
		colly.Async(true),
		/*colly.URLFilters(
			// regexp.MustCompile(`https://www.sudbury.com$`),
			regexp.MustCompile(`(|/police.+)$`),
			// regexp.MustCompile(`https://www.sudbury\.com/membership.+`),
			// regexp.MustCompile(`https://www.sudbury.com/local-news.+`),
			// regexp.MustCompile(`https://www.sudbury\.com/weather.+`),
		),*/
		// colly.Debugger(&debug.LogDebugger{}),
	)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*sudbury.com" glob
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*sudbury.com",
		Parallelism: 2,
		Delay:       1000 * time.Millisecond,
	})

	// Act on every link; <a href="foo">
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		foundUrl := e.Request.AbsoluteURL(e.Attr("href"))

		// Determine if we will submit link to Redis
		if !strings.Contains(foundUrl, "/police/") {
			// log.Printf("INFO: %s not a candidate for Streetcode", foundUrl)
		} else {
			err = publishHrefReceivedEvent(redisClient, foundUrl)
			if err != nil {
				log.Fatal(err)
			}
		}

		collector.Visit(foundUrl)
	})

	collector.Visit(webpage)
	collector.Wait()
}

func publishHrefReceivedEvent(client *redis.Client, href string) error {
	log.Println("Publishing event to Redis")

	err := client.XAdd(ctx, &redis.XAddArgs{
		Stream:       redis_stream,
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
