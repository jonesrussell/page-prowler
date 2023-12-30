package cmd

import (
	"fmt"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
)

var getLinksCmd = &cobra.Command{
	Use:   "getlinks",
	Short: "Get the list of links for a given crawlsiteid",
	RunE: func(cmd *cobra.Command, args []string) error {
		if Crawlsiteid == "" {
			return fmt.Errorf("crawlsiteid is required")
		}

		// Get the manager from the context
		manager, ok := cmd.Context().Value(managerKey).(*crawler.CrawlManager)
		if !ok || manager == nil {
			return fmt.Errorf("CrawlManager is not initialized")
		}

		// Use the Redis client from the manager
		redisClient := manager.Client

		smembersCmd := redisClient.SMembers(cmd.Context(), Crawlsiteid)
		links, err := smembersCmd.Result()
		if err != nil {
			return fmt.Errorf("Failed to get links from Redis: %v", err)
		}

		if len(links) == 0 {
			fmt.Println("No links found for the provided crawlsiteid")
		} else {
			for _, link := range links {
				fmt.Println(link)
			}
		}
		return nil
	},
}

func init() {
	getLinksCmd.Flags().StringVarP(&Crawlsiteid, "crawlsiteid", "s", "", "CrawlSite ID")
	rootCmd.AddCommand(getLinksCmd)
}
