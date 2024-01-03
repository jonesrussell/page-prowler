package cmd

import (
	"fmt"
	"log"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clearlinksCmd = &cobra.Command{
	Use:   "clearlinks",
	Short: "Clear the Redis set for a given crawlsiteid",
	RunE: func(cmd *cobra.Command, args []string) error {
		if Crawlsiteid == "" {
			return fmt.Errorf("crawlsiteid is required")
		}

		manager, ok := cmd.Context().Value(managerKey).(*crawler.CrawlManager)
		if !ok || manager == nil {
			return fmt.Errorf("CrawlManager is not initialized")
		}

		redisClient := manager.Client

		err := redisClient.Del(cmd.Context(), Crawlsiteid)
		if err != nil {
			return fmt.Errorf("failed to clear Redis set: %v", err)
		}

		manager.Logger.Info("Redis set cleared successfully")

		return nil
	},
}

func init() {
	clearlinksCmd.Flags().StringVarP(&Crawlsiteid, "crawlsiteid", "s", "", "CrawlSite ID")
	if err := viper.BindPFlag("crawlsiteid", clearlinksCmd.Flags().Lookup("crawlsiteid")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}
	rootCmd.AddCommand(clearlinksCmd)
}
