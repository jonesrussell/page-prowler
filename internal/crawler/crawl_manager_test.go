package crawler_test

import (
	"testing"

	"github.com/jonesrussell/page-prowler/mocks"
)

func TestNewCrawlManager(t *testing.T) {
	cm := mocks.NewCrawlManager()

	if cm.CrawlingMu == nil {
		t.Fatal("Expected CrawlingMu to be initialized, got nil")
	}
}
