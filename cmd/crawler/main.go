package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/jonesrussell/crawler/internal/drug"
	"github.com/jonesrussell/crawler/internal/rediswrapper"
	"go.uber.org/zap"
)

func main() {
	// Create a logger
	logger := createLogger()
	defer logger.Sync() // Flush the logger before exiting

	// Retrieve URL to crawl from arguments
	if len(os.Args) < 3 {
		logger.Info("Usage: ./crawler https://www.sudbury.com c45fe232-0fbd-4fj8-b097-ff7bb863ae6b")
		os.Exit(0)
	}
	crawlUrl := os.Args[1]
	group := os.Args[2]

	// Load the environment variables
	if godotenv.Load(".env") != nil {
		logger.Warn("Error loading .env file")
	}

	// Create the Redis client and connection
	redisClient := createRedisClient()
	defer redisClient.Close()

	// Create a new crawler
	collector := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(3),
	)

	// Set reasonable limits for responsible crawling
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       3000 * time.Millisecond,
	})

	// When a link is found
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// Extract the full URL
		href := e.Request.AbsoluteURL(e.Attr("href"))

		// Determine if we will submit the link to Redis
		if drug.Related(href) {
			// Announce the drug-related URL
			logger.Info(href)

			// Add URL to the publishing queue
			_, err := rediswrapper.SAdd(href)
			if err != nil {
				logger.Errorw("Error adding URL to Redis set", "error", err)
			}
		}

		if os.Getenv("CRAWL_MODE") != "single" {
			// Check the depth before visiting the link
			if e.Request.Depth < 1 {
				collector.Visit(href)
			}
		}
	})

	// When a url has finished being crawled
	collector.OnScraped(func(r *colly.Response) {
		// Retrieve the urls to be published
		hrefs, err := rediswrapper.SMembers()
		if err != nil {
			log.Fatal(err)
		}

		// Loop over urls
		for i := range hrefs {
			href := hrefs[i]

			// Send url to Redis stream
			err = rediswrapper.PublishHref(os.Getenv("REDIS_STREAM"), href, group)
			if err != nil {
				log.Fatal(err)
			}

			_, err = rediswrapper.Del()
			if err != nil {
				log.Fatal(err)
			}

		}
	})

	// When crawling a url
	collector.OnRequest(func(r *colly.Request) {
		// Announce the url
		fmt.Println("Visiting", r.URL)
	})

	// Set error handler
	collector.OnError(func(r *colly.Response, err error) {
		logger.Errorw("Request URL failed",
			"request_url", r.Request.URL,
			"response", r,
			"error", err,
		)
	})

	// Everything is setup, time to crawl
	logger.Info("Crawler started...")
	collector.Visit(crawlUrl)
	collector.Wait()
}

func createLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize Zap logger: %v", err)
	}
	sugar := logger.Sugar()
	return sugar
}

func createRedisClient() *redis.Client {
	// Setup the Redis connection
	addr := fmt.Sprintf(
		"%s:%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)
	redisClient := rediswrapper.Connect(addr, os.Getenv("REDIS_AUTH"))
	return redisClient
}
