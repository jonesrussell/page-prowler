package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/jonesrussell/page-prowler/crawler"
	"github.com/jonesrussell/page-prowler/internal/consumer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewGetLinksCmd creates a new getlinks command
func NewGetLinksCmd(manager crawler.CrawlManagerInterface) *cobra.Command {
	getLinksCmd := &cobra.Command{
		Use:   "getlinks",
		Short: "Get the list of links for a given siteid",
		RunE: func(cmd *cobra.Command, _ []string) error {
			siteid := viper.GetString("siteid")
			if siteid == "" {
				return ErrSiteidRequired
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

func printLinks(ctx context.Context, manager crawler.CrawlManagerInterface, siteid string) (consumer.Output, error) {
	links, err := consumer.RetrieveAndUnmarshalLinks(ctx, manager, siteid)
	if err != nil {
		return consumer.Output{}, err
	}

	output := consumer.CreateOutput(siteid, links)
	return output, nil
}
