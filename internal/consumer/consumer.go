package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jonesrussell/page-prowler/internal/crawler"
)

type Output struct {
	Siteid    string    `json:"siteid"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Links     []Link    `json:"links"`
}

type Link struct {
	URL           string   `json:"url"`
	MatchingTerms []string `json:"matching_terms"`
}

func RetrieveAndUnmarshalLinks(ctx context.Context, manager crawler.CrawlManagerInterface, siteid string) ([]Link, error) {
	client := manager.Client()
	links, err := client.SMembers(ctx, siteid)
	if err != nil {
		return nil, fmt.Errorf("failed to get links from Redis: %v", err)
	}

	var linkStructs []Link
	for _, link := range links {
		var l Link
		unmarshalErr := json.Unmarshal([]byte(link), &l)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
		linkStructs = append(linkStructs, l)
	}

	return linkStructs, nil
}

func CreateOutput(siteid string, links []Link) Output {
	return Output{
		Siteid:    siteid,
		Timestamp: time.Now(),
		Status:    "success",
		Message:   "Links retrieved successfully",
		Links:     links,
	}
}

func MarshalOutput(output Output) ([]byte, error) {
	return json.Marshal(output)
}
