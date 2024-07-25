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

// NewGetLinksCmd creates a new getlinks command
func NewGetLinksCmd() *cobra.Command {
	getLinksCmd := &cobra.Command{
		Use:   "getlinks",
		Short: "Get the list of links for a given siteid",
		RunE: func(cmd *cobra.Command, _ []string) error {
			siteid := viper.GetString("siteid")
			if siteid == "" {
				return ErrSiteidRequired
			}

			manager, ok := cmd.Context().Value(common.CrawlManagerKey).(*crawler.CrawlManager)
			if !ok || manager == nil {
				return fmt.Errorf("CrawlManager is not initialized")
			}

			output, err := printLinks(cmd.Context(), manager, siteid)
			if err != nil {
				log.Printf("Failed to print links: %v\n", err)
				return err
			}

			jsonOutput, err := consumer.MarshalOutput(output)
			if err != nil {
				return err
			}

			err = printJSON(jsonOutput)
			if err != nil {
				log.Printf("Failed to print JSON output: %v\n", err)
				return err
			}

			return nil
		},
	}

	return getLinksCmd
}

func printJSON(jsonOutput []byte) error {
	_, err := fmt.Println(string(jsonOutput))
	if err != nil {
		return fmt.Errorf("failed to print links: %v", err)
	}
	return nil
}

func printLinks(ctx context.Context, manager *crawler.CrawlManager, siteid string) (consumer.Output, error) {
	links, err := consumer.RetrieveAndUnmarshalLinks(ctx, manager, siteid)
	if err != nil {
		return consumer.Output{}, err
	}

	output := consumer.CreateOutput(siteid, links)
	return output, nil
}
