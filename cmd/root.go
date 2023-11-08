// Package cmd contains the command-line commands for the crawler application.
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jonesrussell/crawler/internal/crawler"
	"github.com/jonesrussell/crawler/internal/crawlresult"
	"github.com/jonesrussell/crawler/internal/logger"
	"github.com/jonesrussell/crawler/internal/rediswrapper"
	"github.com/jonesrussell/crawler/internal/stats"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var linkStats = stats.NewStats() // Define linkStats

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl websites and extract information",
	Long: `Crawl is a CLI tool designed to perform web scraping and data extraction from websites.
           It allows users to specify parameters such as depth of crawl and target elements to extract.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		startCrawling(ctx, viper.GetString("url"), viper.GetString("searchterms"), viper.GetString("crawlsiteid"), viper.GetInt("maxdepth"), viper.GetBool("debug"))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().String("url", "", "URL to crawl")
	rootCmd.PersistentFlags().String("searchterms", "", "Comma-separated search terms")
	rootCmd.PersistentFlags().String("crawlsiteid", "", "CrawlSite ID")
	rootCmd.PersistentFlags().Int("maxdepth", 1, "Maximum depth for the crawler")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	rootCmd.MarkPersistentFlagRequired("url")
	rootCmd.MarkPersistentFlagRequired("searchterms")
	rootCmd.MarkPersistentFlagRequired("crawlsiteid")

	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("searchterms", rootCmd.PersistentFlags().Lookup("searchterms"))
	viper.BindPFlag("crawlsiteid", rootCmd.PersistentFlags().Lookup("crawlsiteid"))
	viper.BindPFlag("maxdepth", rootCmd.PersistentFlags().Lookup("maxdepth"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

func initConfig() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv() // Automatically override values from the .env file with those from the environment.

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error while reading config file", err)
	}
}

// initializeCrawlManager sets up the necessary services for the crawler.
func initializeCrawlManager(ctx context.Context, debug bool) *crawler.CrawlManager {
	// Fetch Redis configuration from Viper
	redisHost := viper.GetString("REDIS_HOST")
	redisPort := viper.GetString("REDIS_PORT")
	redisAuth := viper.GetString("REDIS_AUTH")

	// Initialize Logger with the new logger package
	log := logger.New(debug) // Use the logger package's New function to get a logger instance

	// Initialize RedisWrapper
	redisWrapper, err := rediswrapper.NewRedisWrapper(ctx, redisHost, redisPort, redisAuth)
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
