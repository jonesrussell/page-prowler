package mocks

import (
	"sync"

	colly "github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/crawler"
)

// MockCrawlManager is a mock implementation of the CrawlManagerInterface.
type MockCrawlManager struct {
	crawler.CrawlManager
	// Add any additional fields or methods needed for your mock.
}

// NewCrawlManager returns a new instance of the MockCrawlManager.
func NewCrawlManager() *MockCrawlManager {
	return &MockCrawlManager{
		CrawlManager: crawler.CrawlManager{
			LoggerField:    NewMockLogger(),
			Client:         NewMockClient(),
			MongoDBWrapper: NewMockMongoDBWrapper(),
			CrawlingMu:     &sync.Mutex{},
			Collector:      colly.NewCollector(),
			StatsManager:   crawler.NewStatsManager(),
		},
	}
}

func NewCollector() *colly.Collector {
	return colly.NewCollector()
}
