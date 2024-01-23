package common

import "testing"

func TestCrawlManagerKey(t *testing.T) {
	if CrawlManagerKey == nil {
		t.Errorf("Expected CrawlManagerKey to be assigned, but it was nil")
	}
}
