package crawler_test

import (
	"testing"
	"time"

	"github.com/jonesrussell/page-prowler/internal/crawler"
)

func TestNewCrawlOptions(t *testing.T) {
	// Initialize a new CrawlOptions object
	options := crawler.CrawlOptions{}

	// Set some sample values
	options.CrawlSiteID = "test-site-id"
	options.Debug = true
	options.DelayBetweenRequests = 500 * time.Millisecond
	options.MaxConcurrentRequests = 10
	options.MaxDepth = 3
	options.SearchTerms = []string{"term1", "term2"}
	options.StartURL = "http://example.com"

	// Verify that the fields are correctly set
	if options.CrawlSiteID != "test-site-id" {
		t.Errorf("Expected CrawlSiteID to be 'test-site-id', got '%s'", options.CrawlSiteID)
	}

	if options.Debug != true {
		t.Error("Expected Debug to be true")
	}

	if options.DelayBetweenRequests != 500*time.Millisecond {
		t.Errorf("Expected DelayBetweenRequests to be 500ms, got %v", options.DelayBetweenRequests)
	}

	if options.MaxConcurrentRequests != 10 {
		t.Errorf("Expected MaxConcurrentRequests to be 10, got %d", options.MaxConcurrentRequests)
	}

	if len(options.SearchTerms) != 2 || options.SearchTerms[0] != "term1" || options.SearchTerms[1] != "term2" {
		t.Errorf("Expected SearchTerms to be ['term1', 'term2'], got %v", options.SearchTerms)
	}

	if options.StartURL != "http://example.com" {
		t.Errorf("Expected StartURL to be 'http://example.com', got '%s'", options.StartURL)
	}
}
