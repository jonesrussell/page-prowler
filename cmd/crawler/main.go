package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/jonesrussell/crawler/internal/crawlResult"
	"github.com/jonesrussell/crawler/internal/rediswrapper"
	"github.com/jonesrussell/crawler/internal/stats"
	"github.com/jonesrussell/crawler/internal/termmatcher"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CommandLineArgs struct {
	URL         string
	SearchTerms string
	CrawlSiteID string
	MaxDepth    int
	Debug       bool
}

// Define your struct that matches the environment variables
type EnvConfig struct {
	RedisHost string `envconfig:"REDIS_HOST"`
	RedisPort string `envconfig:"REDIS_PORT"`
	RedisAuth string `envconfig:"REDIS_AUTH"`
}

func main() {
	ctx := context.Background()
	args := processFlags()

	logger, err := initializeLogger(args.Debug)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer safeLoggerSync(logger)

	envCfg, err := loadConfiguration()
	if err != nil {
		logger.Fatal("Failed to load environment configuration: ", zap.Error(err))
	}

	if args.Debug {
		logger.Infof("Redis Host: %s", envCfg.RedisHost)
		logger.Infof("Redis Port: %s", envCfg.RedisPort)
		logger.Infof("Redis Auth: %s", envCfg.RedisAuth)
	}

	// Initialize Redis with the context
	redisWrapper, err := rediswrapper.NewRedisWrapper(ctx, envCfg.RedisHost, envCfg.RedisPort, envCfg.RedisAuth, logger)
	if err != nil {
		logger.Fatal("Failed to initialize Redis", zap.Error(err))
	}

	setupAndStartCrawler(ctx, args, redisWrapper, logger)
}

func loadConfiguration() (*EnvConfig, error) {
	var cfg EnvConfig

	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("error processing environment variables: %w", err)
	}

	return &cfg, nil
}

func initializeLogger(debug bool) (*zap.SugaredLogger, error) {
	var logger *zap.Logger
	var err error

	if debug {
		// Development logger is more verbose and writes to standard output
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // optional, colorizes the output
		logger, err = config.Build()
	} else {
		// Production logger is less verbose and could be set to log to a file
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	return logger.Sugar(), nil
}

func safeLoggerSync(logger *zap.SugaredLogger) {
	if err := logger.Sync(); err != nil {
		// As of Zap 1.x, Sync() returns an error on Windows, which can be ignored.
		// You can log the error if you're running on other platforms or if you want to be thorough.
		if zapErr, ok := err.(*os.PathError); !ok || zapErr.Err != syscall.EINVAL {
			log.Printf("Error syncing logger: %v", err)
		}
	}
}

func setupAndStartCrawler(
	ctx context.Context,
	args CommandLineArgs,
	redisWrapper *rediswrapper.RedisWrapper,
	logger *zap.SugaredLogger,
) {
	// Split search terms
	splitSearchTerms := strings.Split(args.SearchTerms, ",")

	// Configure the collector
	collector := configureCollector([]string{getHostFromURL(args.URL)}, args.MaxDepth)
	if collector == nil {
		logger.Fatal("Failed to configure collector")
		return
	}

	// Setup crawling logic
	var results []crawlResult.PageData
	// Pass args.CrawlSiteID to setupCrawlingLogic along with other parameters
	setupCrawlingLogic(ctx, args.CrawlSiteID, collector, splitSearchTerms, &results, logger, redisWrapper)

	// Start the crawling process
	logger.Info("Crawler started...")
	if err := collector.Visit(args.URL); err != nil {
		logger.Error("Error visiting URL", zap.Error(err))
		return
	}

	// Wait for crawling to complete
	collector.Wait()

	// Handle the results after crawling is done
	jsonData, err := json.Marshal(results)
	if err != nil {
		logger.Error("Error occurred during marshaling", zap.Error(err))
		return
	}

	// Output or process the results as needed
	fmt.Println(string(jsonData))

	logger.Info("Crawling completed.")
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

func handleHTMLParsing(
	ctx context.Context,
	crawlSiteID string, // Include crawlSiteID as a parameter
	logger *zap.SugaredLogger,
	collector *colly.Collector,
	searchTerms []string,
	linkStats *stats.Stats,
	results *[]crawlResult.PageData,
	redisWrapper *rediswrapper.RedisWrapper,
) error { // Removed the named return value 'err' to avoid confusion with the 'err' inside the callback
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Request.AbsoluteURL(e.Attr("href"))
		linkStats.IncrementTotalLinks()

		pageData := crawlResult.PageData{
			URL: href,
			// Add other fields as necessary
		}

		if termmatcher.Related(href, searchTerms) {
			linkStats.IncrementMatchedLinks()
			// Correctly pass the context and crawlSiteID to handleMatchingLinks
			if err := handleMatchingLinks(ctx, logger, collector, href, crawlSiteID, redisWrapper); err != nil {
				logger.Error("Error handling matching links", zap.Error(err))
				// Do not return from the callback, as it will not propagate to the outer function
			}
			pageData.MatchingTerms = searchTerms
		} else {
			linkStats.IncrementNotMatchedLinks()
			handleNonMatchingLinks(logger, href)
		}

		*results = append(*results, pageData)
	})
	return nil // Return nil here as any errors inside the callback do not propagate out
}

