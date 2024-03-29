package crawler_test

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/stats"
	"github.com/jonesrussell/page-prowler/mocks"
	"go.uber.org/mock/gomock"
)

func MockServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Log the request URL
		log.Printf("Mock server received request for URL: %s", req.URL.String())

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
	log := mocks.NewMockLogger()

	// Create a mock CrawlManager with the Logger
	cs := crawler.NewCrawlManager(
		log,
		mocks.NewMockClient(),
		mocks.NewMockMongoDBWrapper(),
	)

	// Initialize the Collector
	collector := colly.NewCollector()
	cs.Collector = crawler.NewCollectorWrapper(collector)

	cs.Collector.SetAllowedDomains([]string{"https://example.com"})

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
			_, err := cs.Crawl(context.Background(), tt.url, tt.searchTerms, tt.crawlSiteID, tt.maxDepth, tt.debug)

			if (err != nil) != tt.wantErr {
				t.Errorf("StartCrawling() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCrawlManager_SetupCrawlingLogic(t *testing.T) {
	cm, options := setupTestEnvironment()

	// Create a gomock Controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock collector to check if the methods are set up
	mockCollector := mocks.NewMockCollectorInterface(ctrl)
	cm.Collector = mockCollector

	// Set up the expectation for the OnHTML method
	mockCollector.EXPECT().OnHTML("a[href]", gomock.Any()).Return()

	// Set up the expectation for the OnError method
	mockCollector.EXPECT().OnError(gomock.Any()).Return()

	// Call the SetupCrawlingLogic method
	err := cm.SetupCrawlingLogic(options)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestCrawlManager_Crawl(t *testing.T) {
	cs, _ := setupTestEnvironment()

	cs.Collector.SetAllowedDomains([]string{"example.com"})

	// Create a mock server
	server := MockServer()
	defer server.Close()

	// Replace the URL with the mock server URL
	myurl := server.URL + `/your-endpoint`

	u, _ := url.Parse(server.URL)
	cs.Collector.SetAllowedDomains([]string{u.Host})

	// Call the function with the mock parameters
	ctx := context.Background()      // Create a context
	searchTerms := "yourSearchTerms" // Example search terms
	crawlSiteID := "yourCrawlSiteID" // Example crawl site ID
	maxDepth := 1                    // Example max depth
	debug := false                   // Example debug flag

	_, err := cs.Crawl(ctx, myurl, searchTerms, crawlSiteID, maxDepth, debug)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// TODO: Add more checks here to verify that the crawl function behaves as expected
}

func TestCrawlManager_CrawlURL(t *testing.T) {
	// Setup: Create a mock CrawlManager and a mock server
	// Call the CrawlURL method with the mock server URL
	// Assert the expected behavior, such as logging the visit
}

func TestCrawlManager_HandleVisitError(t *testing.T) {
	// Setup: Create a mock CrawlManager
	// Call the HandleVisitError method with a mock URL and error
	// Assert that the error is handled correctly, such as logging the error
}

func TestCrawlManager_Logger(t *testing.T) {
	// Setup: Create a mock CrawlManager with a mock Logger
	// Call the Logger method
	// Assert that the correct Logger instance is returned
}

func TestCrawlManager_ProcessMatchingLink(t *testing.T) {
	// Setup: Create a mock CrawlManager and CrawlOptions
	// Call the ProcessMatchingLink method with test parameters
	// Assert the expected behavior, such as processing matching links
}

func TestCrawlManager_UpdateStats(t *testing.T) {
	// Setup: Create a mock CrawlManager and CrawlOptions
	// Call the UpdateStats method with test parameters
	// Assert the expected behavior, such as updating stats
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

func TestCrawlManager_SetupHTMLParsingHandler(t *testing.T) {
	cs, options := setupTestEnvironment()

	// Create a mock server
	server := MockServer()
	defer server.Close()

	u, _ := url.Parse(server.URL)
	cs.Collector.SetAllowedDomains([]string{u.Host})

	// Call the function with the mock parameters
	err := cs.SetupHTMLParsingHandler(cs.GetAnchorElementHandler(options))
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// Call the crawl function with the mock server URL
	ctx := context.Background()      // Create a context
	searchTerms := "yourSearchTerms" // Example search terms
	crawlSiteID := "yourCrawlSiteID" // Example crawl site ID
	maxDepth := 1                    // Example max depth
	debug := false                   // Example debug flag

	_, err = cs.Crawl(ctx, server.URL, searchTerms, crawlSiteID, maxDepth, debug)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// TODO: Add more checks here to verify that the OnHTML method was called correctly
}

func TestCrawlManager_SetupErrorEventHandler(t *testing.T) {
	cs, _ := setupTestEnvironment()

	// Define test cases
	tests := []struct {
		name          string
		simulateError bool
		expectedError error
	}{
		{
			name:          "No Error",
			simulateError: false,
			expectedError: nil,
		},
		{
			name:          "Simulated Error",
			simulateError: true,
			expectedError: errors.New("simulated error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs.Collector.SetAllowedDomains([]string{"example.com"})

			// Simulate an error if needed
			if tt.simulateError {
				// Assuming you have a way to simulate errors, such as a mock HTTP client
				// that can return errors for specific requests.
				// Here, you would set up your mock client to return tt.expectedError for the next request.
				// This part depends on your mock client's implementation.
				// For example, if your mock client has a method to set the next error:
				// mockClient.SetNextError(tt.expectedError)
			}

			// Set up a mock logger to capture log messages
			mockLogger := mocks.NewMockLogger()
			cs.LoggerField = mockLogger

			// Call the function with the mock parameters
			cs.SetupErrorEventHandler(cs.Collector.GetUnderlyingCollector())

			// Trigger a request that would cause the OnError method to be called
			// This could be done by visiting a URL with the collector
			err := cs.Collector.Visit("http://example.com")
			if err != nil && tt.simulateError {
				// Verify that the OnError method was called correctly
				// Assuming your mock logger has a method to assert if a specific log message was recorded
				// mockLogger.AssertCalled(t, "Error", tt.expectedError)
			} else if err == nil && !tt.simulateError {
				// If no error is expected, ensure no error log was recorded
				// mockLogger.AssertNotCalled(t, "Error")
			} else {
				t.Fatalf("Unexpected error: %v", err)
			}
		})
	}
}
