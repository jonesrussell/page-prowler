package cmd_test

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/jonesrussell/page-prowler/cmd"
	"github.com/jonesrussell/page-prowler/cmd/mocks"
	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestClearlinksCmd_WithValidCrawlsiteid(t *testing.T) {
	// Initialize the CrawlManager
	manager := &crawler.CrawlManager{
		Client:         prowlredis.NewMockClient().(*prowlredis.MockClient),
		Logger:         mocks.NewMockLogger(),
		MongoDBWrapper: mocks.NewMockMongoDBWrapper(),
	}

	ctx := context.WithValue(context.Background(), common.CrawlManagerKey, manager)

	// Set crawlsiteid to a non-empty string
	crawlsiteid := "testsite"

	// Set the crawlsiteid environment variable
	os.Setenv("CRAWLSITEID", crawlsiteid)

	// Set the crawlsiteid flag in the viper configuration
	viper.Set("crawlsiteid", crawlsiteid)

	// Create a new instance of ClearlinksCmd
	clearlinksCmd := &cobra.Command{
		Use:   "clearlinks",
		Short: "Clear the Redis set for a given crawlsiteid",
		RunE:  cmd.ClearlinksCmd.RunE,
	}

	// Set the context in the command
	clearlinksCmd.SetContext(ctx)

	// Execute the command and check the error
	err := clearlinksCmd.Execute()
	assert.NoError(t, err)

	// Check if Del was called on the Redis client
	assert.True(t, manager.Client.(*prowlredis.MockClient).WasDelCalled, "Del should have been called on the Redis client")

	// Unset the crawlsiteid environment variable
	os.Unsetenv("CRAWLSITEID")
}

func TestClearlinksCmd_WithEmptyCrawlsiteid(t *testing.T) {
	// Unset the 'crawlsiteid' environment variable
	viper.Reset()
	os.Unsetenv("CRAWLSITEID")

	// Initialize the CrawlManager
	manager := &crawler.CrawlManager{
		Client:         prowlredis.NewMockClient().(*prowlredis.MockClient),
		Logger:         mocks.NewMockLogger(),
		MongoDBWrapper: mocks.NewMockMongoDBWrapper(),
	}

	ctx := context.WithValue(context.Background(), common.CrawlManagerKey, manager)

	// Initialize the command
	clearlinksCmd := &cobra.Command{
		Use:   "clearlinks",
		Short: "Clear the Redis set for a given crawlsiteid",
		RunE:  cmd.ClearlinksCmd.RunE,
	}

	// Set the context in the command
	clearlinksCmd.SetContext(ctx)

	// Set the 'crawlsiteid' flag
	clearlinksCmd.Flags().StringP("crawlsiteid", "s", "", "CrawlSite ID")
	viper.BindPFlag("crawlsiteid", clearlinksCmd.Flags().Lookup("crawlsiteid"))

	// Enable viper to read from environment variables
	viper.AutomaticEnv()

	// Execute the command
	err := clearlinksCmd.Execute()

	// Check for error
	if err == nil || err.Error() != cmd.ErrCrawlsiteidRequired.Error() {
		t.Errorf("clearlinksCmd.Execute() error = %v, wantErr %v", err, cmd.ErrCrawlsiteidRequired)
	}
	// Reset the 'crawlsiteid' environment variable
	os.Setenv("CRAWLSITEID", "")
}

func TestClearlinksCmd_WhenRedisClientReturnsError(t *testing.T) {
	// Initialize the CrawlManager with a mock Redis client that returns an error
	mockRedisClient := prowlredis.NewMockClient().(*prowlredis.MockClient)
	mockRedisClient.DelErr = errors.New("mock Redis error")
	manager := &crawler.CrawlManager{
		Client:         mockRedisClient,
		Logger:         mocks.NewMockLogger(),
		MongoDBWrapper: mocks.NewMockMongoDBWrapper(),
	}

	ctx := context.WithValue(context.Background(), common.CrawlManagerKey, manager)

	// Create a new Cobra command for testing
	clearlinksCmd := &cobra.Command{
		Use:   cmd.ClearlinksCmd.Use,
		Short: cmd.ClearlinksCmd.Short,
		Long:  cmd.ClearlinksCmd.Long,
		RunE:  cmd.ClearlinksCmd.RunE,
	}

	// Set the 'crawlsiteid' flag
	clearlinksCmd.Flags().StringP("crawlsiteid", "s", "", "CrawlSite ID")
	viper.BindPFlag("crawlsiteid", clearlinksCmd.Flags().Lookup("crawlsiteid"))

	// Enable viper to read from environment variables
	viper.AutomaticEnv()

	// Set the 'crawlsiteid' flag value
	viper.Set("crawlsiteid", "test_crawlsiteid")

	// Discard the command's output
	clearlinksCmd.SetOut(io.Discard)

	// Execute the command and check the error
	err := clearlinksCmd.ExecuteContext(ctx)
	assert.Error(t, err)
	assert.Equal(t, "failed to clear Redis set: mock Redis error", err.Error())
}
