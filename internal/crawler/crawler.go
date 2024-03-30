package crawler

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/stats"
)

const (
	// DefaultParallelism defines the default number of concurrent operations allowed during crawling.
	// It is set to 2 to balance between performance and resource usage.
	DefaultParallelism = 2

	// DefaultDelay specifies the default delay between consecutive crawling operations in milliseconds.
	// It is set to 3000 milliseconds (3 seconds) to avoid overwhelming the target server with requests.
	DefaultDelay = 3000 * time.Millisecond
)

// CrawlManagerInterface defines the interface for managing crawling operations.
// It includes methods for setting up crawling logic, handling errors, and starting the crawling process.
//
//go:generate mockery --name=CrawlManagerInterface
type CrawlManagerInterface interface {
	Crawl(ctx context.Context, url string, maxDepth int, debug bool) ([]PageData, error)
	// Search initiates the crawling process for a given URL with the provided options.
	// It returns a slice of PageData and an error if any occurs during the crawling process.
	Search(ctx context.Context, url string, searchTerms, crawlSiteID string, maxDepth int, debug bool) ([]PageData, error)
	// SetupHTMLParsingHandler sets up the handler for HTML parsing with gocolly, using the provided parameters.
	// It returns an error if the setup fails.
	SetupHTMLParsingHandler(handler func(*colly.HTMLElement)) error
	// SetupErrorEventHandler sets up the HTTP request error handling for the colly collector.
	// It configures the collector to handle different types of errors.
	SetupErrorEventHandler(collector *colly.Collector)
	// SetupCrawlingLogic configures and initiates the crawling logic.
	// It sets up the HTML parsing handler and error event handler for the collector.
	// It returns an error if the setup fails.
	SetupCrawlingLogic(*CrawlOptions) error
	// CrawlURL visits the given URL and performs the crawling operation.
	// It logs the visit and waits for the collector to finish its tasks.
	// It returns an error if the visit fails.
	CrawlURL(url string) error
	// HandleVisitError handles the error occurred during the visit of a URL.
	// It logs the error and returns it.
	HandleVisitError(url string, err error) error
	// Logger returns the logger instance associated with the CrawlManager.
	Logger() logger.Logger
	// StartCrawling initiates the crawling process with the given parameters.
	// It validates the input parameters, configures the collector, and starts the crawling process.
	// It returns an error if the crawling process fails to start.
	ProcessMatchingLink(options *CrawlOptions, currentURL string, pageData PageData, matchingTerms []string)
	UpdateStats(options *CrawlOptions, matchingTerms []string)
	Collector(*colly.Collector) *CollectorWrapper
}

// CrawlManager is the implementation of the CrawlManagerInterface.
// It manages the crawling operations, including setting up crawling logic, handling errors, and starting the crawling process.
// The struct fields are initialized with default values or instances of required types.
var _ CrawlManagerInterface = &CrawlManager{
	LoggerField:       nil,                                     // Logger instance for logging messages.
	Client:            nil,                                     // HTTP client for making requests.
	MongoDBWrapper:    nil,                                     // MongoDB wrapper for database operations.
	CollectorInstance: &CollectorWrapper{colly.NewCollector()}, // Colly collector for crawling web pages.
	CrawlingMu:        &sync.Mutex{},                           // Mutex for synchronizing crawling operations.
	StatsManager:      &StatsManager{},                         // Manager for crawling statistics.
}

// Logger returns the logger instance associated with the CrawlManager.
// It provides access to the logging functionality for the crawling operations.
func (cm *CrawlManager) Logger() logger.Logger {
	return cm.LoggerField
}

// validateParameters checks if the provided parameters for crawling are valid.
// It returns an error if any of the parameters are invalid (e.g., empty strings or non-positive maxDepth).
// Parameters:
// - url: The URL to crawl.
// - searchTerms: The search terms to match against the crawled content.
// - crawlSiteID: The ID of the site to crawl.
// - maxDepth: The maximum depth to crawl.
func (cm *CrawlManager) validateParameters(url, searchTerms, crawlSiteID string, maxDepth int) error {
	if url == "" || searchTerms == "" || crawlSiteID == "" || maxDepth <= 0 {
		return errors.New("invalid parameters")
	}
	return nil
}

// initializeStatsManager initializes the StatsManager with default values.
// It sets up the LinkStats and LinkStatsMu fields of the StatsManager.
// The method also locks the CrawlingMu mutex to ensure thread safety during the initialization process.
func (cm *CrawlManager) initializeStatsManager() {
	cm.StatsManager = &StatsManager{
		LinkStats:   &stats.Stats{},
		LinkStatsMu: sync.RWMutex{},
	}
	cm.CrawlingMu.Lock()
	defer cm.CrawlingMu.Unlock()
}

// SetupHTMLParsingHandler sets up the handler for HTML parsing with gocolly, using the provided parameters.
// It configures the collector to handle HTML elements matching the "a[href]" selector by invoking the provided handler function.
// Parameters:
// - handler: A function that takes a *colly.HTMLElement as an argument and performs actions on the element.
// Returns:
// - error: An error if the setup fails.
func (cm *CrawlManager) SetupHTMLParsingHandler(handler func(*colly.HTMLElement)) error {
	cm.CollectorInstance.OnHTML("a[href]", handler)
	return nil
}

// SetupErrorEventHandler sets up the HTTP request error handling for the colly collector.
// It configures the collector to handle different types of errors, specifically logging 500 Internal Server Errors without printing the stack trace, and logging other errors normally.
// Parameters:
// - collector: A pointer to the colly.Collector instance for which the error handling is being set up.
func (cm *CrawlManager) SetupErrorEventHandler(collector *colly.Collector) {
	cm.CollectorInstance.OnError(func(r *colly.Response, err error) {
		statusCode := r.StatusCode
		requestURL := r.Request.URL.String()

		if statusCode == 500 {
			// Handle 500 Internal Server Error without printing the stack trace
			cm.LoggerField.Debug(fmt.Sprintf("[SetupErrorEventHandler] Internal Server Error request_url: %s, status_code: %d", requestURL, statusCode))
		} else if statusCode != 404 {
			// Handle other errors normally
			cm.LoggerField.Debug(fmt.Sprintf("[SetupErrorEventHandler] Request URL failed request_url: %s, status_code: %d", requestURL, statusCode))
		}
	})
}

// SetupCrawlingLogic configures and initiates the crawling logic.
// It sets up the HTML parsing handler and error event handler for the collector.
// It returns an error if the setup fails.
// Parameters:
// - options: The CrawlOptions containing configuration for the crawling process.
// Returns:
// - error: An error if the setup fails.
func (cm *CrawlManager) SetupCrawlingLogic(options *CrawlOptions) error {
	err := cm.SetupHTMLParsingHandler(cm.GetAnchorElementHandler(options))
	if err != nil {
		return cm.handleSetupError(err)
	}

	cm.SetupErrorEventHandler(&colly.Collector{})

	return nil
}

// CrawlURL visits the given URL and performs the crawling operation.
// It logs the visit, waits for the collector to finish its tasks, and logs the completion of the crawling process.
// Parameters:
// - url: The URL to visit and crawl.
// Returns:
// - error: An error if the visit fails or if an error occurs during the crawling process.
func (cm *CrawlManager) CrawlURL(url string) error {
	cm.LoggerField.Debug(fmt.Sprintf("[CrawlURL] Visiting URL: %v", url))

	err := cm.visitWithColly(url)
	if err != nil {
		return cm.HandleVisitError(url, err)
	}

	cm.CollectorInstance.Wait()

	cm.Logger().Info("[CrawlURL] Crawling completed.")
	return nil
}
