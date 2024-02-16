package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/consumer"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
)

var GetLinksCmd = &cobra.Command{
	Use:   "getlinks",
	Short: "Get the list of links for a given siteid",
	RunE: func(cmd *cobra.Command, _ []string) error {
		manager, ok := cmd.Context().Value(common.CrawlManagerKey).(*crawler.CrawlManager)
		if !ok || manager == nil {
			return fmt.Errorf("CrawlManager is not initialized")
		}

		err := printLinks(cmd.Context(), manager, Siteid)
		if err != nil {
			log.Printf("Failed to print links: %v\n", err)
			return err
		}

		return nil
	},
}

func init() {
	resultsCmd.AddCommand(GetLinksCmd)
}

func printJSON(jsonOutput []byte) error {
	_, err := fmt.Println(string(jsonOutput))
	if err != nil {
		return fmt.Errorf("failed to print links: %v", err)
	}
	return nil
}

func printLinks(ctx context.Context, manager *crawler.CrawlManager, siteid string) error {
	links, err := consumer.RetrieveAndUnmarshalLinks(ctx, manager, siteid)
	if err != nil {
		return err
	}

	output := consumer.CreateOutput(siteid, links)

	jsonOutput, err := consumer.MarshalOutput(output)
	if err != nil {
		return err
	}

	return printJSON(jsonOutput)
}
