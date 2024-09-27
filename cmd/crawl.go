package cmd

import (
	"errors"
	"fmt"

	"github.com/jonesrussell/page-prowler/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCrawlCmd creates a new crawl command
func NewCrawlCmd(manager crawler.CrawlManagerInterface) *cobra.Command {
	crawlCmd := &cobra.Command{
		Use:   "crawl",
		Short: "Crawl!",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runCrawlCmd(manager)
		},
	}

	crawlCmd.Flags().StringP("siteid", "s", "", "Site ID for crawling")
	if err := viper.BindPFlag("siteid", crawlCmd.Flags().Lookup("siteid")); err != nil {
		fmt.Println("Error binding flag", err)
	}

	crawlCmd.Flags().StringP("url", "u", "", "URL to crawl")
	if err := viper.BindPFlag("url", crawlCmd.Flags().Lookup("url")); err != nil {
		fmt.Println("Error binding flag", err)
	}

	crawlCmd.Flags().IntP("maxdepth", "m", 1, "Max depth for crawling")
	if err := viper.BindPFlag("maxdepth", crawlCmd.Flags().Lookup("maxdepth")); err != nil {
		fmt.Println("Error binding flag", err)
	}

	crawlCmd.Flags().StringP("searchterms", "t", "", "Search terms for crawling")
	if err := viper.BindPFlag("searchterms", crawlCmd.Flags().Lookup("searchterms")); err != nil {
		fmt.Println("Error binding flag", err)
	}

	return crawlCmd
}

func runCrawlCmd(
	manager crawler.CrawlManagerInterface,
) error {
	// Check if manager is nil
	if manager == nil {
		fmt.Println("Error: manager is nil")
		return errors.New("manager is nil")
	}

	logger := manager.GetLogger()
	// Check if Logger is nil
	if logger == nil {
		fmt.Println("Error: Logger is nil")
		return errors.New("logger is nil")
	}

	options, err := getCrawlOptions()
	if err != nil {
		logger.Error("Error getting options", err)
		return err
	}

	// Print options if Debug is enabled
	if options.Debug {
		logger.Info("CrawlOptions:")
		logger.Info(fmt.Sprintf("  CrawlSiteID: %s", options.CrawlSiteID))
		logger.Info(fmt.Sprintf("  Debug: %t", options.Debug))
		logger.Info(fmt.Sprintf("  DelayBetweenRequests: %s", options.DelayBetweenRequests.String()))
		logger.Info(fmt.Sprintf("  MaxConcurrentRequests: %d", options.MaxConcurrentRequests))
		logger.Info(fmt.Sprintf("  MaxDepth: %d", options.MaxDepth))
		logger.Info(fmt.Sprintf("  SearchTerms: %v", options.SearchTerms))
		logger.Info(fmt.Sprintf("  StartURL: %s", options.StartURL))
	}

	// Call SetOptions to update the manager's options
	err = manager.SetOptions(options)
	if err != nil {
		logger.Error("Error setting options", err)
		return err
	}

	logger.Info("Starting crawling")

	err = manager.Crawl()
	if err != nil {
		logger.Error("Error starting crawling", err)
		return err
	}

	return nil
}

func getCrawlOptions() (*crawler.CrawlOptions, error) {
	// Create an instance of CrawlOptions
	options := &crawler.CrawlOptions{}

	// Populate CrawlOptions fields
	options.CrawlSiteID = viper.GetString("siteid")
	options.Debug = debug
	options.DelayBetweenRequests = viper.GetDuration("delaybetweenrequests")
	options.MaxConcurrentRequests = viper.GetInt("maxconcurrentrequests")
	options.MaxDepth = viper.GetInt("maxdepth")
	options.SearchTerms = viper.GetStringSlice("searchterms")
	options.StartURL = viper.GetString("url")

	return options, nil
}
