package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/crawlresult"
	"github.com/jonesrussell/page-prowler/internal/stats"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl websites and extract information",
	Long: `Crawl is a CLI tool designed to perform web scraping and data extraction from websites.
           It allows users to specify parameters such as depth of crawl and target elements to extract.`,
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()

		if viper.GetBool("debug") {
			fmt.Println("\nFlags:")
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				fmt.Printf("  %-12s : %v\n", flag.Name, flag.Value)
			})

			fmt.Println("\nRedis Environment Variables:")
			fmt.Printf("  %-12s : %s\n", "REDIS_HOST", viper.GetString("REDIS_HOST"))
			fmt.Printf("  %-12s : %s\n", "REDIS_PORT", viper.GetString("REDIS_PORT"))
			fmt.Printf("  %-12s : %s\n", "REDIS_AUTH", viper.GetString("REDIS_AUTH"))
		}

		ctx := context.Background()
		crawlerService, err := initializeManager(ctx, viper.GetBool("debug"))
		if err != nil {
			fmt.Println("Failed to initialize Crawl Manager", "error", err)
			os.Exit(1)
		}

		startCrawling(ctx, viper.GetString("url"), viper.GetString("searchterms"), viper.GetString("crawlsiteid"), viper.GetInt("maxdepth"), viper.GetBool("debug"), crawlerService)
	},
}

func init() {
	crawlCmd.Flags().String("url", "", "URL to crawl")
	crawlCmd.Flags().String("searchterms", "", "Search terms for crawling")
	crawlCmd.Flags().Int("maxdepth", 1, "Maximum depth for crawling")

	viper.BindPFlag("url", crawlCmd.Flags().Lookup("url"))
	viper.BindPFlag("searchterms", crawlCmd.Flags().Lookup("searchterms"))
	viper.BindPFlag("maxdepth", crawlCmd.Flags().Lookup("maxdepth"))

	rootCmd.AddCommand(crawlCmd)
}

func startCrawling(ctx context.Context, url, searchTerms, crawlSiteID string, maxDepth int, debug bool, crawlerService *crawler.CrawlManager) {
	splitSearchTerms := strings.Split(searchTerms, ",")
	host, err := crawler.GetHostFromURL(url, crawlerService.Logger)
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

	options := crawler.CrawlOptions{
		CrawlSiteID: crawlSiteID,
		Collector:   collector,
		SearchTerms: splitSearchTerms,
		Results:     &results,
		LinkStats:   stats.NewStats(),
		Debug:       debug,
	}
	crawlerService.SetupCrawlingLogic(ctx, options)

	crawlerService.Logger.Info("Crawler started...")
	if err := collector.Visit(url); err != nil {
		crawlerService.Logger.Error("Error visiting URL", "url", url, "error", err)
		return
	}

	collector.Wait()

	crawlerService.Logger.Info("Crawling completed.")

	saveResultsToRedis(ctx, crawlerService, results)
	printResults(crawlerService, results)
}

func saveResultsToRedis(ctx context.Context, crawlerService *crawler.CrawlManager, results []crawlresult.PageData) {
	for _, result := range results {
		_, err := crawlerService.RedisWrapper.SAdd(ctx, "yourKeyHere", result)
		if err != nil {
			crawlerService.Logger.Error("Error occurred during saving to Redis", "error", err)
			return
		}
	}
}

func printResults(crawlerService *crawler.CrawlManager, results []crawlresult.PageData) {
	jsonData, err := json.Marshal(results)
	if err != nil {
		crawlerService.Logger.Error("Error occurred during marshaling", "error", err)
		return
	}

	fmt.Println(string(jsonData))
}
