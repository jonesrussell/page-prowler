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

// AsynqLoggerWrapper wraps the logger to provide different log levels.
type AsynqLoggerWrapper struct {
	logger logger.Logger
}

// Debug logs debug level messages.
func (l *AsynqLoggerWrapper) Debug(args ...interface{}) {
	l.logger.Debug(fmt.Sprint(args...), nil)
}

// Info logs info level messages.
func (l *AsynqLoggerWrapper) Info(args ...interface{}) {
	l.logger.Info(fmt.Sprint(args...), nil)
}

// Warn logs warning level messages.
func (l *AsynqLoggerWrapper) Warn(args ...interface{}) {
	l.logger.Warn(fmt.Sprint(args...), nil)
}

// Error logs error level messages.
func (l *AsynqLoggerWrapper) Error(args ...interface{}) {
	l.logger.Error(fmt.Sprint(args...), nil)
}

// Fatal logs fatal level messages.
func (l *AsynqLoggerWrapper) Fatal(args ...interface{}) {
	l.logger.Fatal(fmt.Sprint(args...), nil)
}

// CrawlerWorkerInterface defines the methods a CrawlerWorker should have.
type CrawlerWorkerInterface interface {
	StartCrawlerWorker(
		concurrency int,
		crawlManager crawler.CrawlManagerInterface,
		debug bool,
	)
}

// CrawlerWorker is a struct that implements the CrawlerWorkerInterface.
type CrawlerWorker struct {
}

// StartCrawlerWorker starts the worker with the specified concurrency level.
func (w *CrawlerWorker) StartCrawlerWorker(
	concurrency int,
	crawlManager crawler.CrawlManagerInterface,
	debug bool,
) {
	cm, ok := crawlManager.(*crawler.CrawlManager)
	if !ok {
		cm.Logger().Error("failed to assert type of crawlManager", nil)
		return
	}

	// Initialize a new Asynq server with the default settings.
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cm.Client.Options().Addr,
			Password: cm.Client.Options().Password,
			DB:       cm.Client.Options().DB,
		},
		asynq.Config{
			Concurrency: concurrency,
			Logger:      &AsynqLoggerWrapper{logger: cm.Logger()}, // Use the Logger from CrawlManager
		},
	)

	// mux maps a task type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.CrawlTaskType, func(ctx context.Context, task *asynq.Task) error {
		return handleCrawlTask(ctx, task, cm, debug)
	})

	// Run the server with the handler mux.
	if err := srv.Run(mux); err != nil {
		cm.Logger().Fatal("could not run server", map[string]interface{}{"error": err})
	}
}

// NewCrawlerWorker returns a new instance of CrawlerWorker that implements CrawlerWorkerInterface.
func NewCrawlerWorker() CrawlerWorkerInterface {
	return &CrawlerWorker{}
}

func handleCrawlTask(ctx context.Context, task *asynq.Task, crawlerService *crawler.CrawlManager, debug bool) error {
	var payload tasks.CrawlTaskPayload
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return err
	}

	return crawlerService.StartCrawling(
		ctx,
		payload.URL,
		payload.SearchTerms,
		payload.CrawlSiteID,
		payload.MaxDepth,
		debug,
	)
}
