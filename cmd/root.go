/*
Copyright © 2023 Russell Jones jonesrussell42@gmail.com
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jonesrussell/crawler/internal/crawlResult"
	"github.com/jonesrussell/crawler/internal/crawler"
	"github.com/jonesrussell/crawler/internal/rediswrapper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper" // Import viper
	"go.uber.org/zap"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl websites and extract information",
	Long: `Crawl is a CLI tool designed to perform web scraping and data extraction from websites.
           It allows users to specify parameters such as depth of crawl and target elements to extract.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Your application's logic goes here, replacing the main() content
		ctx := context.Background()
		// Use the flags directly, since they are now package-level variables
		// args represents additional arguments passed to the command
		startCrawling(ctx, url, searchTerms, crawlSiteID, maxDepth, debug)
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

var url string
var searchTerms string
var crawlSiteID string
var maxDepth int
var debug bool

func init() {
	cobra.OnInitialize(initConfig) // Initialize viper when the application starts
	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "URL to crawl")
	rootCmd.PersistentFlags().StringVarP(&searchTerms, "searchterms", "s", "", "Comma-separated search terms")
	rootCmd.PersistentFlags().StringVarP(&crawlSiteID, "crawlsiteid", "c", "", "CrawlSite ID")
	rootCmd.PersistentFlags().IntVarP(&maxDepth, "maxdepth", "m", 1, "Maximum depth for the crawler")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")

	rootCmd.MarkPersistentFlagRequired("url")
	rootCmd.MarkPersistentFlagRequired("searchterms")
	rootCmd.MarkPersistentFlagRequired("crawlsiteid")

	// Bind the flags to Viper parameters
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("searchterms", rootCmd.PersistentFlags().Lookup("searchterms"))
	viper.BindPFlag("crawlsiteid", rootCmd.PersistentFlags().Lookup("crawlsiteid"))
	viper.BindPFlag("maxdepth", rootCmd.PersistentFlags().Lookup("maxdepth"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

func initConfig() {
	viper.SetConfigFile(".env") // name of your env file
	viper.SetConfigType("env")  // config file type to be .env

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error while reading config file", err)
	}
}

func startCrawling(ctx context.Context, url, searchTerms, crawlSiteID string, maxDepth int, debug bool) {
	// Initialize Logger
	logger, err := crawler.InitializeLogger(debug)
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		os.Exit(1)
	}

	// Load configuration
	envCfg, err := crawler.LoadConfiguration()
	if err != nil {
		logger.Errorf("Failed to load environment configuration: %v", err)
		os.Exit(1)
	}

	// Initialize Redis with the context
	redisWrapper, err := rediswrapper.NewRedisWrapper(ctx, envCfg.RedisHost, envCfg.RedisPort, envCfg.RedisAuth, logger)
	if err != nil {
		logger.Errorf("Failed to initialize Redis: %v", err)
		os.Exit(1)
	}

	// Split search terms
	splitSearchTerms := strings.Split(searchTerms, ",")

	// Configure the collector
	collector := crawler.ConfigureCollector([]string{crawler.GetHostFromURL(url)}, maxDepth)
	if collector == nil {
		logger.Fatal("Failed to configure collector")
		return
	}

	// Setup crawling logic
	var results []crawlResult.PageData
	crawler.SetupCrawlingLogic(ctx, crawlSiteID, collector, splitSearchTerms, &results, logger, redisWrapper)

	// Start the crawling process
	logger.Info("Crawler started...")
	if err := collector.Visit(url); err != nil {
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
