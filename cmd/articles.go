package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	MaxDepth    int
	SearchTerms string
	URL         string
)

var articlesCmd = &cobra.Command{
	Use:   "articles",
	Short: "Crawl websites and extract information",
	Long: `Crawl is a CLI tool designed to perform web scraping and data extraction from websites.
           It allows users to specify parameters such as depth of crawl and target elements to extract.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if Crawlsiteid == "" {
			return fmt.Errorf("crawlsiteid is required")
		}

		if SearchTerms == "" {
			return fmt.Errorf("searchterms is required")
		}

		if URL == "" {
			return fmt.Errorf("url is required")
		}

		if Debug {
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
		manager, ok := cmd.Context().Value(managerKey).(*crawler.CrawlManager)
		if !ok || manager == nil {
			log.Fatalf("CrawlManager is not initialized")
		}

		myServerInstance := &CrawlServer{
			CrawlManager: manager,
		}

		if err := StartCrawling(ctx, URL, SearchTerms, Crawlsiteid, MaxDepth, Debug, manager, myServerInstance); err != nil {
			log.Fatalf("Error starting crawling: %v", err)
		}

		return nil
	},
}

func init() {
	articlesCmd.Flags().StringVarP(&Crawlsiteid, "crawlsiteid", "s", "", "CrawlSite ID")
	if err := viper.BindPFlag("crawlsiteid", articlesCmd.Flags().Lookup("crawlsiteid")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	articlesCmd.Flags().StringVarP(&URL, "url", "u", "", "URL to crawl")
	if err := viper.BindPFlag("url", articlesCmd.Flags().Lookup("url")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	articlesCmd.Flags().StringVarP(&SearchTerms, "searchterms", "t", "", "Search terms for crawling")
	if err := viper.BindPFlag("searchterms", articlesCmd.Flags().Lookup("searchterms")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	articlesCmd.Flags().IntVarP(&MaxDepth, "maxdepth", "m", 1, "Max depth for crawling")
	if err := viper.BindPFlag("maxdepth", articlesCmd.Flags().Lookup("maxdepth")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	rootCmd.AddCommand(articlesCmd)
}

func (s *CrawlServer) saveResultsToRedis(ctx context.Context, results []crawler.PageData, key string) error {
	for _, result := range results {
		data, err := result.MarshalBinary()
		if err != nil {
			s.CrawlManager.Logger.Error("Error occurred during marshalling to binary", "error", err)
			return err
		}
		str := string(data)
		count, err := s.CrawlManager.Client.SAdd(ctx, key, str)
		if err != nil {
			s.CrawlManager.Logger.Error("Error occurred during saving to Redis", "error", err)
			return err
		}
		fmt.Println("Added", count, "elements to the set")
	}
	return nil
}

func printResults(crawlerService *crawler.CrawlManager, results []crawler.PageData) {
	jsonData, err := json.Marshal(results)
	if err != nil {
		crawlerService.Logger.Error("Error occurred during marshaling", "error", err)
		return
	}

	fmt.Println(string(jsonData))
}
