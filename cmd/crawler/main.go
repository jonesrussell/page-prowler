package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/jonesrussell/crawler/internal/rediswrapper"
	"github.com/jonesrussell/crawler/internal/stats"
	termmatcher "github.com/jonesrussell/crawler/internal/termmatcher"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type Config struct {
	URL         string
	SearchTerms string
	CrawlsiteID string
	RedisHost   string `envconfig:"REDIS_HOST"`
	RedisPort   string `envconfig:"REDIS_PORT"`
	RedisAuth   string `envconfig:"REDIS_AUTH"`
}

var logger *zap.SugaredLogger

func main() {
	// Create a logger
	initializeLogger()
	defer logger.Sync() // Flush the logger before exiting

	// Log the start of the main function
	logger.Info("Main function started...")

	// Define flags
	urlFlag := flag.String("url", "", "URL to crawl")
	searchTermsFlag := flag.String("searchterms", "", "Comma-separated search terms")
	crawlsiteIDFlag := flag.String("crawlsiteid", "", "Crawlsite ID")

	flag.Parse()

	// Check if flags are set
	if *urlFlag == "" || *searchTermsFlag == "" || *crawlsiteIDFlag == "" {
		log.Fatal("url, searchterms, and crawlsiteid are required")
	}

	var config Config
	config.URL = *urlFlag
	config.SearchTerms = *searchTermsFlag
	config.CrawlsiteID = *crawlsiteIDFlag

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	err = envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	initializeRedis(config)

	logger.Info("Crawling URL:", config.URL)
	searchTerms := strings.Split(config.SearchTerms, ",")
	logger.Info("Search Terms:", searchTerms)

	allowedDomain := getHostFromURL(config.URL)

	collector := configureCollector([]string{allowedDomain})
	setupCrawlingLogic(collector, searchTerms)

	logger.Info("Crawler started...")
	collector.Visit(config.URL)
	collector.Wait()

	logger.Info("Main function completed.")
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

func initializeRedis(config Config) {
	redisAddress := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)
	fmt.Printf("Redis Host: %s\n", config.RedisHost)
	fmt.Printf("Redis Port: %s\n", config.RedisPort)
	fmt.Printf("Redis Address: %s\n", redisAddress)
	rediswrapper.InitializeRedis(logger, redisAddress, config.RedisAuth)
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

// Handle HTML parsing and link extraction
func handleHTMLParsing(collector *colly.Collector, searchTerms []string, linkStats *stats.Stats) {
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Request.AbsoluteURL(e.Attr("href"))

		linkStats.IncrementTotalLinks()

		if termmatcher.Related(href, searchTerms) {
			linkStats.IncrementMatchedLinks()
			handleMatchingLinks(collector, href)
		} else {
			linkStats.IncrementNotMatchedLinks()
			handleNonMatchingLinks(href)
		}
	})
}

// Handle matching links with search terms
func handleMatchingLinks(collector *colly.Collector, href string) {
	logger.Info("Found: ", href)
	if _, err := rediswrapper.SAdd(href); err != nil {
		logger.Errorw("Error adding URL to Redis set", "error", err)
	} else {
		// Visit the URL after adding it to the Redis set
		collector.Visit(href)
	}
}

// Handle non-matching links
func handleNonMatchingLinks(href string) {
	// You can add logic here if needed
}

// Handle Redis operations
func handleRedisOperations() {
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
}

// Handle error events
func handleErrorEvents(collector *colly.Collector) {
	collector.OnError(func(r *colly.Response, err error) {
		statusCode := r.StatusCode
		url := r.Request.URL.String()

		logger.Errorw("Request URL failed",
			"request_url", url,
			"status_code", statusCode,
		)
	})
}

// The refactored setupCrawlingLogic function
func setupCrawlingLogic(collector *colly.Collector, searchTerms []string) {
	linkStats := stats.NewStats()

	handleHTMLParsing(collector, searchTerms, linkStats)
	handleErrorEvents(collector)

	collector.OnScraped(func(r *colly.Response) {
		handleRedisOperations()

		logger.Info("Finished scraping the page:", r.Request.URL.String())

		logger.Info("Total links found:", linkStats.TotalLinks)
		logger.Info("Matched links:", linkStats.MatchedLinks)
		logger.Info("Not matched links:", linkStats.NotMatchedLinks)

	})

	collector.OnRequest(func(r *colly.Request) {
		logger.Info("Visiting: ", r.URL)
	})
}

func getHostFromURL(inputURL string) string {
	u, err := url.Parse(inputURL)
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
	}
	return u.Host
}
