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
	ParentURL     string   `json:"parent_url,omitempty"`
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
		crawlerService.LoggerField.Info("No results to print", map[string]interface{}{})
		return
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		crawlerService.LoggerField.Error("Error occurred during marshaling", map[string]interface{}{"error": err})
		return
	}

	crawlerService.LoggerField.Info(string(jsonData), map[string]interface{}{})
}

func (cm *CrawlManager) SaveResultsToRedis(ctx context.Context, results []PageData, key string) error {
	cm.LoggerField.Debug("SaveResultsToRedis: Number of results before processing", map[string]interface{}{"count": len(results)})

	for _, result := range results {
		cm.LoggerField.Debug("SaveResultsToRedis: Processing result", map[string]interface{}{"result": result})

		data, err := json.Marshal(result)
		if err != nil {
			cm.LoggerField.Error("SaveResultsToRedis: Error occurred during marshalling to JSON", map[string]interface{}{"error": err})
			return err
		}
		str := string(data)
		err = cm.Client.SAdd(ctx, key, str)
		if err != nil {
			cm.LoggerField.Error("SaveResultsToRedis: Error occurred during saving to Redis", map[string]interface{}{"error": err})
			return err
		}
		cm.LoggerField.Debug("SaveResultsToRedis: Added elements to the set", nil)

		// Debugging: Verify that the result was saved correctly
		isMember, err := cm.Client.SIsMember(ctx, key, str)
		if err != nil {
			cm.LoggerField.Error("SaveResultsToRedis: Error occurred during checking membership in Redis set", map[string]interface{}{"error": err})
			return err
		}
		if !isMember {
			cm.LoggerField.Error("SaveResultsToRedis: Result was not saved correctly in Redis set", map[string]interface{}{"result": str})
		} else {
			cm.LoggerField.Debug("SaveResultsToRedis: Result was saved correctly in Redis set", map[string]interface{}{"key": key, "result": str})
		}
	}

	cm.LoggerField.Debug("SaveResultsToRedis: Number of results after processing", map[string]interface{}{"count": len(results)})

	return nil
}

func (p *PageData) UpdatePageData(href string, matchingTerms []string) {
	p.MatchingTerms = matchingTerms
	p.ParentURL = href
}

func (cm *CrawlManager) AppendResult(options *CrawlOptions, pageData PageData) {
	*options.Results = append(*options.Results, pageData)
}
