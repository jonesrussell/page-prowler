package cmd

import (
	"context"

	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/worker"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start the Asynq worker",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		w := worker.NewCrawlerWorker()                                               // Create a new CrawlerWorker instance
		cmd.SetContext(context.WithValue(cmd.Context(), common.CrawlerWorkerKey, w)) // Set the CrawlerWorker in the context
		return nil
	},
	RunE: CrawlerWorkerCmdRun,
}

func CrawlerWorkerCmdRun(cmd *cobra.Command, args []string) error {
	concurrency := 10 // Replace with the concurrency level you want
	manager := cmd.Context().Value(common.CrawlManagerKey).(crawler.CrawlManagerInterface)

	w := cmd.Context().Value(common.CrawlerWorkerKey).(worker.CrawlerWorkerInterface) // Retrieve the CrawlerWorker from the context
	w.StartCrawlerWorker(concurrency, manager, Debug)                                 // Start the CrawlerWorker
	return nil
}

func init() {
	RootCmd.AddCommand(workerCmd)
}
