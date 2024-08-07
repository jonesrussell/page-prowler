package cmd

import (
	"errors"

	"github.com/jonesrussell/loggo"

	"github.com/jonesrussell/page-prowler/crawler"
	"github.com/spf13/cobra"
)

var ErrCrawlManagerNotInitialized = errors.New("CrawlManager is not initialized")
var ErrSiteidRequired = errors.New("siteid is required")

var debug bool

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

	// Add a debug flag
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")

	// Create a new crawl command with the manager
	crawlCmd := NewCrawlCmd(manager)
	resultsCmd := NewResultsCmd(manager)
	// apiCmd := NewAPICmd(manager)
	workerCmd := NewWorkerCmd(manager)
	getLinksCmd := NewGetLinksCmd(manager)
	clearlinksCmd := NewClearlinksCmd(manager)

	// Add the crawl command to the root command
	rootCmd.AddCommand(crawlCmd)
	rootCmd.AddCommand(resultsCmd)
	// rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(workerCmd)
	rootCmd.AddCommand(getLinksCmd)
	rootCmd.AddCommand(clearlinksCmd)

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
