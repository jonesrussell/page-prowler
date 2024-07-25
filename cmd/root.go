package cmd

import (
	"context"
	"errors"

	"github.com/jonesrussell/loggo"

	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
)

var ErrCrawlManagerNotInitialized = errors.New("CrawlManager is not initialized")
var ErrSiteidRequired = errors.New("siteid is required")

// NewRootCmd now returns *cobra.Command
func NewRootCmd(manager *crawler.CrawlManager) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "page-prowler",
		Short: "A tool for finding matchlinks from websites",
		Long: `Page Prowler is a tool that finds matchlinks from websites where the URL matches provided terms. It provides functionalities for:

1. Crawling specific websites and extracting matchlinks that match the provided terms ('matchlinks' command)

	In addition to the command line interface, Page Prowler also provides an HTTP API for interacting with the tool.`,
		SilenceErrors: false,
	}

	// Set the manager to the context
	ctx := context.WithValue(context.Background(), common.CrawlManagerKey, manager)

	// Set the context of the command
	rootCmd.SetContext(ctx)

	// Create a new crawl command with the manager
	crawlCmd := NewCrawlCmd()

	// Add the crawl command to the root command
	rootCmd.AddCommand(crawlCmd)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(rootCmd *cobra.Command, logger loggo.LoggerInterface) {
	err := rootCmd.Execute()
	if err != nil {
		logger.Error("root command execute failed", err)
	}
}
