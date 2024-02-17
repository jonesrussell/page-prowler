package cmd_test

import (
	"context"
	"errors"
	"testing"

	crawlerMocks "github.com/jonesrussell/page-prowler/mocks"
)

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
