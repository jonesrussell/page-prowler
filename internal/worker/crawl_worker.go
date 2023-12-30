package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/tasks"
)

func handleCrawlTask(ctx context.Context, task *asynq.Task) error {
	var payload tasks.CrawlTaskPayload
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return err
	}

	crawlerService := &crawler.CrawlManager{} // Initialize your CrawlManager
	server := &crawler.CrawlServer{}          // Initialize your CrawlServer

	return crawler.StartCrawling(ctx, payload.URL, payload.SearchTerms, payload.CrawlSiteID, payload.MaxDepth, payload.Debug, crawlerService, server)
}

func StartWorker() {
	// Initialize a new Asynq server with the default settings.
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"}, // replace with your Redis server address
		asynq.Config{
			Concurrency: 10, // replace with the concurrency level you want
		},
	)

	// mux maps a task type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.CrawlTaskType, handleCrawlTask)

	// Run the server with the handler mux.
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
