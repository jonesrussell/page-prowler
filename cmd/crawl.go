package cmd

import (
	"fmt"
	"log"

	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
	RunE: runCrawlCmd,
}

func runCrawlCmd(
	cmd *cobra.Command,
	_ []string,
) error {

	ctx := cmd.Context()

	// Access the CrawlManager from the context
	value := ctx.Value(common.CrawlManagerKey)
	if value == nil {
		log.Fatalf("common.CrawlManagerKey not found in context")
	}

	manager, ok := value.(crawler.CrawlManagerInterface)
	if !ok {
		log.Fatalf("Value in context is not of type crawler.CrawlManagerInterface")
	}
	if manager == nil {
		log.Fatalf("manager is nil")
	}

	url := viper.GetString("url")
	if url == "" {
		return fmt.Errorf("url is required")
	}

	if Debug {
		manager.Logger().Info("\nFlags:")
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			manager.Logger().Info(fmt.Sprintf(" %-12s : %s\n", flag.Name, flag.Value.String()))
		})

		manager.Logger().Info("\nRedis Environment Variables:")
		manager.Logger().Info(fmt.Sprintf(" %-12s : %s\n", "REDIS_HOST", viper.GetString("REDIS_HOST")))
		manager.Logger().Info(fmt.Sprintf(" %-12s : %s\n", "REDIS_PORT", viper.GetString("REDIS_PORT")))
		manager.Logger().Info(fmt.Sprintf(" %-12s : %s\n", "REDIS_AUTH", viper.GetString("REDIS_AUTH")))
	}

	_, err := manager.Crawl(ctx, url, viper.GetInt("maxdepth"), viper.GetBool("debug"))
	if err != nil {
		log.Fatalf("Error starting crawling: %v", err)
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

	RootCmd.AddCommand(crawlCmd)
}
