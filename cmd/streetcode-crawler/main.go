package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/jonesrussell/crawler/internal/drug"
)

var ctx = context.Background()

var redisClient = (*redis.Client)(nil)

func main() {
	// Retrieve URL to crawl
	if len(os.Args) < 2 {
		log.Fatalln("usage: ./streetcode-crawler https://www.sudbury.com")
	}
	crawlUrl := os.Args[1]

	if godotenv.Load(".env") != nil {
		log.Println("error loading .env file")
	}

	// Connect to Redis
	redisAddress := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	redisClient = redis.NewClient(&redis.Options{
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
			regexp.MustCompile("https://www.sudbury.com/police"),
			regexp.MustCompile("https://www.midnorthmonitor.com/category/news"),
			regexp.MustCompile("https://www.midnorthmonitor.com/news"),
			regexp.MustCompile("https://www.midnorthmonitor.com/category/news/local-news"),
		),
		// colly.Debugger(&debug.LogDebugger{}),
	)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*sudbury.com" glob
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       3000 * time.Millisecond,
	})

	// Act on every link; <a href="foo">
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		foundHref := e.Request.AbsoluteURL(e.Attr("href"))
		fmt.Println(foundHref)

		// Determine if we will submit link to Redis
		matchedNewsMnm, _ := regexp.MatchString(`^https://www.midnorthmonitor.com/news/`, foundHref)
		matchedPoliceSc, _ := regexp.MatchString(`^https://www.sudbury.com/police/`, foundHref)

		if matchedPoliceSc || matchedNewsMnm {
			if drug.Related(foundHref) {
				doSAdd(foundHref)
			}

			/*if err = publishHref(redisClient, foundHref); err != nil {
				log.Fatal(err)
			}

			if os.Getenv("CSV_WRITE") == "true" {
				writeHrefCsv(foundHref)
			}*/
		}

		if os.Getenv("CRAWL_MODE") != "single" {
			collector.Visit(foundHref)
		}
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// Set error handler
	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	collector.Visit(crawlUrl)
	collector.Wait()
}

func writeHrefCsv(href string) {
	f, err := os.OpenFile(os.Getenv("CSV_FILENAME"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
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

func doSAdd(href string) (delta time.Duration) {
	key := "schref"
	t0 := time.Now()
	redisClient.SAdd(ctx, key, href)
	delta = time.Since(t0)
	redisClient.FlushDB(ctx)
	return
}
