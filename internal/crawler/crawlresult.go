// Package crawlresult defines the data structures used for storing
// results of a web crawl operation.
package crawler

import (
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
func (p PageData) MarshalBinary() ([]byte, error) {
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
