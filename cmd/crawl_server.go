package cmd

import (
	"github.com/hibiken/asynq"
	"github.com/jonesrussell/page-prowler/internal/tasks"
)

func enqueueCrawlTask(client *asynq.Client, payload *tasks.CrawlTaskPayload) error {
	task, err := tasks.NewCrawlTask(payload)
	if err != nil {
		return err
	}
	_, err = client.Enqueue(task)
	return err
}
