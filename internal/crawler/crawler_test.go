package crawler_test

import (
	"context"
	"testing"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/crawler/internal/crawlResult"
	"github.com/jonesrussell/crawler/internal/crawler"
	"github.com/jonesrussell/crawler/internal/rediswrapper/mocks"
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
	redisMock := new(mocks.MockRedisWrapper)
	collector := colly.NewCollector()

	crawler.SetupCrawlingLogic(ctx, crawlSiteID, collector, searchTerms, &results, logger.Sugar(), redisMock)
	assert.NotNil(t, results)
}

func TestGetHostFromURL(t *testing.T) {
	url := "http://example.com"
	host := crawler.GetHostFromURL(url)
	assert.Equal(t, "example.com", host)
}
