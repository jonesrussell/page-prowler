package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "Consume URLs from Redis",
	Long:  `Consume is a CLI tool designed to fetch URLs from a Redis set.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		crawlSiteID := viper.GetString("crawlsiteid")
		debug := viper.GetBool("debug")

		if crawlSiteID == "" {
			fmt.Println("CrawlSiteId is required")
			os.Exit(1)
		}

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

		startConsuming(ctx, crawlSiteID, debug, manager)
	},
}

func init() {
	rootCmd.AddCommand(consumeCmd)
}

func startConsuming(ctx context.Context, crawlSiteID string, debug bool, manager *crawler.CrawlManager) {
	urls, err := manager.Client.SMembers(ctx, crawlSiteID) // Use Client instead of RedisWrapper
	if err != nil {
		manager.Logger.Error("Error fetching URLs from Redis", "error", err)
		return
	}

	for _, url := range urls {
		manager.Logger.Info("Fetched URL", "url", url)
	}
}
