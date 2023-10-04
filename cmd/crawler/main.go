// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/jonesrussell/crawler/internal/rediswrapper"
	termmatcher "github.com/jonesrussell/crawler/internal/termmatcher"
	"go.uber.org/zap"
)

type Config struct {
	URL         string
	SearchTerms string
	CrawlsiteID string
}

var logger *zap.SugaredLogger

func main() {
	// Create a logger
	initializeLogger()
	defer logger.Sync() // Flush the logger before exiting

	// Log the start of the main function
	logger.Info("Main function started...")

	// Retrieve URL to crawl and search terms from command line arguments
	config, err := parseCommandLineArguments()
	if err != nil {
		logger.Error("Error:", err)
		return // Return to exit the function gracefully
	}

	// Initialize Redis client
	rediswrapper.InitializeRedis(logger, os.Getenv("REDIS_HOST"), os.Getenv("REDIS_AUTH"))

	// Set the Crawlsite ID
	rediswrapper.SetCrawlsiteID(config.CrawlsiteID)

	crawlURL := config.URL
	searchTerms := strings.Split(config.SearchTerms, ",")

	// Log the URL being crawled
	logger.Info("Crawling URL:", crawlURL)

	// Log the search terms
	logger.Info("Search Terms:", searchTerms)

	// Load environment variables
	loadEnvironmentVariables()

	// Dynamically set allowed domain based on input URL
	allowedDomain := getHostFromURL(crawlURL)

	// Configure Colly collector with user agent and increased MaxDepth
	collector := configureCollector([]string{allowedDomain})
	// collector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"

	// Set up the crawling logic
	setupCrawlingLogic(collector, searchTerms)

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
	flag.StringVar(&config.CrawlsiteID, "crawlsite", "", "Crawlsite ID") // Add a new flag for Crawlsite ID
	flag.Parse()

	if config.URL == "" {
		return Config{}, fmt.Errorf("URL is required")
	}

	if config.CrawlsiteID == "" {
		return Config{}, fmt.Errorf("crawlsite id is required") // Ensure Crawlsite ID is provided
	}

	return config, nil
}

func loadEnvironmentVariables() {
	if godotenv.Load(".env") != nil {
		logger.Warn("Error loading .env file")
	}
}

func initializeLogger() {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.OutputPaths = []string{"stdout"} // Write logs to stdout
	zapLogger, err := loggerConfig.Build()
	if err != nil {
		log.Fatalf("Failed to initialize Zap logger: %v", err)
	}
	logger = zapLogger.Sugar()
}

func configureCollector(allowedDomains []string) *colly.Collector {
	collector := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(3),
	)

	collector.AllowedDomains = allowedDomains

	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       3000 * time.Millisecond,
	})

	return collector
}

func setupCrawlingLogic(collector *colly.Collector, searchTerms []string) {
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Request.AbsoluteURL(e.Attr("href"))

		if termmatcher.Related(href, searchTerms) {
			logger.Info("Found: ", href)

			if _, err := rediswrapper.SAdd(href); err != nil {
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

			err = rediswrapper.PublishHref("streetcode", href)
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
