package crawler

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Results holds the results of the crawling process.
type Results struct {
	Pages []PageData
}

// NewResults creates a new instance of Results.
func NewResults() *Results {
	return &Results{
		Pages: make([]PageData, 0), // Initialize Pages slice
	}
}

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

// UpdatePageData updates the PageData with the provided href and matching terms.
// It sets the ParentURL and MatchingTerms fields of the PageData.
func (p *PageData) UpdatePageData(href string, matchingTerms []string) {
	p.MatchingTerms = matchingTerms
	p.ParentURL = href
}
