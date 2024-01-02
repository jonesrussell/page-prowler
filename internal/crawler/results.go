package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// PageData represents the data of a crawled page.
type PageData struct {
	URL           string   `json:"url,omitempty"`            // The URL of the web page
	Links         []string `json:"links,omitempty"`          // The hyperlinks found on the web page
	SearchTerms   []string `json:"search_terms,omitempty"`   // The search terms used during the crawl
	MatchingTerms []string `json:"matching_terms,omitempty"` // The terms that matched the search criteria
	Error         string   `json:"error,omitempty"`          // Any error encountered during crawling of this page
}

func (p *PageData) Validate() error {
	// Check if the URL field is a valid URL
	_, err := url.ParseRequestURI(p.URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	// Add more checks as needed
	return nil
}

// MarshalBinary marshals the PageData into binary form.
func (p *PageData) MarshalBinary() ([]byte, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(p)
}

// UnmarshalBinary unmarshals binary data into PageData.
func (p *PageData) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, p); err != nil {
		return err
	}
	return p.Validate()
}

// printResults prints the results of the crawl.
func printResults(crawlerService *CrawlManager, results []PageData) {
	jsonData, err := json.Marshal(results)
	if err != nil {
		crawlerService.Logger.Error("Error occurred during marshaling", "error", err)
		return
	}

	crawlerService.Logger.Info(string(jsonData))
}

// SaveResultsToRedis saves the results of the crawl to Redis.
func (s *CrawlServer) SaveResultsToRedis(ctx context.Context, results []PageData, key string) error {
	// Debugging statement
	if ctx.Err() != nil {
		s.CrawlManager.Logger.Error("Crawl: context error:", ctx.Err())
	} else {
		s.CrawlManager.Logger.Info("Crawl: context is not done")
	}

	for _, result := range results {
		data, err := result.MarshalBinary()
		if err != nil {
			s.CrawlManager.Logger.Error("Error occurred during marshalling to binary", "error", err)
			return err
		}
		str := string(data)
		cmd := s.CrawlManager.Client.SAdd(ctx, key, str)
		count, err := cmd.Result()

		if err != nil {
			s.CrawlManager.Logger.Error("Error occurred during saving to Redis", "error", err)
			return err
		}
		s.CrawlManager.Logger.Info("Added", count, "elements to the set")
	}
	return nil
}
