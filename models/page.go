package models

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// PageData represents the data of a crawled page.
type PageData struct {
	URL             string   `json:"url,omitempty"`
	MatchingTerms   []string `json:"matching_terms,omitempty"`
	SimilarityScore float64  `json:"similarity_score,omitempty"`
	Error           string   `json:"error,omitempty"`
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

// UpdatePageData updates the PageData with the matching terms, and similarity score.
// It sets the MatchingTerms and SimilarityScore fields of the PageData.
func (p *PageData) UpdatePageData(matchingTerms []string, similarityScore float64) {
	p.MatchingTerms = matchingTerms
	p.SimilarityScore = similarityScore
}
