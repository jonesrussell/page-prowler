package crawler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/stats"
	"github.com/jonesrussell/page-prowler/mocks"
	"go.uber.org/zap"
)

func MockServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Write your mock response here
		_, err := rw.Write([]byte(`Hello, World!`))
		if err != nil {
			return
		}
	}))

	return server
}

func setupTestEnvironment() (*crawler.CrawlManager, *crawler.CrawlOptions) {
	// Create a mock Logger
	zapLogger, _ := zap.NewDevelopment()
	log := &logger.ZapLoggerWrapper{Logger: zapLogger.Sugar()}

	// Create a mock CrawlManager with the Logger
	cs := crawler.NewCrawlManager(
		log,
		mocks.NewMockClient(),
		mocks.NewMockMongoDBWrapper(),
	)

	// Initialize the Collector
	cs.Collector = colly.NewCollector()

	// Initialize StatsManager with non-nil LinkStats
	statsManager := &crawler.StatsManager{
		LinkStats: &stats.Stats{},
	}

	cs.StatsManager = statsManager

	// Initialize CrawlingMu
	cs.CrawlingMu = &sync.Mutex{}

	// Create a mock CrawlOptions with initialized LinkStats
	options := crawler.NewCrawlOptions(
		"crawlSiteID",
		[]string{"term"},
		true,                  // Debug
		&[]crawler.PageData{}, // Initialize Results
	)

	return cs, options
}

func TestCrawler_StartCrawling(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		searchTerms string
		crawlSiteID string
		maxDepth    int
		debug       bool
		wantErr     bool
	}{
		{
			name:        "Valid parameters",
			url:         "https://example.com",
			searchTerms: "term",
			crawlSiteID: "crawlSiteID",
			maxDepth:    1,
			debug:       true,
			wantErr:     false,
		},
		{
			name:        "Invalid parameters",
			url:         "",
			searchTerms: "",
			crawlSiteID: "",
			maxDepth:    0,
			debug:       false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs, _ := setupTestEnvironment()

			// Create a mock server
			server := MockServer()
			defer server.Close()

			// Call the function with the test parameters
			err := cs.StartCrawling(context.Background(), server.URL, tt.searchTerms, tt.crawlSiteID, tt.maxDepth, tt.debug)

			if (err != nil) != tt.wantErr {
				t.Errorf("StartCrawling() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
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
	myurl := server.URL + `/your-endpoint`

	// Call the function with the mock parameters
	_, err := cs.Crawl(myurl, options)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// TODO: Add more checks here to verify that the crawl function behaves as expected
}

func TestGetAnchorElementHandler(t *testing.T) {
	cs, options := setupTestEnvironment()

	// Call the function with the mock parameters
	handler := cs.GetAnchorElementHandler(options)

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
}

func TestSetupHTMLParsingHandler(t *testing.T) {
	cs, options := setupTestEnvironment()

	// Create a mock server
	server := MockServer()
	defer server.Close()

	// Call the function with the mock parameters
	err := cs.SetupHTMLParsingHandler(cs.GetAnchorElementHandler(options))
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