func handleMatchingLinks(
	ctx context.Context,
	logger *zap.SugaredLogger,
	collector *colly.Collector,
	href string,
	crawlSiteID string, // Assuming this is passed into the function
	redisWrapper *rediswrapper.RedisWrapper,
) error {
	logger.Info("Found: ", zap.String("url", href))
	if _, err := redisWrapper.SAdd(ctx, crawlSiteID, href); err != nil { // Use CrawlSiteID as the key name
		logger.Error("Error adding URL to Redis set", zap.String("set", crawlSiteID), zap.Error(err))
		return err
	}
	// Visit the URL after adding it to the Redis set
	err := collector.Visit(href)
	if err != nil {
		logger.Error("Error visiting URL", zap.String("url", href), zap.Error(err))
		return err
	}
	return nil
}

// Handle non-matching links
func handleNonMatchingLinks(logger *zap.SugaredLogger, href string) {
	logger.Info("Non-matching link: ", zap.String("url", href))
}

func handleRedisOperations(ctx context.Context, redisWrapper *rediswrapper.RedisWrapper, logger *zap.SugaredLogger) error {
	// You need to pass the context and the appropriate key to SMembers
	hrefs, err := redisWrapper.SMembers(ctx, "yourKeyHere") // Replace "yourKeyHere" with the actual key you're interested in
	if err != nil {
		logger.Error("Error getting members from Redis", zap.Error(err))
		return err
	}

	for _, href := range hrefs {
		err = redisWrapper.PublishHref(ctx, "streetcode", href) // Make sure to pass ctx to PublishHref if it requires it
		if err != nil {
			logger.Error("Error publishing href to Redis", zap.Error(err))
			return err
		}

		if _, err = redisWrapper.Del(ctx, href); err != nil { // Make sure to pass ctx to Del if it requires it
			logger.Error("Error deleting from Redis", zap.Error(err))
			return err
		}
	}
	return nil
}

func handleErrorEvents(collector *colly.Collector, logger *zap.SugaredLogger) {
	collector.OnError(func(r *colly.Response, err error) {
		statusCode := r.StatusCode
		requestURL := r.Request.URL.String()

		if statusCode != 404 {
			logger.Error("Request URL failed",
				zap.String("request_url", requestURL),
				zap.Int("status_code", statusCode),
				zap.Error(err),
			)
		}
	})
}

func setupCrawlingLogic(
	ctx context.Context,
	crawlSiteID string, // Added CrawlSiteID as a parameter
	collector *colly.Collector,
	searchTerms []string,
	results *[]crawlResult.PageData,
	logger *zap.SugaredLogger,
	redisWrapper *rediswrapper.RedisWrapper,
) {
	linkStats := stats.NewStats()

	// Pass crawlSiteID to handleHTMLParsing along with other parameters
	err := handleHTMLParsing(ctx, crawlSiteID, logger, collector, searchTerms, linkStats, results, redisWrapper)
	if err != nil {
		logger.Error("Error during HTML parsing", zap.Error(err))
		return
	}

	handleErrorEvents(collector, logger)

	// Assuming handleRedisOperations does not need CrawlSiteID; if it does, it needs to be added as a parameter
	collector.OnScraped(func(r *colly.Response) {
		err := handleRedisOperations(ctx, redisWrapper, logger)
		if err != nil {
			logger.Error("Error with Redis operations", zap.Error(err))
			return
		}

		logger.Info("Finished scraping the page", zap.String("url", r.Request.URL.String()))
		logger.Info("Total links found", zap.Int("total_links", linkStats.TotalLinks))
		logger.Info("Matched links", zap.Int("matched_links", linkStats.MatchedLinks))
		logger.Info("Not matched links", zap.Int("not_matched_links", linkStats.NotMatchedLinks))
		// Here, you would add code to populate the 'results' slice with data
	})

	collector.OnRequest(func(r *colly.Request) {
		logger.Info("Visiting", zap.String("url", r.URL.String()))
	})
}

func getHostFromURL(inputURL string) string {
	u, err := url.Parse(inputURL)
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
	}
	return u.Host
}
