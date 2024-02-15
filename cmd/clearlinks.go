package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
)

var ClearlinksCmd = &cobra.Command{
	Use:   "clearlinks",
	Short: "Clear the Redis set for a given siteid",
	RunE: func(cmd *cobra.Command, _ []string) error {
		log.Println("RunE function started")

		siteid, _ := cmd.Flags().GetString("siteid")
		if siteid == "" {
			return errors.New("siteid is required")
		}

		manager, ok := cmd.Context().Value(common.CrawlManagerKey).(*crawler.CrawlManager)
		if !ok || manager == nil {
			return ErrCrawlManagerNotInitialized
		}

		redisClient := manager.Client

		err := redisClient.Del(cmd.Context(), siteid)
		if err != nil {
			return fmt.Errorf("failed to clear Redis set: %v", err)
		}

		if Debug {
			manager.LoggerField.Debug("Debugging enabled. Clearing Redis set...")
		}

		manager.Logger().Info("Redis set cleared successfully")

		log.Println("RunE function ended")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ClearlinksCmd)
}
