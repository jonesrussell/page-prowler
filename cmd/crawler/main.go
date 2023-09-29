package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/jonesrussell/crawler/internal/rediswrapper"
	termmatcher "github.com/jonesrussell/crawler/internal/termmatcher"
	"go.uber.org/zap"
)

func main() {
	// Create a logger
	logger := createLogger()
	defer logger.Sync() // Flush the logger before exiting

	// Log the start of the main function
	logger.Info("Main function started...")

	// Retrieve URL to crawl and search terms from command line arguments
	crawlURL, searchTerms, err := parseCommandLineArguments()
	if err != nil {
		logger.Error("Error:", err)
		return // Return to exit the function gracefully
	}

	// Log the URL being crawled
	logger.Info("Crawling URL:", crawlURL)

	// Load environment variables
	loadEnvironmentVariables(logger)

	// Create and configure dependencies
	redisClient := createRedisClient()
	defer redisClient.Close()

	collector := configureCollector()

	// Set up the crawling logic
	setupCrawlingLogic(collector, logger, searchTerms)

	// Start crawling
	logger.Info("Crawler started...")
	collector.Visit(crawlURL)
	collector.Wait()

	// Log the completion of the main function
	logger.Info("Main function completed.")
}

func parseCommandLineArguments() (string, []string, error) {
	// Define flags for URL and search terms
	var crawlURL string
	var searchTerm string

	// Parse the command-line flags
	flag.StringVar(&crawlURL, "url", "", "URL to crawl")
	flag.StringVar(&searchTerm, "search", "", "Search terms (comma-separated)")
	flag.Parse()

	// Check if the URL flag is empty
	if crawlURL == "" {
		return "", nil, fmt.Errorf("URL is required")
	}

	// Split search terms if provided
	var searchTerms []string
	if searchTerm != "" {
		searchTerms = strings.Split(searchTerm, ",")
	}

	return crawlURL, searchTerms, nil
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

func setupCrawlingLogic(collector *colly.Collector, logger *zap.SugaredLogger, searchTerms []string) {
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Request.AbsoluteURL(e.Attr("href"))

		if termmatcher.Related(href, searchTerms) {
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

			err = rediswrapper.PublishHref(os.Getenv("REDIS_STREAM"), href)
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
