package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/tasks"
)

type AsynqLoggerWrapper struct {
	logger logger.Logger
}

func (l *AsynqLoggerWrapper) Debug(args ...interface{}) {
	l.logger.Debug(fmt.Sprint(args...), nil)
}

func (l *AsynqLoggerWrapper) Info(args ...interface{}) {
	l.logger.Info(fmt.Sprint(args...), nil)
}

func (l *AsynqLoggerWrapper) Warn(args ...interface{}) {
	l.logger.Warn(fmt.Sprint(args...), nil)
}

func (l *AsynqLoggerWrapper) Error(args ...interface{}) {
	l.logger.Error(fmt.Sprint(args...), nil)
}

func (l *AsynqLoggerWrapper) Fatal(args ...interface{}) {
	l.logger.Fatal(fmt.Sprint(args...), nil)
}

// Define an interface for the worker
type WorkerInterface interface {
	StartWorker(concurrency int, crawlerService crawler.CrawlManagerInterface, debug bool)
}

// Define a struct that implements the Workerinterface
type Worker struct {
}

// Implement the StartWorker method for the crawlWorkerImpl struct
func (w *Worker) StartWorker(concurrency int, crawlerService crawler.CrawlManagerInterface, debug bool) {
	// Initialize a new Asynq server with the default settings.
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     crawlerService.(*crawler.CrawlManager).Client.Options().Addr,
			Password: crawlerService.(*crawler.CrawlManager).Client.Options().Password,
			DB:       crawlerService.(*crawler.CrawlManager).Client.Options().DB,
		},
		asynq.Config{
			Concurrency: concurrency,
			Logger:      &AsynqLoggerWrapper{logger: crawlerService.Logger()}, // Use the Logger from CrawlManager
		},
	)

	// mux maps a task type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.CrawlTaskType, func(ctx context.Context, task *asynq.Task) error {
		return handleCrawlTask(ctx, task, crawlerService.(*crawler.CrawlManager), debug)
	})

	// Run the server with the handler mux.
	if err := srv.Run(mux); err != nil {
		crawlerService.Logger().Fatal("could not run server", map[string]interface{}{"error": err})
	}
}

// NewWorker returns a new instance of CrawlWorker
func NewWorker() Worker {
	return Worker{}
}

func handleCrawlTask(ctx context.Context, task *asynq.Task, crawlerService *crawler.CrawlManager, debug bool) error {
	var payload tasks.CrawlTaskPayload
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return err
	}

	return crawlerService.StartCrawling(ctx, payload.URL, payload.SearchTerms, payload.CrawlSiteID, payload.MaxDepth, debug)
}
