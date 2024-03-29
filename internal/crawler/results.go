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

// Validate checks if the PageData fields are valid.
// It returns an error if the URL field is not a valid URL.
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
// It returns an error if the PageData is not valid or if the marshaling fails.
func (p *PageData) MarshalBinary() ([]byte, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(p)
}

// UnmarshalBinary unmarshals binary data into PageData.
// It returns an error if the unmarshaling fails or if the PageData is not valid after unmarshaling.
func (p *PageData) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, p); err != nil {
		return err
	}
	return p.Validate()
}

// SaveResultsToRedis saves the crawling results to a Redis set.
// It marshals each PageData into a JSON string and adds it to the set.
// Parameters:
// - ctx: The context for the Redis operation.
// - results: The slice of PageData to save.
// - key: The Redis key to use for the set.
// Returns:
// - error: An error if the marshaling or saving to Redis fails.
func (cm *CrawlManager) SaveResultsToRedis(ctx context.Context, results []PageData, key string) error {
	cm.LoggerField.Debug(fmt.Sprintf("SaveResultsToRedis: Number of results before processing: %d", len(results)))

	for _, result := range results {
		cm.LoggerField.Debug(fmt.Sprintf("SaveResultsToRedis: Processing result %v", result))

		data, err := json.Marshal(result)
		if err != nil {
			cm.LoggerField.Error(fmt.Sprintf("SaveResultsToRedis: Error occurred during marshalling to JSON: %v", err))
			return err
		}
		str := string(data)
		err = cm.Client.SAdd(ctx, key, str)
		if err != nil {
			cm.LoggerField.Error(fmt.Sprintf("SaveResultsToRedis: Error occurred during saving to Redis: %v", err))
			return err
		}
		cm.LoggerField.Debug("SaveResultsToRedis: Added elements to the set")

		// Debugging: Verify that the result was saved correctly
		isMember, err := cm.Client.SIsMember(ctx, key, str)
		if err != nil {
			cm.LoggerField.Error(fmt.Sprintf("SaveResultsToRedis: Error occurred during checking membership in Redis set: %v", err))
			return err
		}
		if !isMember {
			cm.LoggerField.Error(fmt.Sprintf("SaveResultsToRedis: Result was not saved correctly in Redis set: %v", str))
		} else {
			cm.LoggerField.Debug(fmt.Sprintf("SaveResultsToRedis: Result was saved correctly in Redis set, key: %s, result: %s", key, str))
		}
	}

	cm.LoggerField.Debug(fmt.Sprintf("SaveResultsToRedis: Number of results after processing: %d", len(results)))

	return nil
}

// UpdatePageData updates the PageData with the provided href and matching terms.
// It sets the ParentURL and MatchingTerms fields of the PageData.
func (p *PageData) UpdatePageData(href string, matchingTerms []string) {
	p.MatchingTerms = matchingTerms
	p.ParentURL = href
}

// AppendResult appends a PageData to the Results slice in the CrawlOptions.
// Parameters:
// - options: The CrawlOptions containing the Results slice.
// - pageData: The PageData to append to the Results slice.
func (cm *CrawlManager) AppendResult(options *CrawlOptions, pageData PageData) {
	*options.Results = append(*options.Results, pageData)
}
