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
	fmt.Println("Main function started...")

	// Create a logger
	logger := createLogger()
	defer logger.Sync() // Flush the logger before exiting

	// Retrieve URL to crawl from arguments
	crawlURL, group, err := parseCommandLineArguments(os.Args)
	if err != nil {
		fmt.Println("Error:", err)
		return // Return to exit the function gracefully
	}

	fmt.Println("Crawling URL:", crawlURL)

	// Load environment variables
	loadEnvironmentVariables(logger)

	// Create and configure dependencies
	redisClient := createRedisClient()
	defer redisClient.Close()

	collector := configureCollector()

	// Set up the crawling logic
	setupCrawlingLogic(collector, logger, group)

	// Start crawling
	logger.Info("Crawler started...")
	collector.Visit(crawlURL)
	collector.Wait()

	fmt.Println("Main function completed.")
}

func parseCommandLineArguments(args []string) (string, string, error) {
	if len(args) < 3 {
		return "", "", fmt.Errorf("Usage: ./crawler https://www.sudbury.com c45fe232-0fbd-4fj8-b097-ff7bb863ae6b")
	}
	return args[1], args[2], nil
}

func loadEnvironmentVariables(logger *zap.SugaredLogger) {
	if godotenv.Load(".env") != nil {
		logger.Warn("Error loading .env file")
	}
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
	addr := fmt.Sprintf(
		"%s:%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)
	redisClient := rediswrapper.Connect(addr, os.Getenv("REDIS_AUTH"))
	return redisClient
}

func configureCollector() *colly.Collector {
	collector := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(3),
	)

	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       3000 * time.Millisecond,
	})

	return collector
}

func setupCrawlingLogic(collector *colly.Collector, logger *zap.SugaredLogger, group string) {
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Request.AbsoluteURL(e.Attr("href"))

		if drug.Related(href) {
			logger.Info(href)

			_, err := rediswrapper.SAdd(href)
			if err != nil {
				logger.Errorw("Error adding URL to Redis set", "error", err)
			}
		}

		if os.Getenv("CRAWL_MODE") != "single" {
			if e.Request.Depth < 1 {
				collector.Visit(href)
			}
		}
	})

	collector.OnScraped(func(r *colly.Response) {
		hrefs, err := rediswrapper.SMembers()
		if err != nil {
			log.Fatal(err)
		}

		for i := range hrefs {
			href := hrefs[i]

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

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	collector.OnError(func(r *colly.Response, err error) {
		logger.Errorw("Request URL failed",
			"request_url", r.Request.URL,
			"response", r,
			"error", err,
		)
	})
}
