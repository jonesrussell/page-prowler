package cmd

import (
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jonesrussell/page-prowler/internal/tasks"
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

		// Retrieve the Redis connection details
		redisHost := viper.GetString("REDIS_HOST")
		redisPort := viper.GetString("REDIS_PORT")
		redisAuth := viper.GetString("REDIS_AUTH")

		// Create a new asynq.Client using the same Redis connection details
		client := asynq.NewClient(asynq.RedisClientOpt{
			Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
			Password: redisAuth,
		})

		payload := &tasks.CrawlTaskPayload{
			URL:         URL,
			SearchTerms: SearchTerms,
			CrawlSiteID: Crawlsiteid,
			MaxDepth:    MaxDepth,
			Debug:       Debug,
		}

		err := enqueueCrawlTask(client, payload)
		if err != nil {
			log.Fatalf("Error enqueuing crawl task: %v", err)
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
