package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// crawlCmd represents the crawl command
var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		return fx.New(
			fx.Provide(
				func(logger logger.Logger, client prowlredis.ClientInterface) *crawler.CrawlManager {
					options := &crawler.CrawlOptions{} // Assuming options are configured here or passed differently
					return crawler.NewCrawlManager(logger, client, options)
				},
			),
			fx.Invoke(func(manager *crawler.CrawlManager, lc fx.Lifecycle) {
				lc.Append(fx.Hook{
					OnStart: func(context.Context) error {
						return runCrawlCmd(cmd, manager)
					},
				})
			}),
		).Start(context.Background())
	},
}

func runCrawlCmd(cmd *cobra.Command, manager crawler.CrawlManagerInterface) error {
	if manager == nil {
		log.Fatalf("manager is nil")
	}

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

	// Print options if Debug is enabled
	if options.Debug {
		fmt.Printf("CrawlOptions:\n")
		fmt.Printf("  CrawlSiteID: %s\n", options.CrawlSiteID)
		fmt.Printf("  Debug: %t\n", options.Debug)
		fmt.Printf("  DelayBetweenRequests: %s\n", options.DelayBetweenRequests.String())
		fmt.Printf("  MaxConcurrentRequests: %d\n", options.MaxConcurrentRequests)
		fmt.Printf("  MaxDepth: %d\n", options.MaxDepth)
		fmt.Printf("  SearchTerms: %v\n", options.SearchTerms)
		fmt.Printf("  StartURL: %s\n", options.StartURL)
	}

	// Call SetOptions to update the manager's options
	err := manager.SetOptions(options)
	if err != nil {
		log.Fatalf("Error setting options: %v", err)
	}

	// Now you can use options in your crawl operation
	err = manager.Crawl()
	if err != nil {
		log.Fatalf("Error starting crawling: %v", err)
	}

	if options.Debug {
		manager.Logger().Info("\nFlags:")
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			manager.Logger().Info(fmt.Sprintf(" %-12s : %s\n", flag.Name, flag.Value.String()))
		})

		manager.Logger().Info("\nRedis Environment Variables:")
		manager.Logger().Info(fmt.Sprintf(" %-12s : %s\n", "REDIS_HOST", viper.GetString("REDIS_HOST")))
		manager.Logger().Info(fmt.Sprintf(" %-12s : %s\n", "REDIS_PORT", viper.GetString("REDIS_PORT")))
		manager.Logger().Info(fmt.Sprintf(" %-12s : %s\n", "REDIS_AUTH", viper.GetString("REDIS_AUTH")))
	}

	return nil
}

// init initializes the matchlinks command.
func init() {
	crawlCmd.Flags().StringP("url", "u", "", "URL to crawl")
	if err := viper.BindPFlag("url", crawlCmd.Flags().Lookup("url")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	crawlCmd.Flags().IntP("maxdepth", "m", 1, "Max depth for crawling")
	if err := viper.BindPFlag("maxdepth", crawlCmd.Flags().Lookup("maxdepth")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	crawlCmd.Flags().StringP("searchterms", "t", "", "Search terms for crawling")
	if err := viper.BindPFlag("searchterms", crawlCmd.Flags().Lookup("searchterms")); err != nil {
		log.Fatalf("Error binding flag: %v", err)
	}

	RootCmd.AddCommand(crawlCmd)
}
