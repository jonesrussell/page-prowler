package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/jonesrussell/crawler/internal/crawlResult"
	"github.com/jonesrussell/crawler/internal/rediswrapper"
	"github.com/jonesrussell/crawler/internal/stats"
	"github.com/jonesrussell/crawler/internal/termmatcher"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type CommandLineArgs struct {
	URL         string
	SearchTerms string
	CrawlSiteID string
	MaxDepth    int
	Debug       bool
}

type Config struct {
	URL         string
	SearchTerms string
	CrawlSiteID string
	RedisHost   string `envconfig:"REDIS_HOST"`
	RedisPort   string `envconfig:"REDIS_PORT"`
	RedisAuth   string `envconfig:"REDIS_AUTH"`
}

var logger *zap.SugaredLogger

func main() {
	args := processFlags()

	initializeLogger(args.Debug)
	defer func(logger *zap.SugaredLogger) {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("Error syncing logger: %v", err)
		}
	}(logger) // Flush the logger before exiting

	// Log the start of the main function
	logger.Info("Main function started...")

	var config Config
	config.URL = args.URL
	config.SearchTerms = args.SearchTerms
	config.CrawlSiteID = args.CrawlSiteID

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	err = envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	initializeRedis(config, args.Debug)

	logger.Info("Crawling URL:", config.URL)
	splitSearchTerms := strings.Split(config.SearchTerms, ",") // Use a new variable for the split search terms
	logger.Info("Search Terms:", splitSearchTerms)

	allowedDomain := getHostFromURL(config.URL)

	var results []crawlResult.PageData
	collector := configureCollector([]string{allowedDomain}, args.MaxDepth)
	setupCrawlingLogic(collector, splitSearchTerms, &results)

	logger.Info("Crawler started...")

	if err := collector.Visit(config.URL); err != nil {
		logger.Errorw("Error visiting URL", "error", err)
		return
	}

	collector.Wait()

	// After all crawling jobs are done, print the results
	jsonData, err := json.Marshal(results)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
	}
	fmt.Println(string(jsonData))

	logger.Info("Main function completed.")
}

func processFlags() CommandLineArgs {
	args := CommandLineArgs{}

	flag.StringVar(&args.URL, "url", "", "URL to crawl")
	flag.StringVar(&args.SearchTerms, "searchterms", "", "Comma-separated search terms")
	flag.StringVar(&args.CrawlSiteID, "crawlsiteid", "", "CrawlSite ID")
	flag.IntVar(&args.MaxDepth, "maxdepth", 1, "Maximum depth for the crawler")
	flag.BoolVar(&args.Debug, "debug", false, "Enable debug mode")

	flag.Parse()

	if args.URL == "" || args.SearchTerms == "" || args.CrawlSiteID == "" {
		fmt.Println("The following flags are required: url, searchterms, crawlsiteid")
		flag.PrintDefaults()
		os.Exit(2)
	}

	if args.Debug {
		fmt.Printf("URL: %s\n", args.URL)
		fmt.Printf("SearchTerms: %s\n", args.SearchTerms)
		fmt.Printf("CrawlSiteID: %s\n", args.CrawlSiteID)
		fmt.Printf("MaxDepth: %d\n", args.MaxDepth)
		fmt.Printf("Debug: %v\n", args.Debug)
	}

	return args
}

func initializeLogger(debug bool) {
	loggerConfig := zap.NewProductionConfig()
	if debug {
		loggerConfig.Level.SetLevel(zap.DebugLevel) // Log all messages in debug mode
	} else {
		loggerConfig.Level.SetLevel(zap.ErrorLevel) // Only log errors in non-debug mode
	}
	loggerConfig.OutputPaths = []string{"stdout"} // Write logs to stdout
	zapLogger, err := loggerConfig.Build()
	if err != nil {
		log.Fatalf("Failed to initialize Zap logger: %v", err)
	}
	logger = zapLogger.Sugar()
}

func initializeRedis(config Config, debug bool) {
	redisAddress := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)
	if debug {
		fmt.Printf("Redis Host: %s\n", config.RedisHost)
		fmt.Printf("Redis Port: %s\n", config.RedisPort)
		fmt.Printf("Redis Address: %s\n", redisAddress)
	}

	rediswrapper.InitializeRedis(logger, redisAddress, config.RedisAuth)
}

func configureCollector(allowedDomains []string, maxDepth int) *colly.Collector {
	collector := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(maxDepth),
	)

	collector.AllowedDomains = allowedDomains

	err := collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       3000 * time.Millisecond,
	})
	if err != nil {
		return nil
	}

	return collector
}

func handleHTMLParsing(collector *colly.Collector, searchTerms []string, linkStats *stats.Stats, results *[]crawlResult.PageData) (err error) {
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Request.AbsoluteURL(e.Attr("href"))

		linkStats.IncrementTotalLinks()

		pageData := crawlResult.PageData{
			URL: href,
			// Add other fields as necessary
		}

		if termmatcher.Related(href, searchTerms) {
			linkStats.IncrementMatchedLinks()
			err = handleMatchingLinks(collector, href)
			if err != nil {
				logger.Errorw("Error handling matching links", "error", err)
				return
			}
			pageData.MatchingTerms = searchTerms
		} else {
			linkStats.IncrementNotMatchedLinks()
			handleNonMatchingLinks(href)
		}

		*results = append(*results, pageData)
	})
	return
}

func handleMatchingLinks(collector *colly.Collector, href string) error {
	logger.Info("Found: ", href)
	if _, err := rediswrapper.SAdd(href); err != nil {
		logger.Errorw("Error adding URL to Redis set", "error", err)
		return err
	}
	// Visit the URL after adding it to the Redis set
	err := collector.Visit(href)
	if err != nil {
		return err
	} // Ignore any errors from visiting the URL
	return nil
}

// Handle non-matching links
func handleNonMatchingLinks(href string) {
	logger.Infof("Non-matching link: %s", href)
}

func handleRedisOperations() error {
	hrefs, err := rediswrapper.SMembers()
	if err != nil {
		logger.Errorw("Error getting members from Redis", "error", err)
		return err
	}

	for i := range hrefs {
		href := hrefs[i]

		err = rediswrapper.PublishHref("streetcode", href)
		if err != nil {
			logger.Errorw("Error publishing href to Redis", "error", err)
			return err
		}

		if _, err = rediswrapper.Del(); err != nil {
			logger.Errorw("Error deleting from Redis", "error", err)
			return err
		}
	}
	return nil
}

func handleErrorEvents(collector *colly.Collector) {
	collector.OnError(func(r *colly.Response, err error) {
		statusCode := r.StatusCode
		requestUrl := r.Request.URL.String()

		if statusCode != 404 {
			logger.Errorw("Request URL failed",
				"request_url", requestUrl,
				"status_code", statusCode,
			)
		}
	})
}

func setupCrawlingLogic(collector *colly.Collector, searchTerms []string, results *[]crawlResult.PageData) {
	linkStats := stats.NewStats()

	err := handleHTMLParsing(collector, searchTerms, linkStats, results)
	if err != nil {
		return
	}
	handleErrorEvents(collector)

	collector.OnScraped(func(r *colly.Response) {
		err := handleRedisOperations()
		if err != nil {
			return
		}

		logger.Info("Finished scraping the page:", r.Request.URL.String())

		logger.Info("Total links found:", linkStats.TotalLinks)
		logger.Info("Matched links:", linkStats.MatchedLinks)
		logger.Info("Not matched links:", linkStats.NotMatchedLinks)

		// Here, you would add code to populate the 'results' slice with data
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
