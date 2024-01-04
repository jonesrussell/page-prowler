package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/tasks"
	"github.com/jonesrussell/page-prowler/redis"
)

type CustomLogger struct {
	logger.Logger
}

func (cl *CustomLogger) Debug(args ...interface{}) {
	msg := fmt.Sprintf("%v", args)
	cl.Logger.Debug(msg)
}

func (cl *CustomLogger) Info(args ...interface{}) {
	msg := fmt.Sprintf("%v", args)
	cl.Logger.Info(msg)
}

func (cl *CustomLogger) Warn(args ...interface{}) {
	msg := fmt.Sprintf("%v", args)
	cl.Logger.Warn(msg)
}

func (cl *CustomLogger) Error(args ...interface{}) {
	msg := fmt.Sprintf("%v", args)
	cl.Logger.Error(msg)
}

func (cl *CustomLogger) Fatal(args ...interface{}) {
	msg := fmt.Sprintf("%v", args)
	cl.Logger.Fatal(msg)
}

func handleCrawlTask(ctx context.Context, task *asynq.Task, crawlerService *crawler.CrawlManager, debug bool) error {
	var payload tasks.CrawlTaskPayload
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return err
	}

	return crawlerService.StartCrawling(ctx, payload.URL, payload.SearchTerms, payload.CrawlSiteID, payload.MaxDepth, debug)
}

func StartWorker(concurrency int, crawlerService *crawler.CrawlManager, debug bool) {
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
			Logger:      &CustomLogger{logger.New(debug, logger.DefaultLogLevel)},
		},
	)

	// mux maps a task type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.CrawlTaskType, func(ctx context.Context, task *asynq.Task) error {
		return handleCrawlTask(ctx, task, crawlerService, debug)
	})

	// Run the server with the handler mux.
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
