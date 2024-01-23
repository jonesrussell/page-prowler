package crawler

import (
	"testing"
)

func TestNewCrawlOptions(t *testing.T) {
	crawlSiteID := "testSiteID"
	searchTerms := []string{"term1", "term2"}
	debug := true
	var results []PageData

	co := NewCrawlOptions(crawlSiteID, searchTerms, debug, &results)

	if co.CrawlSiteID != crawlSiteID {
		t.Errorf("Expected CrawlSiteID to be %v, got %v", crawlSiteID, co.CrawlSiteID)
	}

	if len(co.SearchTerms) != len(searchTerms) {
		t.Errorf("Expected SearchTerms length to be %v, got %v", len(searchTerms), len(co.SearchTerms))
	}

	for i, term := range searchTerms {
		if co.SearchTerms[i] != term {
			t.Errorf("Expected SearchTerms[%d] to be %v, got %v", i, term, co.SearchTerms[i])
		}
	}

	if co.Debug != debug {
		t.Errorf("Expected Debug to be %v, got %v", debug, co.Debug)
	}
}
