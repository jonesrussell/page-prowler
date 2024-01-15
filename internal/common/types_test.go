package common

import "testing"

func TestCrawlManagerKey(t *testing.T) {
	if CrawlManagerKey != "cm" {
		t.Errorf("Expected CrawlManagerKey to be 'cm', but got '%s'", CrawlManagerKey)
	}
}
