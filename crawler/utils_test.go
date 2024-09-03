package crawler

import (
	"context"
	"testing"

	"github.com/gocolly/redisstorage"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/dbmanager"
	"github.com/jonesrussell/page-prowler/internal/termmatcher"
	"github.com/jonesrussell/page-prowler/models"
	"github.com/stretchr/testify/assert"
)

func TestHandleMatchingTerms(t *testing.T) {
	logger := loggo.NewMockLogger()
	dbManager := dbmanager.NewMockDBManager()
	contentProcessor := termmatcher.NewDefaultContentProcessor()
	termMatcher := termmatcher.NewTermMatcher(logger, 0.6, contentProcessor)

	collectorWrapper := &CollectorWrapper{}
	crawlOptions := &CrawlOptions{}
	redisStorage := &redisstorage.Storage{}

	cm := NewCrawlManager(logger, dbManager, collectorWrapper, crawlOptions, redisStorage, termMatcher)
	cm.initializeStatsManager()

	options := &CrawlOptions{
		CrawlSiteID: "test_crawl",
	}
	currentURL := "https://www.example.com/the-cat-has-been-abducted"
	pageData := models.PageData{URL: currentURL}
	matchingTerms := []string{"abduct"}

	err := cm.handleMatchingTerms(options, currentURL, pageData, matchingTerms)

	assert.NoError(t, err)

	ctx := context.Background()
	key := options.CrawlSiteID
	savedResults, err := dbManager.GetResultsFromRedis(ctx, key)

	assert.NoError(t, err)
	assert.Len(t, savedResults, 1)
	actualPageData := savedResults[0]

	assert.Equal(t, currentURL, actualPageData.URL)
	assert.Equal(t, matchingTerms, actualPageData.MatchingTerms)

	// Check if SimilarityScore is set
	assert.GreaterOrEqual(t, actualPageData.SimilarityScore, 0.0)
	assert.LessOrEqual(t, actualPageData.SimilarityScore, 1.0)

	// Check if stats were updated correctly
	assert.Equal(t, 1, cm.StatsManager.LinkStats.GetMatchedLinks())
	assert.Equal(t, 0, cm.StatsManager.LinkStats.GetNotMatchedLinks())
}

func TestUpdateStats(t *testing.T) {
	logger := loggo.NewMockLogger()
	dbManager := dbmanager.NewMockDBManager()
	contentProcessor := termmatcher.NewDefaultContentProcessor()
	termMatcher := termmatcher.NewTermMatcher(logger, 0.6, contentProcessor)

	collectorWrapper := &CollectorWrapper{}
	crawlOptions := &CrawlOptions{}
	redisStorage := &redisstorage.Storage{}

	cm := NewCrawlManager(logger, dbManager, collectorWrapper, crawlOptions, redisStorage, termMatcher)
	cm.initializeStatsManager()

	options := &CrawlOptions{}

	// Test with matching terms
	cm.UpdateStats(options, []string{"term1", "term2"})
	assert.Equal(t, 1, cm.StatsManager.LinkStats.GetMatchedLinks())
	assert.Equal(t, 0, cm.StatsManager.LinkStats.GetNotMatchedLinks())

	// Test with no matching terms
	cm.UpdateStats(options, []string{})
	assert.Equal(t, 1, cm.StatsManager.LinkStats.GetMatchedLinks())
	assert.Equal(t, 1, cm.StatsManager.LinkStats.GetNotMatchedLinks())
}
