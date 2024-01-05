package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

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
func EnqueueCrawlTask(client *asynq.Client, payload *CrawlTaskPayload) error {
	task, err := NewCrawlTask(payload)
	if err != nil {
		return err
	}
	_, err = client.Enqueue(task)
	return err
}

func NewCrawlTask(payload *CrawlTaskPayload) (*asynq.Task, error) {
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
