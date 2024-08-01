package cmd

import (
	"github.com/jonesrussell/page-prowler/crawler"
	"github.com/jonesrussell/page-prowler/internal/worker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewWorkerCmd creates a new worker command
func NewWorkerCmd(manager crawler.CrawlManagerInterface) *cobra.Command {
	workerCmd := &cobra.Command{
		Use:   "worker",
		Short: "Start the Asynq worker",
		Run: func(_ *cobra.Command, _ []string) {
			concurrency := 10 // Replace with the concurrency level you want
			debug := viper.GetBool("debug")
			worker.StartWorker(concurrency, manager, debug)
		},
	}

	return workerCmd
}
