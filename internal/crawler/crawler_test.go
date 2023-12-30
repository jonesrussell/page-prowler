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

func TestGetHostFromURL(t *testing.T) {
	// Create a mock Logger
	zapLogger, _ := zap.NewDevelopment()
	log := &logger.ZapLoggerWrapper{Logger: zapLogger.Sugar()}

	// Define test cases
	testCases := []struct {
		url      string
		expected string
	}{
		{"http://example.com/path", "example.com"},
		{"https://www.example.com", "www.example.com"},
		// add more test cases here
	}

	// Run test cases
	for _, tc := range testCases {
		host, err := GetHostFromURL(tc.url, log)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if host != tc.expected {
			t.Errorf("Expected host %v, but got %v", tc.expected, host)
		}
	}
}

func TestSetupCrawlingLogic(t *testing.T) {
	// Create a mock Logger
	zapLogger, _ := zap.NewDevelopment()
	log := &logger.ZapLoggerWrapper{Logger: zapLogger.Sugar()}

	// Create a mock CrawlManager with the Logger
	cs := &CrawlManager{
		Logger: log,
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
