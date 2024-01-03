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

		links, err := redisClient.SMembers(cmd.Context(), Crawlsiteid)
		if err != nil {
			return fmt.Errorf("failed to get links from Redis: %v", err)
		}
		if err != nil {
			return fmt.Errorf("failed to get links from Redis: %v", err)
		}

		if len(links) == 0 {
			manager.Logger.Info("No links found for the provided crawlsiteid")
		} else {
			for _, link := range links {
				manager.Logger.Info(link)
			}
		}

		return nil
	},
}

func init() {
	getLinksCmd.Flags().StringVarP(&Crawlsiteid, "crawlsiteid", "s", "", "CrawlSite ID")
	rootCmd.AddCommand(getLinksCmd)
}
