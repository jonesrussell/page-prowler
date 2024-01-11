package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jonesrussell/page-prowler/internal/crawler"
)

type Output struct {
	Crawlsiteid string    `json:"crawlsiteid"`
	Timestamp   time.Time `json:"timestamp"`
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	Links       []Link    `json:"links"`
}

type Link struct {
	URL           string   `json:"url"`
	MatchingTerms []string `json:"matching_terms"`
}

func RetrieveAndUnmarshalLinks(ctx context.Context, manager *crawler.CrawlManager, crawlsiteid string) ([]Link, error) {
	links, err := manager.Client.SMembers(ctx, crawlsiteid)
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

func CreateOutput(crawlsiteid string, links []Link) Output {
	return Output{
		Crawlsiteid: crawlsiteid,
		Timestamp:   time.Now(),
		Status:      "success",
		Message:     "Links retrieved successfully",
		Links:       links,
	}
}

func MarshalOutput(output Output) ([]byte, error) {
	return json.Marshal(output)
}
