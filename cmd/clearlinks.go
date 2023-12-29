package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clearlinksCmd = &cobra.Command{
	Use:   "clearlinks",
	Short: "Clear the Redis set for a given crawlsiteid",
	Run: func(cmd *cobra.Command, args []string) {
		crawlsiteid := viper.GetString("crawlsiteid")
		if crawlsiteid == "" {
			fmt.Println("crawlsiteid is required")
			os.Exit(1)
		}

		// Get the manager from the context
		manager, ok := cmd.Context().Value(managerKey).(*crawler.CrawlManager)
		if !ok || manager == nil {
			log.Fatalf("CrawlManager is not initialized")
		}

		// Use the Redis client from the manager
		redisClient := manager.Client

		_, err := redisClient.Del(cmd.Context(), crawlsiteid)
		if err != nil {
			fmt.Println("Failed to clear Redis set", "error", err)
			os.Exit(1)
		}

		fmt.Println("Redis set cleared successfully")
	},
}

func init() {
	rootCmd.AddCommand(clearlinksCmd)
}
