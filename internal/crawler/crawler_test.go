package crawler

import (
	"context"
	"reflect"
	"testing"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"go.uber.org/zap"
)

func TestConfigureCollector(t *testing.T) {
	allowedDomains := []string{"test.com"}
	maxDepth := 2

	collector := ConfigureCollector(allowedDomains, maxDepth)

	if collector == nil {
		t.Errorf("Expected collector to be not nil")
		return
	}

	if !reflect.DeepEqual(collector.AllowedDomains, allowedDomains) {
		t.Errorf("Expected allowed domains to be %v, but got %v", allowedDomains, collector.AllowedDomains)
	}

	if collector.MaxDepth != maxDepth {
		t.Errorf("Expected max depth to be %d, but got %d", maxDepth, collector.MaxDepth)
	}
}

func TestHandleHTMLParsing(t *testing.T) {
	// Create a mock CrawlManager
	cs := &CrawlManager{}

	// Create a mock CrawlOptions with a mock Collector
	options := &CrawlOptions{
		Collector: colly.NewCollector(),
	}

	// Call the function with the mock parameters
	err := cs.HandleHTMLParsing(context.Background(), options)

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
	cs.handleErrorEvents(collector)

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
	cs.SetupCrawlingLogic(context.Background(), options)

	// TODO: Add more checks here to verify that the OnHTML, OnError, OnScraped, and OnRequest methods were called correctly
}

func TestGetHostFromURL(t *testing.T) {
	// Create a mock Logger
	zapLogger, _ := zap.NewDevelopment()
	logger := &logger.ZapLoggerWrapper{Logger: zapLogger.Sugar()}

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
		host, err := GetHostFromURL(tc.url, logger)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if host != tc.expected {
			t.Errorf("Expected host %v, but got %v", tc.expected, host)
		}
	}
}
