package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/consumer"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var GetLinksCmd = &cobra.Command{
	Use:   "getlinks",
	Short: "Get the list of links for a given crawlsiteid",
	RunE: func(cmd *cobra.Command, _ []string) error {
		log.Println("RunE function started")

		crawlsiteid := viper.GetString("crawlsiteid")
		log.Println(crawlsiteid)

		if crawlsiteid == "" {
			return ErrCrawlsiteidRequired
		}

		log.Println("CrawlSiteID: ", crawlsiteid) // Print the crawlsiteid

		if Debug {
			log.Println("Debug mode is enabled.")
		} else {
			log.Println("Debug mode is disabled.")
		}

		manager, ok := cmd.Context().Value(common.CrawlManagerKey).(*crawler.CrawlManager)
		if !ok || manager == nil {
			return fmt.Errorf("CrawlManager is not initialized")
		}

		log.Println("RunE function ending")
		err := printLinks(cmd.Context(), manager, crawlsiteid)
		if err != nil {
			log.Printf("Failed to print links: %v\n", err)
			return err
		}

		return nil
	},
}

func init() {
	GetLinksCmd.Flags().StringP("crawlsiteid", "s", "", "CrawlSite ID")

	rootCmd.AddCommand(GetLinksCmd)
}

func printJSON(jsonOutput []byte) error {
	_, err := fmt.Println(string(jsonOutput))
	if err != nil {
		return fmt.Errorf("failed to print links: %v", err)
	}
	return nil
}

func printLinks(ctx context.Context, manager *crawler.CrawlManager, crawlsiteid string) error {
	log.Println("printLinks function started")
	links, err := consumer.RetrieveAndUnmarshalLinks(ctx, manager, crawlsiteid)
	if err != nil {
		return err
	}

	output := consumer.CreateOutput(crawlsiteid, links)

	jsonOutput, err := consumer.MarshalOutput(output)
	if err != nil {
		return err
	}

	return printJSON(jsonOutput)
}
