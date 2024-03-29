package cmd

import (
	"github.com/spf13/cobra"
)

// resultsCmd represents the results command
var resultsCmd = &cobra.Command{
	Use:   "results",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	RootCmd.AddCommand(resultsCmd)
}
