package crawler_test

import (
	"sync"
	"testing"

	"github.com/jonesrussell/page-prowler/cmd/mocks"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	crawlerMocks "github.com/jonesrussell/page-prowler/internal/crawler/mocks"
	"github.com/jonesrussell/page-prowler/internal/stats"
	"github.com/stretchr/testify/mock"
)

func TestProcessMatchingLinkAndUpdateStats(t *testing.T) {
	// Define some test data
	href := "http://example.com"
	pageData := crawler.PageData{}
	matchingTerms := []string{"term1", "term2"}

	// Create an instance of the mock processor
	mockProcessor := new(crawlerMocks.MatchedLinkProcessor)

	// Define the behavior of the mock methods as needed
	mockProcessor.On("IncrementMatchedLinks", mock.AnythingOfType("*crawler.CrawlOptions")).Return()
	mockProcessor.On("HandleMatchingLinks", href).Return(nil)
	mockProcessor.On("UpdatePageData", &pageData, href, matchingTerms).Return()
	mockProcessor.On("AppendResult", mock.AnythingOfType("*crawler.CrawlOptions"), mock.AnythingOfType("crawler.PageData")).Return()

	// Create an instance of CrawlManager and set the MatchedLinkProcessor to the mock processor
	crawlManager := &crawler.CrawlManager{
		MatchedLinkProcessor: mockProcessor,
		LoggerField:          mocks.NewMockLogger(),
	}

	// Call the function with the test data
	options := &crawler.CrawlOptions{
		CrawlSiteID: "testSite",
		SearchTerms: []string{"term1", "term2"},
		Results:     &[]crawler.PageData{},
		LinkStats:   stats.NewStats(),
		LinkStatsMu: sync.Mutex{},
		Debug:       false,
	}
	crawlManager.ProcessMatchingLinkAndUpdateStats(options, href, pageData, matchingTerms)

	// Assert that the mock methods were called with the correct parameters
	mockProcessor.AssertExpectations(t)
}

