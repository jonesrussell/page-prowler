package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

// AsynqClient defines an interface with the methods you use from asynq.Client.
type AsynqClient interface {
	Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

const (
	CrawlTaskType = "crawl"
)

type CrawlTaskPayload struct {
	URL         string `json:"url"`
	SearchTerms string `json:"search_terms"`
	CrawlSiteID string `json:"crawl_site_id"`
	MaxDepth    int    `json:"max_depth"`
	Debug       bool   `json:"debug"`
}

// EnqueueCrawlTask creates asynq task
func EnqueueCrawlTask(client AsynqClient, payload *CrawlTaskPayload) (string, error) {
	task, err := NewCrawlTask(payload)
	if err != nil {
		return "", err
	}
	info, err := client.Enqueue(task)
	if err != nil {
		return "", err
	}
	return info.ID, nil
}

func NewCrawlTask(payload *CrawlTaskPayload) (*asynq.Task, error) {
	// Validate the payload
	if payload.URL == "" || payload.SearchTerms == "" || payload.CrawlSiteID == "" || payload.MaxDepth < 0 {
		return nil, fmt.Errorf("invalid payload")
	}

	data, err := json.Marshal(map[string]interface{}{
		"url":           payload.URL,
		"search_terms":  payload.SearchTerms,
		"crawl_site_id": payload.CrawlSiteID,
		"max_depth":     payload.MaxDepth,
		"debug":         payload.Debug,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(CrawlTaskType, data), nil
}
