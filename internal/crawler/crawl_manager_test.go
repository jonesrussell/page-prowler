package crawler_test

import (
	"testing"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/jonesrussell/page-prowler/mocks"
)

func TestNewCrawlManager(t *testing.T) {
	loggerField := mocks.NewMockLogger()
	client := prowlredis.NewMockClient()
	mongoDBWrapper := mocks.NewMockMongoDBWrapper()

	cm := crawler.NewCrawlManager(loggerField, client, mongoDBWrapper)

	if cm.CrawlingMu == nil {
		t.Fatal("Expected CrawlingMu to be initialized, got nil")
	}
}