/*func equalPageData(a, b PageData) bool {
	return a.URL == b.URL &&
		a.ParentURL == b.ParentURL &&
		reflect.DeepEqual(a.SearchTerms, b.SearchTerms) &&
		reflect.DeepEqual(a.MatchingTerms, b.MatchingTerms) &&
		a.Error == b.Error
}

func TestGetHref(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Write a small HTML page with a single link
		rw.Write([]byte(`<a href="/test">Test</a>`))
	}))
	defer server.Close()

	// Create a new CrawlManager
	cm := &CrawlManager{
		Logger: mocks.NewMockLogger(),
	}

	// Create a new Collector
	c := colly.NewCollector()

	var href string
	c.OnHTML("a", func(e *colly.HTMLElement) {
		href = cm.getHref(e)
	})

	// Visit the mock server's URL
	err := c.Visit(server.URL)
	if err != nil {
		t.Fatalf("Failed to visit server: %v", err)
	}

	// Check the result
	expectedHref := server.URL + "/test"
	if href != expectedHref {
		t.Errorf("Expected '%s', got '%s'", expectedHref, href)
	}

	// Additional assertions
	if !strings.HasPrefix(href, "http") {
		t.Errorf("Expected href to start with 'http', got '%s'", href)
	}

	if !strings.HasSuffix(href, "/test") {
		t.Errorf("Expected href to end with '/test', got '%s'", href)
	}
}

func TestIncrementMatchedLinks(t *testing.T) {
	// Create a new CrawlManager
	cm := &CrawlManager{}

	// Create a new CrawlOptions with a new Stats object
	options := &CrawlOptions{
		LinkStats:   stats.NewStats(),
		LinkStatsMu: sync.Mutex{},
	}

	// Call the incrementMatchedLinks function
	cm.incrementMatchedLinks(options)

	// Check the result
	if options.LinkStats.GetMatchedLinks() != 1 {
		t.Errorf("Expected matched links count to be 1, got %d", options.LinkStats.GetMatchedLinks())
	}
}

func TestUpdatePageData(t *testing.T) {
	// Create a new CrawlManager
	cm := &CrawlManager{}

	// Define a test PageData
	pageData := &PageData{}

	// Define test inputs
	href := "https://example.com"
	matchingTerms := []string{"term1", "term2"}

	// Call the updatePageData function
	cm.updatePageData(pageData, href, matchingTerms)

	// Check the result
	if !reflect.DeepEqual(pageData.MatchingTerms, matchingTerms) {
		t.Errorf("Expected MatchingTerms to be '%v', got '%v'", matchingTerms, pageData.MatchingTerms)
	}

	if pageData.ParentURL != href {
		t.Errorf("Expected ParentURL to be '%s', got '%s'", href, pageData.ParentURL)
	}
}

func TestAppendResult(t *testing.T) {
	// Create a new CrawlManager
	cm := &CrawlManager{}

	// Create a new CrawlOptions with an empty Results slice
	options := &CrawlOptions{
		Results:     &[]PageData{},
		LinkStatsMu: sync.Mutex{},
	}

	// Define a test PageData
	pageData := PageData{
		// Initialize the fields of PageData as needed
	}

	// Call the appendResult function
	cm.appendResult(options, pageData)

	// Check the result
	if len(*options.Results) != 1 {
		t.Errorf("Expected length of Results to be 1, got %d", len(*options.Results))
	}

	if !equalPageData((*options.Results)[0], pageData) {
		t.Errorf("Expected first element of Results to be '%v', got '%v'", pageData, (*options.Results)[0])
	}
}

func TestIncrementNonMatchedLinkCount(t *testing.T) {
	// Create a new CrawlManager
	cm := &CrawlManager{
		Logger: mocks.NewMockLogger(),
	}

	// Create a new CrawlOptions with a new Stats object
	options := &CrawlOptions{
		LinkStats:   stats.NewStats(),
		LinkStatsMu: sync.Mutex{},
	}

	// Call the incrementNonMatchedLinkCount function
	cm.incrementNonMatchedLinkCount(options)

	// Check the result
	if options.LinkStats.GetNotMatchedLinks() != 1 {
		t.Errorf("Expected not matched links count to be 1, got %d", options.LinkStats.NotMatchedLinks)
	}
}

func TestCreateLimitRule(t *testing.T) {
	// Create a new CrawlManager
	cm := &CrawlManager{}

	// Call the createLimitRule function
	rule := cm.createLimitRule()

	// Check the result
	if rule.DomainGlob != "*" {
		t.Errorf("Expected DomainGlob to be '*', got '%s'", rule.DomainGlob)
	}

	if rule.Parallelism != DefaultParallelism {
		t.Errorf("Expected Parallelism to be %d, got %d", DefaultParallelism, rule.Parallelism)
	}

	if rule.Delay != DefaultDelay {
		t.Errorf("Expected Delay to be %v, got %v", DefaultDelay, rule.Delay)
	}
}

func TestSplitSearchTerms(t *testing.T) {
	cm := &CrawlManager{}

	// Test with multiple terms
	terms := cm.splitSearchTerms("term1,term2,term3")
	expectedTerms := []string{"term1", "term2", "term3"}
	if !reflect.DeepEqual(terms, expectedTerms) {
		t.Errorf("Expected '%v', got '%v'", expectedTerms, terms)
	}

	// Test with a single term
	terms = cm.splitSearchTerms("term1")
	expectedTerms = []string{"term1"}
	if !reflect.DeepEqual(terms, expectedTerms) {
		t.Errorf("Expected '%v', got '%v'", expectedTerms, terms)
	}

	// Test with invalid terms
	terms = cm.splitSearchTerms("term1,,term3")
	expectedTerms = []string{"term1", "term3"} // Empty terms are ignored
	if !reflect.DeepEqual(terms, expectedTerms) {
		t.Errorf("Expected '%v', got '%v'", expectedTerms, terms)
	}
}

func TestCreateStartCrawlingOptions(t *testing.T) {
	// Create a new CrawlManager
	cm := &CrawlManager{}

	// Define test inputs
	crawlSiteID := "testSite"
	searchTerms := []string{"term1", "term2"}
	debug := true

	// Call the createStartCrawlingOptions function
	options := cm.createStartCrawlingOptions(crawlSiteID, searchTerms, debug)

	// Check the result
	if options.CrawlSiteID != crawlSiteID {
		t.Errorf("Expected CrawlSiteID to be '%s', got '%s'", crawlSiteID, options.CrawlSiteID)
	}

	if !reflect.DeepEqual(options.SearchTerms, searchTerms) {
		t.Errorf("Expected SearchTerms to be '%v', got '%v'", searchTerms, options.SearchTerms)
	}

	if options.Debug != debug {
		t.Errorf("Expected Debug to be '%v', got '%v'", debug, options.Debug)
	}

	if options.Results == nil {
		t.Error("Expected Results to be initialized, got nil")
	}

	if options.LinkStats == nil {
		t.Error("Expected LinkStats to be initialized, got nil")
	}
}

func TestCrawlManager_ProcessMatchingLinkAndUpdateStats(t *testing.T) {
	type fields struct {
		Logger               logger.Logger
		Client               prowlredis.ClientInterface
		MongoDBWrapper       mongodbwrapper.MongoDBInterface
		Collector            *colly.Collector
		CrawlingMu           *sync.Mutex
		VisitedPages         map[string]bool
		MatchedLinkProcessor MatchedLinkProcessor
	}

	type args struct {
		options       *CrawlOptions
		href          string
		pageData      PageData
		matchingTerms []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test Case 1: Empty href",
			fields: fields{
				Logger:               mocks.NewMockLogger(),
				Client:               prowlredis.NewMockClient(),
				MongoDBWrapper:       mocks.NewMockMongoDBWrapper(),
				Collector:            colly.NewCollector(),
				CrawlingMu:           &sync.Mutex{},
				VisitedPages:         make(map[string]bool),
				MatchedLinkProcessor: mocks.NewMockMatchedLinkProcessor(),
			},
			args: args{
				options:       &CrawlOptions{ },
				href:          "",
				pageData:      PageData{  },
				matchingTerms: []string{"term1", "term2"},
			},
			wantErr: true,
		},
		{
			name: "Test Case 2: Normal operation",
			fields: fields{
				Logger:               mocks.NewMockLogger(),
				Client:               prowlredis.NewMockClient(),
				MongoDBWrapper:       mocks.NewMockMongoDBWrapper(),
				Collector:            colly.NewCollector(),
				CrawlingMu:           &sync.Mutex{},
				VisitedPages:         make(map[string]bool),
				MatchedLinkProcessor: mocks.NewMockMatchedLinkProcessor(),
			},
			args: args{
				options:       &CrawlOptions{  },
				href:          "https://example.com",
				pageData:      PageData{  },
				matchingTerms: []string{"term1", "term2"},
			},
			wantErr: false,
		},
		// Add more test cases as needed...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &CrawlManager{
				Logger:               tt.fields.Logger,
				Client:               tt.fields.Client,
				MongoDBWrapper:       tt.fields.MongoDBWrapper,
				Collector:            tt.fields.Collector,
				CrawlingMu:           tt.fields.CrawlingMu,
				VisitedPages:         tt.fields.VisitedPages,
				MatchedLinkProcessor: tt.fields.MatchedLinkProcessor,
			}
			cs.ProcessMatchingLinkAndUpdateStats(tt.args.options, tt.args.href, tt.args.pageData, tt.args.matchingTerms)
			// Add assertions here to check for expected behavior
		})
	}
}*/
