package cmd

import (
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/worker"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start the Asynq worker",
	Run: func(cmd *cobra.Command, args []string) {
		concurrency := 10 // Replace with the concurrency level you want
		manager := cmd.Context().Value(managerKey).(*crawler.CrawlManager)
		worker.StartWorker(concurrency, manager, Debug)
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
