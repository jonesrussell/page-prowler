package cmd

import (
	"github.com/spf13/cobra"
)

func CreateTestCommand(use, short, long string, runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		RunE:  runE,
	}
}
