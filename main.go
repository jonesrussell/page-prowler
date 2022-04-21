package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gocolly/colly"
)

const (
	redis_uri    = "localhost"
	redis_port   = "6379"
	redis_stream = "streetcode"
)

var ctx = context.Background()

func main() {
	log.Println("Publisher started")

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redis_uri, redis_port),
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Unable to connect to Redis", err)
	}

	log.Println("Connected to Redis server")

	// Retrieve URL parameter
	webpage := os.Args[1]
	u, err := url.Parse(webpage)
	if err != nil {
		panic(err)
	}

	collector := colly.NewCollector(
		colly.AllowedDomains(u.Host),
		colly.Async(true),
		colly.URLFilters(
			// regexp.MustCompile(`https://www.sudbury.com$`),
			regexp.MustCompile(`https://www.sudbury.com(|/police.+)$`),
			// regexp.MustCompile(`https://www.sudbury\.com/membership.+`),
			// regexp.MustCompile(`https://www.sudbury.com/local-news.+`),
			// regexp.MustCompile(`https://www.sudbury\.com/weather.+`),
		),
		// colly.Debugger(&debug.LogDebugger{}),
	)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*sudbury.com" glob
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*sudbury.com",
		Parallelism: 2,
		Delay:       1000 * time.Millisecond,
	})

	collector.OnHTML("a.section-item", func(element *colly.HTMLElement) {
		href := element.Attr("href")

		if strings.Contains(href, "/police") {
			err = publishHrefReceivedEvent(redisClient, href)
			if err != nil {
				log.Fatal(err)
			}
		}

	})

	collector.OnHTML(`a[href]`, func(e *colly.HTMLElement) {
		foundURL := e.Request.AbsoluteURL(e.Attr("href"))
		collector.Visit(foundURL)
	})

	/*collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})*/

	fmt.Println("Webpage:", webpage)
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
