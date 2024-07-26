package crawler_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/mocks"
)

func TestNewCrawlManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggerField := mocks.NewMockLogger()
	client := mocks.NewMockClient()
	collector := mocks.NewMockCollectorInterface(ctrl)

	cm := crawler.NewCrawlManager(loggerField, client, collector, &crawler.CrawlOptions{})

	if cm.CrawlingMu == nil {
		t.Fatal("Expected CrawlingMu to be initialized, got nil")
	}
}
