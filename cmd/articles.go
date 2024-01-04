package cmd

import (
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jonesrussell/page-prowler/internal/crawler"
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

var matchlinksCmd = &cobra.Command{
	Use:   "matchlinks",
	Short: "Crawl websites and extract information",
	Long: `Crawl is a CLI tool designed to perform web scraping and data extraction from websites.
           It allows users to specify parameters such as depth of crawl and target elements to extract.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Access the CrawlManager from the context
		value := cmd.Context().Value(managerKey)
		if value == nil {
			log.Fatalf("managerKey not found in context")
		}
		manager, ok := value.(*crawler.CrawlManager)
		if !ok {
			log.Fatalf("managerKey in context is not of type *crawler.CrawlManager")
		}
		if manager == nil {
			log.Fatalf("manager is nil")
		}

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
			manager.Logger.Info("\nFlags:")
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				manager.Logger.Infof(" %-12s : %v\n", flag.Name, flag.Value)
			})

			manager.Logger.Info("\nRedis Environment Variables:")
			manager.Logger.Infof(" %-12s : %s\n", "REDIS_HOST", viper.GetString("REDIS_HOST"))
			manager.Logger.Infof(" %-12s : %s\n", "REDIS_PORT", viper.GetString("REDIS_PORT"))
			manager.Logger.Infof(" %-12s : %s\n", "REDIS_AUTH", viper.GetString("REDIS_AUTH"))
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

		err := tasks.EnqueueCrawlTask(client, payload)
		if err != nil {
			log.Fatalf("Error enqueuing crawl task: %v", err)
		}

		return nil
	},
}

func init() {
	matchlinksCmd.Flags().StringVarP(&Crawlsiteid, "crawlsiteid", "s", "", "CrawlSite ID")
	if err := viper.BindPFlag("crawlsiteid", matchlinksCmd.Flags().Lookup("crawlsiteid")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	matchlinksCmd.Flags().StringVarP(&URL, "url", "u", "", "URL to crawl")
	if err := viper.BindPFlag("url", matchlinksCmd.Flags().Lookup("url")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	matchlinksCmd.Flags().StringVarP(&SearchTerms, "searchterms", "t", "", "Search terms for crawling")
	if err := viper.BindPFlag("searchterms", matchlinksCmd.Flags().Lookup("searchterms")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	matchlinksCmd.Flags().IntVarP(&MaxDepth, "maxdepth", "m", 1, "Max depth for crawling")
	if err := viper.BindPFlag("maxdepth", matchlinksCmd.Flags().Lookup("maxdepth")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	rootCmd.AddCommand(matchlinksCmd)
}
