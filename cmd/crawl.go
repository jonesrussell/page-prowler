package cmd

import (
	"errors"
	"fmt"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// NewCrawlCmd creates a new crawl command
func NewCrawlCmd(manager crawler.CrawlManagerInterface) *cobra.Command {
	crawlCmd := &cobra.Command{
		Use:   "crawl",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCrawlCmd(cmd, args, manager)
		},
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
	cmd *cobra.Command,
	_ []string,
	manager crawler.CrawlManagerInterface, // remove the pointer here
) error {
	// Check if manager is nil
	if manager == nil {
		fmt.Println("Error: manager is nil")
		return errors.New("manager is nil")
	}

	logger := manager.Logger()
	// Check if Logger is nil
	if logger == nil {
		fmt.Println("Error: Logger is nil")
		return errors.New("Logger is nil")
	}

	options, err := getCrawlOptions()
	if err != nil {
		logger.Error("Error getting options", err)
		return err
	}

	// Print options if Debug is enabled
	if options.Debug {
		manager.Logger().Info("CrawlOptions:")
		manager.Logger().Info(fmt.Sprintf("  CrawlSiteID: %s", options.CrawlSiteID))
		manager.Logger().Info(fmt.Sprintf("  Debug: %t", options.Debug))
		manager.Logger().Info(fmt.Sprintf("  DelayBetweenRequests: %s", options.DelayBetweenRequests.String()))
		manager.Logger().Info(fmt.Sprintf("  MaxConcurrentRequests: %d", options.MaxConcurrentRequests))
		manager.Logger().Info(fmt.Sprintf("  MaxDepth: %d", options.MaxDepth))
		manager.Logger().Info(fmt.Sprintf("  SearchTerms: %v", options.SearchTerms))
		manager.Logger().Info(fmt.Sprintf("  StartURL: %s", options.StartURL))
	}

	// Call SetOptions to update the manager's options
	err = manager.SetOptions(options)
	if err != nil {
		manager.Logger().Error("Error setting options", err)
		return err
	}

	manager.Logger().Info("Starting crawling")
	// Now you can use options in your crawl operation
	err = manager.Crawl()
	if err != nil {
		manager.Logger().Error("Error starting crawling", err)
		return err
	}

	if options.Debug {
		manager.Logger().Info("\nFlags:")
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			manager.Logger().Info(fmt.Sprintf(" %-12s : %s", flag.Name, flag.Value.String()))
		})

		manager.Logger().Info("\nRedis Environment Variables:")
		manager.Logger().Info(fmt.Sprintf(" %-12s : %s", "REDIS_HOST", viper.GetString("REDIS_HOST")))
		manager.Logger().Info(fmt.Sprintf(" %-12s : %s", "REDIS_PORT", viper.GetString("REDIS_PORT")))
		manager.Logger().Info(fmt.Sprintf(" %-12s : %s", "REDIS_AUTH", viper.GetString("REDIS_AUTH")))
	}

	return nil
}

func getCrawlOptions() (*crawler.CrawlOptions, error) {
	// Create an instance of CrawlOptions
	options := &crawler.CrawlOptions{}

	// Populate CrawlOptions fields
	options.CrawlSiteID = viper.GetString("siteid")
	options.Debug = viper.GetBool("debug")
	options.DelayBetweenRequests = viper.GetDuration("delaybetweenrequests")
	options.MaxConcurrentRequests = viper.GetInt("maxconcurrentrequests")
	options.MaxDepth = viper.GetInt("maxdepth")
	options.SearchTerms = viper.GetStringSlice("searchterms")
	options.StartURL = viper.GetString("url")

	return options, nil
}
