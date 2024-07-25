package cmd

import (
	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/worker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewWorkerCmd creates a new worker command
func NewWorkerCmd() *cobra.Command {
	workerCmd := &cobra.Command{
		Use:   "worker",
		Short: "Start the Asynq worker",
		Run: func(cmd *cobra.Command, _ []string) {
			concurrency := 10 // Replace with the concurrency level you want
			manager := cmd.Context().Value(common.CrawlManagerKey).(*crawler.CrawlManager)
			debug := viper.GetBool("debug")
			worker.StartWorker(concurrency, manager, debug)
		},
	}

	return workerCmd
}
