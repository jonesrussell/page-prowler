package cmd

import (
	"context"
	"fmt"

	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/consumer"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
)

func init() {
	getLinksCmd.Flags().StringVarP(&Crawlsiteid, "crawlsiteid", "s", "", "CrawlSite ID")
}

var getLinksCmd = &cobra.Command{
	Use:   "getlinks",
	Short: "Get the list of links for a given crawlsiteid",
	RunE:  getLinks,
}

func getLinks(cmd *cobra.Command, _ []string) error {
	if Crawlsiteid == "" {
		return ErrCrawlsiteidRequired
	}

	manager, ok := cmd.Context().Value(common.ManagerKey).(*crawler.CrawlManager)
	if !ok || manager == nil {
		return fmt.Errorf("CrawlManager is not initialized")
	}

	return printLinks(cmd.Context(), manager, Crawlsiteid)
}

func printJSON(jsonOutput []byte) error {
	_, err := fmt.Println(string(jsonOutput))
	if err != nil {
		return fmt.Errorf("failed to print links: %v", err)
	}
	return nil
}

func printLinks(ctx context.Context, manager *crawler.CrawlManager, crawlsiteid string) error {
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
