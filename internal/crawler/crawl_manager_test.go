package crawler

import (
	"testing"

	"github.com/jonesrussell/page-prowler/cmd/mocks"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
)

func TestNewCrawlManager(t *testing.T) {
	loggerField := mocks.NewMockLogger()
	client := prowlredis.NewMockClient()
	mongoDBWrapper := mocks.NewMockMongoDBWrapper()

	cm := NewCrawlManager(loggerField, client, mongoDBWrapper)

	if cm.CrawlingMu == nil {
		t.Fatal("Expected CrawlingMu to be initialized, got nil")
	}
}
