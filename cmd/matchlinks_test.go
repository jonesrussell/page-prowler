package cmd

import (
	"context"
	"errors"
	"testing"

	"github.com/jonesrussell/page-prowler/cmd/mocks"
)

/*func TestMatchlinksCmd(t *testing.T) {
	mockCrawlManager := &mocks.MockCrawlManager{}
	ctx := context.WithValue(context.Background(), common.CrawlManagerKey, &MockCrawlManager{})

	cmd := CreateTestCommand(matchlinksCmd.Use, matchlinksCmd.Short, matchlinksCmd.Long, matchlinksCmd.RunE)

	// Define the flags
	cmd.Flags().StringP("crawlsiteid", "s", "", "CrawlSite ID")
	cmd.Flags().StringP("url", "u", "", "URL to crawl")
	cmd.Flags().StringP("searchterms", "t", "", "Search terms for crawling")
	cmd.Flags().IntP("maxdepth", "m", 1, "Max depth for crawling")

	cmd.SetArgs([]string{"--crawlsiteid", "test", "--url", "http://example.com", "--searchterms", "test", "--maxdepth", "1"})

	// Create a buffer to capture the output
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	// Execute the command
	err := cmd.ExecuteContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Add your assertions here based on the expected output
	mockCrawlManager.AssertExpectations(t)
}*/

func TestStartCrawlingIsCalledWithCorrectArguments(t *testing.T) {
	mockManager := new(mocks.MockCrawlManager)
	ctx := context.Background()
	url := "http://example.com"
	searchterms := "test"
	crawlsiteid := "123"
	maxdepth := 1
	debug := false

	mockManager.On("StartCrawling", ctx, url, searchterms, crawlsiteid, maxdepth, debug).Return(nil)

	// Call the function that uses StartCrawling here
	err := mockManager.StartCrawling(ctx, url, searchterms, crawlsiteid, maxdepth, debug)
	if err != nil {
		t.Errorf("error running matchlinks command: %v", err)
	}

	mockManager.AssertExpectations(t)
}

func TestStartCrawlingHandlesErrors(t *testing.T) {
	mockManager := new(mocks.MockCrawlManager)
	ctx := context.Background()
	url := "http://example.com"
	searchterms := "test"
	crawlsiteid := "123"
	maxdepth := 1
	debug := false

	// Set up the mock to return an error
	mockManager.On("StartCrawling", ctx, url, searchterms, crawlsiteid, maxdepth, debug).Return(errors.New("mock error"))

	// Call the function that uses StartCrawling here
	err := mockManager.StartCrawling(ctx, url, searchterms, crawlsiteid, maxdepth, debug)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	mockManager.AssertExpectations(t)
}

/*func TestFlagsAreBoundCorrectly(t *testing.T) {
	// Define a flag
	flagName := "crawlsiteid"
	flagValue := "testvalue"

	// Set the flag value
	matchlinksCmd.SetArgs([]string{"--" + flagName, flagValue})

	// Parse the flags
	if err := matchlinksCmd.ParseFlags([]string{"--" + flagName, flagValue}); err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	// Check if the Viper key has the correct value
	if viper.GetString(flagName) != flagValue {
		t.Errorf("Viper key does not have the correct value. Expected %s, got %s", flagValue, viper.GetString(flagName))
	}
}*/
