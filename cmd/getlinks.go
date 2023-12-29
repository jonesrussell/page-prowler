package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getLinksCmd = &cobra.Command{
	Use:   "getlinks",
	Short: "Get the list of links for a given crawlsiteid",
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

		// Print the manager
		fmt.Printf("Manager: %+v\n", manager)

		// Use the Redis client from the manager
		redisClient := manager.Client

		links, err := redisClient.SMembers(cmd.Context(), crawlsiteid)
		if err != nil {
			fmt.Println("Failed to get links from Redis", "error", err)
			os.Exit(1)
		}

		for _, link := range links {
			fmt.Println(link)
		}
	},
}

func init() {
	rootCmd.AddCommand(getLinksCmd)
}
