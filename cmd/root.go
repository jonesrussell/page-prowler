package cmd

import (
	"errors"

	"github.com/jonesrussell/page-prowler/news"

	"github.com/jonesrussell/page-prowler/crawler"
	"github.com/spf13/cobra"
)

var ErrCrawlManagerNotInitialized = errors.New("CrawlManager is not initialized")
var ErrSiteidRequired = errors.New("siteid is required")

var debug bool

// NewRootCmd now returns *cobra.Command and accepts newsService
func NewRootCmd(manager *crawler.CrawlManager, newsService news.Service) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "page-prowler",
		Short: "A tool for finding matchlinks from websites",
		Long: `Page Prowler is a tool that finds matchlinks from websites where the URL matches provided terms. It provides functionalities for:

1. Crawling specific websites and extracting matchlinks that match the provided terms ('matchlinks' command)
2. Generating static news sites ('gensite' command)

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
	genSiteCmd := NewGenSiteCmd(newsService) // Pass newsService to NewGenSiteCmd

	serveCmd := NewServeCmd(newsService)

	// Add the commands to the root command
	rootCmd.AddCommand(crawlCmd)
	rootCmd.AddCommand(resultsCmd)
	// rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(workerCmd)
	rootCmd.AddCommand(getLinksCmd)
	rootCmd.AddCommand(clearlinksCmd)
	rootCmd.AddCommand(genSiteCmd)
	rootCmd.AddCommand(serveCmd)

	return rootCmd
}
