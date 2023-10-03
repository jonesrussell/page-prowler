package main

import (
	"flag"
	"fmt"
	"log"
	"net/url" // Import the net/url package
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

type Config struct {
	URL         string
	SearchTerms string
}

func main() {
	// Create a logger
	logger := createLogger()
	defer logger.Sync() // Flush the logger before exiting

	// Log the start of the main function
	logger.Info("Main function started...")

	// Retrieve URL to crawl and search terms from command line arguments
	config, err := parseCommandLineArguments()
	if err != nil {
		logger.Error("Error:", err)
		return // Return to exit the function gracefully
	}

	crawlURL := config.URL
	searchTerms := strings.Split(config.SearchTerms, ",")

	// Log the URL being crawled
	logger.Info("Crawling URL:", crawlURL)

	// Log the search terms
	logger.Info("Search Terms:", searchTerms)

	// Load environment variables
	loadEnvironmentVariables(logger)

	// Create and configure dependencies
	redisClient := createRedisClient()
	defer redisClient.Close()

	// Dynamically set allowed domain based on input URL
	allowedDomain := getHostFromURL(crawlURL)

	// Configure Colly collector with user agent and increased MaxDepth
	collector := configureCollector([]string{allowedDomain}, logger)
	// collector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"

	// Set up the crawling logic
	setupCrawlingLogic(collector, logger, searchTerms)

	// Start crawling
	logger.Info("Crawler started...")
	collector.Visit(crawlURL)
	collector.Wait()

	// Log the completion of the main function
	logger.Info("Main function completed.")
}

func parseCommandLineArguments() (Config, error) {
	var config Config

	flag.StringVar(&config.URL, "url", "", "URL to crawl")
	flag.StringVar(&config.SearchTerms, "search", "", "Search terms (comma-separated)")
	flag.Parse()

	if config.URL == "" {
		return Config{}, fmt.Errorf("URL is required")
	}

	return config, nil
}

func loadEnvironmentVariables(logger *zap.SugaredLogger) {
	if godotenv.Load(".env") != nil {
		logger.Warn("Error loading .env file")
	}
}

func createLogger() *zap.SugaredLogger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.OutputPaths = []string{"stdout"} // Write logs to stdout
	logger, err := loggerConfig.Build()
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

func configureCollector(allowedDomains []string, logger *zap.SugaredLogger) *colly.Collector {
	collector := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(3),
	)

	// Log the allowed domains using the logger
	logger.Info("Allowed Domains:", allowedDomains)

	// Set allowed domains based on the provided domains
	collector.AllowedDomains = allowedDomains

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
			logger.Info("Found: ", href)

			_, err := rediswrapper.SAdd(href)
			if err != nil {
				logger.Errorw("Error adding URL to Redis set", "error", err)
			} else {
				// Visit the URL after adding it to the Redis set
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
		logger.Info("Visiting: ", r.URL)
	})

	collector.OnError(func(r *colly.Response, err error) {
		// Extract the status code and URL
		statusCode := r.StatusCode
		url := r.Request.URL.String()

		logger.Errorw("Request URL failed",
			"request_url", url,
			"status_code", statusCode,
		)
	})

}

func getHostFromURL(inputURL string) string {
	u, err := url.Parse(inputURL)
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
	}
	return u.Host
}
