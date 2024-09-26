package crawler

import (
	"context"
	"testing"

	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/dbmanager"
	"github.com/jonesrussell/page-prowler/internal/matcher"
	"github.com/jonesrussell/page-prowler/internal/termmatcher"
	"github.com/jonesrussell/page-prowler/models"
	"github.com/stretchr/testify/assert"
)

// MockMatcher is a simple implementation of the matcher interface for testing
type MockMatcher struct{}

func (m *MockMatcher) Match(content string) bool {
	// Implement mock logic for testing
	return content == "test"
}

func TestNewTermMatcher(t *testing.T) {
	logger := loggo.NewMockLogger()                   // Create a mock logger
	mockMatchers := []matcher.Matcher{&MockMatcher{}} // Create a slice of matchers

	// Pass both logger and matchers to NewTermMatcher
	tm := termmatcher.NewTermMatcher(logger, mockMatchers)

	// Add your test cases here
	if tm == nil {
		t.Error("Expected TermMatcher to be initialized, got nil")
	}
}

func TestHandleMatchingTerms(t *testing.T) {
	// Create a mock logger
	logger := loggo.NewMockLogger()

	// Create a mock DBManager
	dbManager := dbmanager.NewMockDBManager()
	// Create mock matchers for testing
	mockMatchers := []matcher.Matcher{&MockMatcher{}}

	// Create an actual TermMatcher with logger and mock matchers
	termMatcher := termmatcher.NewTermMatcher(logger, mockMatchers)

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
