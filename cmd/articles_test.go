package cmd

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/jonesrussell/page-prowler/cmd/mocks"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestArticlesCmd(t *testing.T) {
	if articlesCmd == nil {
		t.Errorf("articlesCmd is not initialized")
	}
}

func TestArticlesCmdFlags(t *testing.T) {
	// Set the flags
	if err := articlesCmd.Flags().Set("crawlsiteid", "test"); err != nil {
		t.Fatalf("Error setting crawlsiteid flag: %v", err)
	}
	if err := articlesCmd.Flags().Set("searchterms", "test"); err != nil {
		t.Fatalf("Error setting searchterms flag: %v", err)
	}
	if err := articlesCmd.Flags().Set("url", "test"); err != nil {
		t.Fatalf("Error setting url flag: %v", err)
	}

	// Check if the flags are correctly set
	assert.Equal(t, "test", viper.GetString("crawlsiteid"))
	assert.Equal(t, "test", viper.GetString("searchterms"))
	assert.Equal(t, "test", viper.GetString("url"))
}

func TestSaveResultsToRedis(t *testing.T) {
	ctx := context.Background()
	mockRedisClient := mocks.NewMockRedisClient() // Create the MockRedisClient using the NewMockRedisClient function
	manager := &crawler.CrawlManager{
		Client: mockRedisClient,
		Logger: &mocks.MockLogger{}, // replace with your actual mock logger
	}
	server := &CrawlServer{
		CrawlManager: manager,
	}

	results := []crawler.PageData{
		{
			URL:           "https://www.jonesrussell42.xyz/privacy-policy",
			MatchingTerms: []string{"PRIVACI", "POLICI"},
		},
	}
	key := "testKey"

	err := server.saveResultsToRedis(ctx, results, key)
	assert.NoError(t, err)

	// Add assertions to check if the results were correctly saved to Redis
	savedResults, err := mockRedisClient.SMembers(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, len(results), len(savedResults))
	for i, result := range results {
		var savedResult crawler.PageData
		err := json.Unmarshal([]byte(savedResults[i]), &savedResult)
		assert.NoError(t, err)
		assert.Equal(t, result, savedResult)
	}
}
