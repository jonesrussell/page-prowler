package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/tasks"
	"github.com/jonesrussell/page-prowler/redis"
)

func handleCrawlTask(ctx context.Context, task *asynq.Task, crawlerService *crawler.CrawlManager) error {
	var payload tasks.CrawlTaskPayload
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return err
	}

	// Debugging statements
	log.Println("Payload URL:", payload.URL)
	log.Println("Payload SearchTerms:", payload.SearchTerms)
	log.Println("Payload CrawlSiteID:", payload.CrawlSiteID)
	log.Println("Payload MaxDepth:", payload.MaxDepth)
	log.Println("Payload Debug:", payload.Debug)

	server := &crawler.CrawlServer{} // Initialize your CrawlServer

	return crawler.StartCrawling(ctx, payload.URL, payload.SearchTerms, payload.CrawlSiteID, payload.MaxDepth, payload.Debug, crawlerService, server)
}

func StartWorker(concurrency int, crawlerService *crawler.CrawlManager) {
	ctx := context.Background()

	// Create a new Redis client
	redisClient, err := redis.NewClient(ctx, "localhost", "", "6379")
	if err != nil {
		log.Fatalf("Failed to create Redis client: %v", err)
	}

	// Initialize a new Asynq server with the default settings.
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     redisClient.Options().Addr,
			Password: redisClient.Options().Password,
			DB:       redisClient.Options().DB,
		},
		asynq.Config{
			Concurrency: concurrency,
		},
	)

	// mux maps a task type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.CrawlTaskType, func(ctx context.Context, task *asynq.Task) error {
		// Debugging statement
		log.Println("Task payload:", string(task.Payload()))
		return handleCrawlTask(ctx, task, crawlerService)
	})

	// Run the server with the handler mux.
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
