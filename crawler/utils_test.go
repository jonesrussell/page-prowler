package crawler

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/dbmanager"
	"github.com/jonesrussell/page-prowler/internal/matcher"
	"github.com/jonesrussell/page-prowler/internal/termmatcher"
	"github.com/jonesrussell/page-prowler/models"
	"github.com/stretchr/testify/assert"
)

// MockMatcher is a simple implementation of the matcher interface for testing
type MockMatcher struct{}

func (m *MockMatcher) Match(content string, pattern string) (bool, error) {
	// Implement mock logic for testing
	if content == "" || pattern == "" {
		return false, errors.New("content or pattern cannot be empty")
	}
	return strings.Contains(content, pattern), nil // Example logic
}

func TestNewTermMatcher(t *testing.T) {
	logger := loggo.NewMockLogger(gomock.NewController(t))
	mockMatchers := []matcher.Matcher{&MockMatcher{}}

	tm := termmatcher.NewTermMatcher(logger, mockMatchers)

	// Add your test cases here
	if tm == nil {
		t.Error("Expected TermMatcher to be initialized, got nil")
	}
}

func TestHandleMatchingTerms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock logger
	logger := loggo.NewMockLogger(ctrl)
	logger.EXPECT().Debug(gomock.Any()).AnyTimes()

	// Create a mock DBManager
	dbManager := dbmanager.NewMockDBManager()

	// Create an actual TermMatcher with a mock matcher
	mockMatcher := &MockMatcher{}
	termMatcher := termmatcher.NewTermMatcher(logger, []matcher.Matcher{mockMatcher})

	cm := NewCrawlManager(logger, dbManager, nil, nil, nil)
	cm.TermMatcher = termMatcher
	cm.initializeStatsManager()

	// Define the input parameters
	options := &CrawlOptions{CrawlSiteID: "test_crawl"}
	currentURL := "https://www.example.com/the-cat-has-been-abducted"
	pageData := models.PageData{URL: currentURL}
	matchingTerms := []string{"abduct"}

	// Call the function
	err := cm.handleMatchingTerms(options, currentURL, pageData, matchingTerms)

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
	savedResults, err := dbManager.GetResultsFromRedis(ctx, options.CrawlSiteID)

	// Additional assertions
	assert.NoError(t, err)
	assert.Contains(t, savedResults, expectedPageData)
	assert.Equal(t, 1, len(savedResults), "Expected one result to be saved")
	assert.Equal(t, expectedPageData, savedResults[0], "Saved result does not match expected data")
}
