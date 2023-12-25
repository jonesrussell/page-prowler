package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/crawlresult"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var articlesCmd = &cobra.Command{
	Use:   "articles",
	Short: "Crawl websites and extract information",
	Long: `Crawl is a CLI tool designed to perform web scraping and data extraction from websites.
           It allows users to specify parameters such as depth of crawl and target elements to extract.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		myServerInstance := &MyServer{
			CrawlManager: crawlerService,
		}

		StartCrawling(ctx, viper.GetString("url"), viper.GetString("searchterms"), viper.GetString("crawlsiteid"), viper.GetInt("maxdepth"), viper.GetBool("debug"), crawlerService, myServerInstance)
	},
}

func init() {
	articlesCmd.Flags().String("url", "", "URL to crawl")
	articlesCmd.Flags().String("searchterms", "", "Search terms for crawling")
	articlesCmd.Flags().Int("maxdepth", 1, "Maximum depth for crawling")

	viper.BindPFlag("url", articlesCmd.Flags().Lookup("url"))
	viper.BindPFlag("searchterms", articlesCmd.Flags().Lookup("searchterms"))
	viper.BindPFlag("maxdepth", articlesCmd.Flags().Lookup("maxdepth"))

	rootCmd.AddCommand(articlesCmd)
}

func (s *MyServer) saveResultsToRedis(ctx context.Context, results []crawlresult.PageData) error {
	for _, result := range results {
		data, err := result.MarshalBinary()
		if err != nil {
			s.CrawlManager.Logger.Error("Error occurred during marshalling to binary", "error", err)
			return err
		}
		str := string(data)
		count, err := s.CrawlManager.RedisClient.SAdd(ctx, "yourKeyHere", str)
		if err != nil {
			s.CrawlManager.Logger.Error("Error occurred during saving to Redis", "error", err)
			return err
		}
		fmt.Println("Added", count, "elements to the set")
	}
	return nil
}

func printResults(crawlerService *crawler.CrawlManager, results []crawlresult.PageData) {
	jsonData, err := json.Marshal(results)
	if err != nil {
		crawlerService.Logger.Error("Error occurred during marshaling", "error", err)
		return
	}

	fmt.Println(string(jsonData))
}
