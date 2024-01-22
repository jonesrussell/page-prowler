package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// PageData represents the data of a crawled page.
type PageData struct {
	URL           string   `json:"url,omitempty"`
	ParentURL     string   `json:"parent_url,omitempty"` // New field for the parent URL
	SearchTerms   []string `json:"search_terms,omitempty"`
	MatchingTerms []string `json:"matching_terms,omitempty"`
	Error         string   `json:"error,omitempty"`
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

// logResults prints the results of the crawl.
func logResults(crawlerService *CrawlManager, results []PageData) {
	if len(results) == 0 {
		crawlerService.Logger.Info("No results to print")
		return
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		crawlerService.Logger.Error("Error occurred during marshaling", "error", err)
		return
	}

	crawlerService.Logger.Info(string(jsonData))
}

func (cs *CrawlManager) SaveResultsToRedis(ctx context.Context, results []PageData, key string) error {
	for _, result := range results {
		data, err := result.MarshalBinary()
		if err != nil {
			cs.Logger.Error("Error occurred during marshalling to binary", "error", err)
			return err
		}
		str := string(data)
		err = cs.Client.SAdd(ctx, key, str)
		if err != nil {
			cs.Logger.Error("Error occurred during saving to Redis", "error", err)
			return err
		}
		cs.Logger.Debug("Added elements to the set")

		// Debugging: Verify that the result was saved correctly
		isMember, err := cs.Client.SIsMember(ctx, key, str)
		if err != nil {
			cs.Logger.Error("Error occurred during checking membership in Redis set", "error", err)
			return err
		}
		if !isMember {
			cs.Logger.Error("Result was not saved correctly in Redis set", "result", str)
		} else {
			cs.Logger.Debug("Result was saved correctly in Redis set", "result", str)
		}
	}
	return nil
}
