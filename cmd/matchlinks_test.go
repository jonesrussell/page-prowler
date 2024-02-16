package cmd_test

import (
	"context"
	"errors"
	"testing"

	crawlerMocks "github.com/jonesrussell/page-prowler/mocks"
)

/*func TestMatchlinksCmd(t *testing.T) {
	mockCrawlManager := &mocks.MockCrawlManager{}
	ctx := context.WithValue(context.Background(), common.CrawlManagerKey, &MockCrawlManager{})

	cmd := CreateTestCommand(matchlinksCmd.Use, matchlinksCmd.Short, matchlinksCmd.Long, matchlinksCmd.RunE)

	// Define the flags
	cmd.Flags().StringP("siteid", "s", "", "CrawlSite ID")
	cmd.Flags().StringP("url", "u", "", "URL to crawl")
	cmd.Flags().StringP("searchterms", "t", "", "Search terms for crawling")
	cmd.Flags().IntP("maxdepth", "m", 1, "Max depth for crawling")

	cmd.SetArgs([]string{"--siteid", "test", "--url", "http://example.com", "--searchterms", "test", "--maxdepth", "1"})

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
	mockManager := new(crawlerMocks.CrawlManagerInterface)
	ctx := context.Background()
	url := "https://example.com"
	searchterms := "test"
	siteid := "123"
	maxdepth := 1

	mockManager.On("StartCrawling", ctx, url, searchterms, siteid, maxdepth, false).Return(nil)

	// Call the function that uses StartCrawling here
	err := mockManager.StartCrawling(ctx, url, searchterms, siteid, maxdepth, false)
	if err != nil {
		t.Errorf("error running matchlinks command: %v", err)
	}

	mockManager.AssertExpectations(t)
}

func TestStartCrawlingHandlesErrors(t *testing.T) {
	mockManager := new(crawlerMocks.CrawlManagerInterface)
	ctx := context.Background()
	url := "https://example.com"
	searchterms := "test"
	siteid := "123"
	maxdepth := 1

	// Set up the mock to return an error
	mockManager.On("StartCrawling", ctx, url, searchterms, siteid, maxdepth, false).Return(errors.New("mock error"))

	// Call the function that uses StartCrawling here
	err := mockManager.StartCrawling(ctx, url, searchterms, siteid, maxdepth, false)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	mockManager.AssertExpectations(t)
}

/*func TestFlagsAreBoundCorrectly(t *testing.T) {
	// Define a flag
	flagName := "siteid"
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
