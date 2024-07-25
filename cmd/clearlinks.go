package cmd

import (
	"fmt"

	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewClearlinksCmd creates a new clearlinks command
func NewClearlinksCmd() *cobra.Command {
	clearlinksCmd := &cobra.Command{
		Use:   "clearlinks",
		Short: "Clear the Redis set for a given siteid",
		RunE:  ClearlinksMain,
	}

	return clearlinksCmd
}

func ClearlinksMain(cmd *cobra.Command, _ []string) error {
	siteid := viper.GetString("siteid")
	if siteid == "" {
		return ErrSiteidRequired
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

	debug := viper.GetBool("debug")
	if debug {
		manager.LoggerField.Debug("Debugging enabled. Clearing Redis set...")
	}

	manager.Logger().Info("Redis set cleared successfully")

	return nil
}
