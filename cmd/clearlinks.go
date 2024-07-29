package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewClearlinksCmd creates a new clearlinks command
func NewClearlinksCmd(manager crawler.CrawlManagerInterface) *cobra.Command {
	clearlinksCmd := &cobra.Command{
		Use:   "clearlinks",
		Short: "Clear the Redis set for a given siteid",
		RunE: func(_ *cobra.Command, _ []string) error {
			return ClearlinksMain(manager)
		},
	}

	return clearlinksCmd
}

func ClearlinksMain(manager crawler.CrawlManagerInterface) error {
	siteid := viper.GetString("siteid")
	if siteid == "" {
		return ErrSiteidRequired
	}

	// Check if manager is nil
	if manager == nil {
		fmt.Println("Error: manager is nil")
		return errors.New("manager is nil")
	}

	dbManager := manager.GetDBManager()

	err := dbManager.ClearRedisSet(context.Background(), siteid)
	if err != nil {
		return fmt.Errorf("failed to clear Redis set: %v", err)
	}

	debug := viper.GetBool("debug")
	if debug {
		manager.Logger().Debug("Debugging enabled. Clearing Redis set...")
	}

	manager.Logger().Info("Redis set cleared successfully")

	return nil
}
