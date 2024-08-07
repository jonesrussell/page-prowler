package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/crawler"
	"github.com/jonesrussell/page-prowler/internal/tasks"
)

type AsynqLoggerWrapper struct {
	logger loggo.LoggerInterface
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

func handleCrawlTask(task *asynq.Task, cm crawler.CrawlManagerInterface, debug bool) error {
	var payload tasks.CrawlTaskPayload
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return err
	}

	searchTermsSlice := strings.Split(payload.SearchTerms, ",")

	options := crawler.CrawlOptions{
		StartURL:    payload.URL,
		MaxDepth:    payload.MaxDepth,
		SearchTerms: searchTermsSlice,
		Debug:       debug,
	}

	err = cm.SetOptions(&options)
	if err != nil {
		return err
	}

	err = cm.Crawl()
	return err
}

func StartWorker(concurrency int, manager crawler.CrawlManagerInterface, debug bool) {
	dbManager := manager.GetDBManager()

	// Initialize a new Asynq server with the default settings.
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     dbManager.RedisOptions().Addr,
			Password: dbManager.RedisOptions().Password,
			DB:       dbManager.RedisOptions().DB,
		},
		asynq.Config{
			Concurrency: concurrency,
			Logger:      &AsynqLoggerWrapper{logger: manager.GetLogger()}, // Use the Logger from CrawlManager
		},
	)

	// mux maps a task type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.CrawlTaskType, func(_ context.Context, task *asynq.Task) error {
		return handleCrawlTask(task, manager, debug) // Pass the manager to the handleCrawlTask function
	})

	// Run the server with the handler mux.
	if err := srv.Run(mux); err != nil {
		manager.GetLogger().Fatal(fmt.Sprintf("could not run server: %v", err), nil)
	}
}
