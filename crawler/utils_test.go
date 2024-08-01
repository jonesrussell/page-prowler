package crawler

import (
	"context"
	"testing"

	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/dbmanager"
	"github.com/jonesrussell/page-prowler/internal/termmatcher"
	"github.com/jonesrussell/page-prowler/models"
	"github.com/stretchr/testify/assert"
)

func TestHandleMatchingTerms(t *testing.T) {
	// Create a mock logger
	logger := loggo.NewMockLogger()

	// Create a mock DBManager
	dbManager := dbmanager.NewMockDBManager()

	// Create an actual TermMatcher
	termMatcher := termmatcher.NewTermMatcher(logger)

	cm := NewCrawlManager(logger, dbManager, nil, nil, nil)
	cm.TermMatcher = termMatcher
	cm.initializeStatsManager()

	// Define the input parameters
	options := &CrawlOptions{}
	currentURL := "https://www.example.com/the-cat-has-been-abducted"
	pageData := models.PageData{URL: currentURL}
	matchingTerms := []string{"abduct"}

	// Call the function
	err := cm.handleMatchingTerms(options, currentURL, pageData, matchingTerms)

	// Print the actual PageData
	t.Logf("Actual PageData: %+v\n", pageData)

	// Assert that there was no error
	assert.NoError(t, err)

	// Define the expected PageData
	expectedPageData := models.PageData{
		URL:             currentURL,
		MatchingTerms:   matchingTerms,
		SimilarityScore: 1, // Update this to the expected similarity score
	}

	// Assert that the result was saved to Redis
	ctx := context.Background()
	key := options.CrawlSiteID
	savedResults, err := dbManager.GetResultsFromRedis(ctx, key)

	// Print the actual results saved to Redis
	t.Logf("Actual results saved to Redis: %+v\n", savedResults)

	assert.NoError(t, err)
	assert.Contains(t, savedResults, expectedPageData) // Use the expected PageData
}
