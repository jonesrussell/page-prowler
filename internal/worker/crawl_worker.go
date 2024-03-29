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
	l.logger.Debug(fmt.Sprint(args...))
}

func (l *AsynqLoggerWrapper) Info(args ...interface{}) {
	l.logger.Debug(fmt.Sprint(args...))
}

func (l *AsynqLoggerWrapper) Warn(args ...interface{}) {
	l.logger.Debug(fmt.Sprint(args...))
}

func (l *AsynqLoggerWrapper) Error(args ...interface{}) {
	l.logger.Debug(fmt.Sprint(args...))
}

func (l *AsynqLoggerWrapper) Fatal(args ...interface{}) {
	l.logger.Debug(fmt.Sprint(args...))
}

// Implement the rest of the asynq.Logger methods in a similar way

func handleCrawlTask(ctx context.Context, task *asynq.Task, cm *crawler.CrawlManager, debug bool) error {
	var payload tasks.CrawlTaskPayload
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return err
	}

	_, err = cm.Crawl(ctx, payload.URL, payload.SearchTerms, payload.CrawlSiteID, payload.MaxDepth, debug)
	return err
}

func StartWorker(concurrency int, cm *crawler.CrawlManager, debug bool) {
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
		cm.Logger().Fatal(fmt.Sprintf("could not run server: %v", err))
	}
}
