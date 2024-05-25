package crawler_test

import (
	"testing"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/mocks"
)

func TestNewCrawlManager(t *testing.T) {
	loggerField := mocks.NewMockLogger()
	client := mocks.NewMockClient()

	cm := crawler.NewCrawlManager(loggerField, client, &crawler.CrawlOptions{})

	if cm.CrawlingMu == nil {
		t.Fatal("Expected CrawlingMu to be initialized, got nil")
	}
}
