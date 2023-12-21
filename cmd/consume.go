package cmd

import (
	"context"
	"fmt"
	"os"

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

		startConsuming(ctx, crawlSiteID, debug)
	},
}

func init() {
	rootCmd.AddCommand(consumeCmd)
}

func startConsuming(ctx context.Context, crawlSiteID string, debug bool) {
	crawlerService, err := initializeManager(ctx, debug)
	if err != nil {
		crawlerService.Logger.Error("Failed to initialize Consume Manager", "error", err)
		os.Exit(1)
	}

	urls, err := crawlerService.RedisClient.SMembers(ctx, crawlSiteID) // Use RedisClient instead of RedisWrapper
	if err != nil {
		crawlerService.Logger.Error("Error fetching URLs from Redis", "error", err)
		return
	}

	for _, url := range urls {
		crawlerService.Logger.Info("Fetched URL", "url", url)
	}
}
