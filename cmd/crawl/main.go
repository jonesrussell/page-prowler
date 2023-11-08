package crawl

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/crawlresult"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/rediswrapper"
	"github.com/jonesrussell/page-prowler/internal/stats"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var linkStats = stats.NewStats() // Define linkStats

// CrawlCmd represents the crawl command
var CrawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl websites and extract information",
	Long: `Crawl is a CLI tool designed to perform web scraping and data extraction from websites.
           It allows users to specify parameters such as depth of crawl and target elements to extract.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		debug := viper.GetBool("debug")
		fmt.Println("Debug:", debug) // Print the value of debug
		startCrawling(ctx, viper.GetString("url"), viper.GetString("searchterms"), viper.GetString("crawlsiteid"), viper.GetInt("maxdepth"), debug)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	CrawlCmd.PersistentFlags().String("url", "", "URL to crawl")
	CrawlCmd.PersistentFlags().String("searchterms", "", "Comma-separated search terms")
	CrawlCmd.PersistentFlags().String("crawlsiteid", "", "CrawlSite ID")
	CrawlCmd.PersistentFlags().Int("maxdepth", 1, "Maximum depth for the crawler")
	CrawlCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	CrawlCmd.MarkPersistentFlagRequired("url")
	CrawlCmd.MarkPersistentFlagRequired("searchterms")
	CrawlCmd.MarkPersistentFlagRequired("crawlsiteid")

	viper.BindPFlag("url", CrawlCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("searchterms", CrawlCmd.PersistentFlags().Lookup("searchterms"))
	viper.BindPFlag("crawlsiteid", CrawlCmd.PersistentFlags().Lookup("crawlsiteid"))
	viper.BindPFlag("maxdepth", CrawlCmd.PersistentFlags().Lookup("maxdepth"))
	viper.BindPFlag("debug", CrawlCmd.PersistentFlags().Lookup("debug"))
}

func initConfig() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv() // Automatically override values from the .env file with those from the environment.

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error while reading config file", err)
	}

	// Bind the current command's flags to viper
	viper.BindPFlags(CrawlCmd.PersistentFlags())
}

// initializeCrawlManager sets up the necessary services for the crawler.
func initializeCrawlManager(ctx context.Context, debug bool) *crawler.CrawlManager {
	// Fetch Redis configuration from Viper
	redisHost := viper.GetString("REDIS_HOST")
	redisPort := viper.GetString("REDIS_PORT")
	redisAuth := viper.GetString("REDIS_AUTH")

	// Initialize Logger with the new logger package
	log := logger.New(debug) // Use the logger package's New function to get a logger instance

	// Create a new Redis client instance.
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisAuth, // no password set
		DB:       0,         // use default DB
	})

	// Initialize RedisWrapper
	redisWrapper, err := rediswrapper.NewRedisWrapper(ctx, rdb)
	if err != nil {
		log.Error("Failed to initialize Redis", "error", err)
		os.Exit(1)
	}

	// Return the CrawlManager instance
	return &crawler.CrawlManager{
		Logger:       log, // Use the new logger instance here
		RedisWrapper: redisWrapper,
	}
}

func startCrawling(ctx context.Context, url, searchTerms, crawlSiteID string, maxDepth int, debug bool) {
	// Initialize CrawlManager
	crawlerService := initializeCrawlManager(ctx, debug)

	splitSearchTerms := strings.Split(searchTerms, ",")
	host, err := crawler.GetHostFromURL(url, crawlerService.Logger) // Now it needs two arguments
	if err != nil {
		crawlerService.Logger.Error("Failed to parse URL", "url", url, "error", err)
		return
	}

	collector := crawler.ConfigureCollector([]string{host}, maxDepth)
	if collector == nil {
		crawlerService.Logger.Fatal("Failed to configure collector")
		return
	}

	var results []crawlresult.PageData
	// Create a CrawlOptions struct and pass it to SetupCrawlingLogic
	options := crawler.CrawlOptions{
		CrawlSiteID: crawlSiteID,
		Collector:   collector,
		SearchTerms: splitSearchTerms,
		Results:     &results,
		LinkStats:   linkStats,
		Debug:       debug,
	}
	crawlerService.SetupCrawlingLogic(ctx, options)

	crawlerService.Logger.Info("Crawler started...")
	if err := collector.Visit(url); err != nil {
		crawlerService.Logger.Error("Error visiting URL", "url", url, "error", err)
		return
	}

	collector.Wait()

	jsonData, err := json.Marshal(results)
	if err != nil {
		crawlerService.Logger.Error("Error occurred during marshaling", "error", err)
		return
	}

	fmt.Println(string(jsonData))
	crawlerService.Logger.Info("Crawling completed.")
}
