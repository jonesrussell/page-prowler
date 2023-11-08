package crawler_test

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/gocolly/colly"
	"github.com/jonesrussell/crawler/internal/crawlResult"
	"github.com/jonesrussell/crawler/internal/crawler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestLoadConfiguration(t *testing.T) {
	cfg, err := crawler.LoadConfiguration()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestInitializeLogger(t *testing.T) {
	logger, err := crawler.InitializeLogger(true)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestConfigureCollector(t *testing.T) {
	allowedDomains := []string{"example.com"}
	maxDepth := 3
	collector := crawler.ConfigureCollector(allowedDomains, maxDepth)
	assert.NotNil(t, collector)
}

func TestHandleHTMLParsing(t *testing.T) {
	// Mock objects and the setup logic for the test will go here
}

func TestHandleMatchingLinks(t *testing.T) {
	// Mock objects and the setup logic for the test will go here
}

func TestHandleNonMatchingLinks(t *testing.T) {
	// Mock objects and the setup logic for the test will go here
}

func TestHandleRedisOperations(t *testing.T) {
	// Mock objects and the setup logic for the test will go here
}

func TestHandleErrorEvents(t *testing.T) {
	// Mock objects and the setup logic for the test will go here
}

func TestSetupCrawlingLogic(t *testing.T) {
	ctx := context.Background()
	crawlSiteID := "testSiteID"
	searchTerms := []string{"test"}
	results := []crawlResult.PageData{}
	logger, _ := zap.NewProduction()
	db, mock := redismock.NewClientMock() // Create a new mock Redis client
	collector := colly.NewCollector()

	// Set your expectations here
	// For example, if you expect the SAdd function to be called once with "testSiteID" and "testURL",
	// and it should return a NewIntResult(1, nil), you can do:
	mock.ExpectSAdd(crawlSiteID, "testURL").SetVal(1)

	crawler.SetupCrawlingLogic(ctx, crawlSiteID, collector, searchTerms, &results, logger.Sugar(), db) // Use the mock client here

	// Check if all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetHostFromURL(t *testing.T) {
	url := "http://example.com"
	host := crawler.GetHostFromURL(url)
	assert.Equal(t, "example.com", host)
}
