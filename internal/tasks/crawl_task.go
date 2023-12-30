package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	CrawlTaskType = "crawl"
)

type CrawlTaskPayload struct {
	URL         string
	SearchTerms string
	CrawlSiteID string
	MaxDepth    int
	Debug       bool
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
