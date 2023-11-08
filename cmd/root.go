// Package cmd contains the command-line commands for the crawler application.
package cmd

import (
	"os"

	"github.com/jonesrussell/page-prowler/cmd/consume"
	"github.com/jonesrussell/page-prowler/cmd/crawl"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "page-prowler",
	Short: "Page Prowler is a CLI tool for web crawling and data extraction",
	Long: `Page Prowler is a CLI tool designed to perform web scraping and data extraction from websites.
           It allows users to specify parameters such as depth of crawl and target elements to extract.`,
}

func init() {
	rootCmd.AddCommand(consume.ConsumeCmd)
	rootCmd.AddCommand(crawl.CrawlCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
