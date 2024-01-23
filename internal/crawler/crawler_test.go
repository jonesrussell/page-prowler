package crawler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/cmd/mocks"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/jonesrussell/page-prowler/internal/stats"
	"go.uber.org/zap"
)

func MockServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Write your mock response here
		rw.Write([]byte(`Hello, World!`))
	}))

	return server
}

func setupTestEnvironment() (*CrawlManager, *CrawlOptions) {
	// Create a mock Logger
	zapLogger, _ := zap.NewDevelopment()
	log := &logger.ZapLoggerWrapper{Logger: zapLogger.Sugar()}

	// Create a mock CrawlManager with the Logger
	cs := NewCrawlManager(
		log,
		prowlredis.NewMockClient(),
		mocks.NewMockMongoDBWrapper(),
	)

	// Create a mock CrawlOptions with initialized LinkStats
	options := NewCrawlOptions(
		"crawlSiteID",
		[]string{"term"},
		true,             // Debug
		&[]PageData{},    // Initialize Results
		stats.NewStats(), // Initialize LinkStats
	)

	return cs, options
}

func TestHandleHTMLParsing(t *testing.T) {
	cs, options := setupTestEnvironment()

	// Call the function with the mock parameters
	err := cs.SetupHTMLParsingHandler(cs.getAnchorElementHandler(options))

	// Check that no error was returned
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// TODO: Add more checks here to verify that the OnHTML method was called correctly
}

func TestHandleErrorEvents(t *testing.T) {
	cs, _ := setupTestEnvironment()

	// Call the function with the mock parameters
	cs.SetupErrorEventHandler(cs.Collector)

	// Trigger an error in the collector
	cs.Collector.OnError(func(r *colly.Response, err error) {
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
		{"https://example.com/path", "example.com"},
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

func TestVisitURL(t *testing.T) {
	cs, options := setupTestEnvironment()

	// Call the function with the mock parameters
	cs.CrawlURL("https://example.com", options)

	// Check that the VisitedPages map has been populated correctly
	if _, ok := cs.VisitedPages["https://example.com"]; !ok {
		t.Errorf("Expected https://example.com to be in VisitedPages, but it wasn't")
	}
}

func TestConfigureCollector(t *testing.T) {
	cs, _ := setupTestEnvironment()

	// Call the function with the mock parameters
	err := cs.ConfigureCollector([]string{"example.com"}, 1)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// TODO: Add more checks here to verify that the Collector is configured correctly
}

func TestStartCrawling(t *testing.T) {
	cs, _ := setupTestEnvironment()

	// Call the function with the mock parameters
	err := cs.StartCrawling(context.Background(), "https://example.com", "term", "crawlSiteID", 1, true)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// TODO: Add more checks here to verify that the crawling process is started correctly
}

func TestSetupCrawlingLogic(t *testing.T) {
	cs, options := setupTestEnvironment()

	// Call the function with the mock parameters
	err := cs.SetupCrawlingLogic(options)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// TODO: Add more checks here to verify that the OnHTML, OnError, OnScraped, and OnRequest methods were called correctly
}

func TestCrawl(t *testing.T) {
	cs, options := setupTestEnvironment()

	// Create a mock server
	server := MockServer()
	defer server.Close()

	// Replace the URL with the mock server URL
	url := server.URL + "/your-endpoint"

	// Call the function with the mock parameters
	_, err := cs.Crawl(url, options)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// TODO: Add more checks here to verify that the crawl function behaves as expected
}

func TestGetAnchorElementHandler(t *testing.T) {
	cs, options := setupTestEnvironment()

	// Call the function with the mock parameters
	handler := cs.getAnchorElementHandler(options)

	// Create a mock HTMLElement
	var u url.URL
	u.Scheme = "https"
	u.Host = "example.com"
	u.Path = "/path"
	element := &colly.HTMLElement{
		Request: &colly.Request{
			URL: &u,
		},
		Text: "Link Text",
	}

	// Call the handler function with the mock element
	handler(element)

	// Check that the IncrementTotalLinks method was called correctly
	if options.LinkStats.TotalLinks != 1 {
		t.Errorf("Expected TotalLinks to be 1, but got %d", options.LinkStats.TotalLinks)
	}
}

func TestSetupHTMLParsingHandler(t *testing.T) {
	cs, options := setupTestEnvironment()

	// Create a mock server
	server := MockServer()
	defer server.Close()

	// Call the function with the mock parameters
	err := cs.SetupHTMLParsingHandler(cs.getAnchorElementHandler(options))
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// Call the crawl function with the mock server URL
	_, err = cs.Crawl(server.URL, options)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// TODO: Add more checks here to verify that the OnHTML method was called correctly
}

func TestSetupErrorEventHandler(t *testing.T) {
	cs, _ := setupTestEnvironment()

	// Call the function with the mock parameters
	cs.SetupErrorEventHandler(cs.Collector)

	// TODO: Add more checks here to verify that the OnError method was called correctly
}
