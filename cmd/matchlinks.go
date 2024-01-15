package cmd

import (
	"fmt"
	"log"

	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var matchlinksCmd = &cobra.Command{
	Use:   "matchlinks",
	Short: "Crawl websites and extract information",
	Long: `Crawl is a CLI tool designed to perform web scraping and data extraction from websites.
           It allows users to specify parameters such as depth of crawl and target elements to extract.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// Access the CrawlManager from the context
		value := ctx.Value(common.CrawlManagerKey)
		if value == nil {
			log.Fatalf("common.ManagerKey not found in context")
		}

		manager, ok := value.(*crawler.CrawlManager)
		if !ok {
			log.Fatalf("common.ManagerKey in context is not of type *crawler.CrawlManager")
		}
		if manager == nil {
			log.Fatalf("manager is nil")
		}

		crawlsiteid := viper.GetString("crawlsiteid")
		if crawlsiteid == "" {
			return ErrCrawlsiteidRequired
		}

		searchterms := viper.GetString("searchterms")
		if searchterms == "" {
			return fmt.Errorf("searchterms is required")
		}

		url := viper.GetString("url")
		if url == "" {
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

		err := manager.StartCrawling(ctx, url, searchterms, crawlsiteid, viper.GetInt("maxdepth"), Debug)
		if err != nil {
			log.Fatalf("Error starting crawling: %v", err)
		}

		return nil
	},
}

func init() {
	matchlinksCmd.Flags().StringP("crawlsiteid", "s", "", "CrawlSite ID")
	if err := viper.BindPFlag("crawlsiteid", matchlinksCmd.Flags().Lookup("crawlsiteid")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	matchlinksCmd.Flags().StringP("url", "u", "", "URL to crawl")
	if err := viper.BindPFlag("url", matchlinksCmd.Flags().Lookup("url")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	matchlinksCmd.Flags().StringP("searchterms", "t", "", "Search terms for crawling")
	if err := viper.BindPFlag("searchterms", matchlinksCmd.Flags().Lookup("searchterms")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	matchlinksCmd.Flags().IntP("maxdepth", "m", 1, "Max depth for crawling")
	if err := viper.BindPFlag("maxdepth", matchlinksCmd.Flags().Lookup("maxdepth")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}
}
