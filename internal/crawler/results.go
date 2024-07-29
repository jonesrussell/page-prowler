package crawler

import "github.com/jonesrussell/page-prowler/models"

// Results holds the results of the crawling process.
type Results struct {
	Pages []models.PageData
}

// NewResults creates a new instance of Results.
func NewResults() *Results {
	return &Results{
		Pages: make([]models.PageData, 0), // Initialize Pages slice
	}
}
