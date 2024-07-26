package cmd

import (
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/cobra"
)

// NewResultsCmd creates a new results command
func NewResultsCmd(manager crawler.CrawlManagerInterface) *cobra.Command {
	resultsCmd := &cobra.Command{
		Use:   "results",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			// Your command execution logic goes here
			// You can use similar logic as in your NewCrawlCmd function
			return nil
		},
	}

	return resultsCmd
}
