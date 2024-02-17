package cmd

import (
	"fmt"

	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
)

var ClearlinksCmd = &cobra.Command{
	Use:   "clearlinks",
	Short: "Clear the Redis set for a given siteid",
	RunE:  ClearlinksMain,
}

func ClearlinksMain(cmd *cobra.Command, _ []string) error {
	if Siteid == "" {
		return ErrSiteidRequired
	}

	manager, ok := cmd.Context().Value(common.CrawlManagerKey).(*crawler.CrawlManager)
	if !ok || manager == nil {
		return ErrCrawlManagerNotInitialized
	}

	redisClient := manager.Client

	err := redisClient.Del(cmd.Context(), Siteid)
	if err != nil {
		return fmt.Errorf("failed to clear Redis set: %v", err)
	}

	if Debug {
		manager.LoggerField.Debug("Debugging enabled. Clearing Redis set...", map[string]interface{}{})
	}

	manager.Logger().Info("Redis set cleared successfully", map[string]interface{}{})

	return nil
}

func init() {
	resultsCmd.AddCommand(ClearlinksCmd)
}
