package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/jonesrussell/page-prowler/internal/common"
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
			return ErrCrawlsiteidRequired
		}

		ctx := context.Background()
		debug := viper.GetBool("debug")

		// Get the manager from the context
		manager, ok := cmd.Context().Value(common.ManagerKey).(*crawler.CrawlManager)
		if !ok || manager == nil {
			log.Fatalf("CrawlManager is not initialized")
		}

		if debug {
			manager.Logger.Info("All configuration keys and values:")
			for _, key := range viper.AllKeys() {
				manager.Logger.Info(fmt.Sprintf("%s: %v\n", key, viper.Get(key)))
			}
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
	urls, err := manager.Client.SMembers(ctx, crawlSiteID)
	if err != nil {
		manager.Logger.Error("Error fetching URLs from Redis", "error", err)
		return
	}

	for _, url := range urls {
		manager.Logger.Info("Fetched URL", "url", url)
	}
}
