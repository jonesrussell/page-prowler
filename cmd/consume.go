package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "Consume URLs from Redis",
	Long:  `Consume is a CLI tool designed to fetch URLs from a Redis set.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if Crawlsiteid == "" {
			return fmt.Errorf("crawlsiteid is required")
		}

		ctx := context.Background()
		debug := viper.GetBool("debug")

		if debug {
			fmt.Println("All configuration keys and values:")
			for _, key := range viper.AllKeys() {
				fmt.Printf("%s: %v\n", key, viper.Get(key))
			}
		}

		// Get the manager from the context
		manager, ok := cmd.Context().Value(managerKey).(*crawler.CrawlManager)
		if !ok || manager == nil {
			log.Fatalf("CrawlManager is not initialized")
		}

		startConsuming(ctx, Crawlsiteid, manager)

		return nil
	},
}

func init() {
	consumeCmd.Flags().StringVarP(&Crawlsiteid, "crawlsiteid", "s", "", "CrawlSite ID")

	rootCmd.AddCommand(consumeCmd)
}

func startConsuming(ctx context.Context, crawlSiteID string, manager *crawler.CrawlManager) {
	urls, err := manager.Client.SMembers(ctx, crawlSiteID) // Use Client instead of RedisWrapper
	if err != nil {
		manager.Logger.Error("Error fetching URLs from Redis", "error", err)
		return
	}

	for _, url := range urls {
		manager.Logger.Info("Fetched URL", "url", url)
	}
}
