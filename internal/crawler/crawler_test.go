package crawler

import (
	"context"
	"testing"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"go.uber.org/zap"
)

func TestHandleHTMLParsing(t *testing.T) {
	// Create a mock CrawlManager
	cs := &CrawlManager{}

	// Create a mock CrawlOptions with a mock Collector
	options := &CrawlOptions{
		Collector: colly.NewCollector(),
	}

	// Call the function with the mock parameters
	err := cs.setupHTMLParsingHandler(context.Background(), options)

	// Check that no error was returned
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// TODO: Add more checks here to verify that the OnHTML method was called correctly
}

func TestHandleErrorEvents(t *testing.T) {
	// Create a mock CrawlManager
	cs := &CrawlManager{}

	// Create a mock Collector
	collector := colly.NewCollector()

	// Call the function with the mock parameters
	cs.setupErrorEventHandler(collector)

	// Trigger an error in the collector
	collector.OnError(func(r *colly.Response, err error) {
		// Check that the error handling function was called with the correct parameters
		if r.StatusCode != 404 {
			t.Errorf("Expected status code to be 404, but got %d", r.StatusCode)
		}
	})
}

func TestSetupCrawlingLogic(t *testing.T) {
	// Create a mock Logger
	zapLogger, _ := zap.NewDevelopment()
	logger := &logger.ZapLoggerWrapper{Logger: zapLogger.Sugar()}

	// Create a mock CrawlManager with the Logger
	cs := &CrawlManager{
		Logger: logger,
	}

	// Create a mock CrawlOptions with a mock Collector
	options := &CrawlOptions{
		Collector: colly.NewCollector(),
	}

	// Call the function with the mock parameters
	err := cs.setupCrawlingLogic(context.Background(), options)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// TODO: Add more checks here to verify that the OnHTML, OnError, OnScraped, and OnRequest methods were called correctly
}
