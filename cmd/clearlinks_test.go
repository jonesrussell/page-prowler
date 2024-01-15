package cmd

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/jonesrussell/page-prowler/cmd/mocks"
	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/stretchr/testify/assert"
)

func initializeTestManager() (*crawler.CrawlManager, error) {
	return initializeManager(
		prowlredis.NewMockClient().(*prowlredis.MockClient),
		mocks.NewMockLogger(),
		mocks.NewMockMongoDBWrapper(),
	)
}

func TestClearlinksCmd_WithoutInitializedManager(t *testing.T) {
	// Set Crawlsiteid to a non-empty string
	Crawlsiteid = "testsite"

	// Do not initialize CrawlManager
	ctx := context.WithValue(context.Background(), common.CrawlManagerKey, nil)

	// Create a new Cobra command for testing
	cmd := CreateTestCommand(clearlinksCmd.Use, clearlinksCmd.Short, clearlinksCmd.Long, clearlinksCmd.RunE)

	// Discard the command's output
	cmd.SetOut(io.Discard)

	// Execute the command and check the error
	err := cmd.ExecuteContext(ctx)
	assert.Error(t, err)
	assert.Equal(t, ErrCrawlManagerNotInitialized, err)
}

func TestClearlinksCmd_WithValidCrawlsiteid(t *testing.T) {
	// Save the original value of Crawlsiteid and reset it after the test
	originalCrawlsiteid := Crawlsiteid
	defer func() { Crawlsiteid = originalCrawlsiteid }()

	// Set Crawlsiteid to a non-empty string
	Crawlsiteid = "testsite"

	// Initialize the CrawlManager
	manager, err := initializeTestManager()
	if err != nil {
		t.Fatalf("Failed to initialize CrawlManager: %v", err)
	}

	ctx := context.WithValue(context.Background(), common.CrawlManagerKey, manager)

	// Create a new Cobra command for testing
	cmd := CreateTestCommand(clearlinksCmd.Use, clearlinksCmd.Short, clearlinksCmd.Long, clearlinksCmd.RunE)

	// Execute the command and check the error
	err = cmd.ExecuteContext(ctx)
	assert.NoError(t, err)
	assert.True(t, manager.Client.(*prowlredis.MockClient).WasDelCalled, "Del should have been called on the Redis client")
}

func TestClearlinksCmd_WithEmptyCrawlsiteid(t *testing.T) {
	// Save the original value of Crawlsiteid and reset it after the test
	originalCrawlsiteid := Crawlsiteid
	defer func() { Crawlsiteid = originalCrawlsiteid }()

	// Set Crawlsiteid to a non-empty string
	Crawlsiteid = ""

	// Initialize the CrawlManager
	manager, err := initializeTestManager()
	if err != nil {
		t.Fatalf("Failed to initialize CrawlManager: %v", err)
	}

	ctx := context.WithValue(context.Background(), common.CrawlManagerKey, manager)

	// Create a new Cobra command for testing
	cmd := CreateTestCommand(clearlinksCmd.Use, clearlinksCmd.Short, clearlinksCmd.Long, clearlinksCmd.RunE)

	// Discard the command's output
	cmd.SetOut(io.Discard)

	// Execute the command and check the error
	err = cmd.ExecuteContext(ctx)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrCrawlsiteidRequired)
}

func TestClearlinksCmd_WhenRedisClientReturnsError(t *testing.T) {
	// Save the original value of Crawlsiteid and reset it after the test
	originalCrawlsiteid := Crawlsiteid
	defer func() { Crawlsiteid = originalCrawlsiteid }()

	// Set Crawlsiteid to a non-empty string
	Crawlsiteid = "testsite"

	// Initialize the CrawlManager with a mock Redis client that returns an error
	mockRedisClient := prowlredis.NewMockClient().(*prowlredis.MockClient)
	mockRedisClient.DelErr = errors.New("mock Redis error")
	manager, err := initializeTestManager()
	if err != nil {
		t.Fatalf("Failed to initialize CrawlManager: %v", err)
	}
	manager.Client = mockRedisClient

	ctx := context.WithValue(context.Background(), common.CrawlManagerKey, manager)

	// Create a new Cobra command for testing
	cmd := CreateTestCommand(clearlinksCmd.Use, clearlinksCmd.Short, clearlinksCmd.Long, clearlinksCmd.RunE)

	// Discard the command's output
	cmd.SetOut(io.Discard)

	// Execute the command and check the error
	err = cmd.ExecuteContext(ctx)
	assert.Error(t, err)
	assert.Equal(t, "failed to clear Redis set: mock Redis error", err.Error())
}

func TestClearlinksCmd_CheckLoggingOutput(t *testing.T) {
	// Save the original value of Crawlsiteid and reset it after the test
	originalCrawlsiteid := Crawlsiteid
	defer func() { Crawlsiteid = originalCrawlsiteid }()

	// Set Crawlsiteid to a non-empty string
	Crawlsiteid = "testsite"

	// Initialize the CrawlManager
	manager, err := initializeTestManager()
	if err != nil {
		t.Fatalf("Failed to initialize CrawlManager: %v", err)
	}

	ctx := context.WithValue(context.Background(), common.CrawlManagerKey, manager)

	// Create a new Cobra command for testing
	cmd := CreateTestCommand(clearlinksCmd.Use, clearlinksCmd.Short, clearlinksCmd.Long, clearlinksCmd.RunE)

	// Discard the command's output
	cmd.SetOut(io.Discard)

	// Execute the command and check the error
	err = cmd.ExecuteContext(ctx)
	assert.NoError(t, err)

	// Check the log messages
	entries := manager.Logger.(*mocks.MockLogger).AllEntries()
	assert.Len(t, entries, 1)
	assert.Equal(t, "Redis set cleared successfully", entries[0].Message)
}
