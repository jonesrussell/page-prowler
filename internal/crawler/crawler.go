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
	// Crawl initiates the crawling process for a given URL with the provided options.
	// It returns a slice of PageData and an error if any occurs during the crawling process.
	Crawl(url string, options *CrawlOptions) ([]PageData, error)
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
	StartCrawling(ctx context.Context, url string, searchterms string, siteid string, maxdepth int, debug bool) error
	ProcessMatchingLinkAndUpdateStats(*CrawlOptions, string, PageData, []string)
}

// CrawlManager is the implementation of the CrawlManagerInterface.
// It manages the crawling operations, including setting up crawling logic, handling errors, and starting the crawling process.
// The struct fields are initialized with default values or instances of required types.
var _ CrawlManagerInterface = &CrawlManager{
	LoggerField:    nil,                // Logger instance for logging messages.
	Client:         nil,                // HTTP client for making requests.
	MongoDBWrapper: nil,                // MongoDB wrapper for database operations.
	Collector:      &colly.Collector{}, // Colly collector for crawling web pages.
	CrawlingMu:     &sync.Mutex{},      // Mutex for synchronizing crawling operations.
	StatsManager:   &StatsManager{},    // Manager for crawling statistics.
}

// CrawlOptions represents the options for configuring and initiating the crawling logic.
// It includes parameters that control the crawling process, such as the site ID to crawl, search terms to match, a pointer to store results, and a debug flag.
type CrawlOptions struct {
	CrawlSiteID string      // The ID of the site to crawl.
	SearchTerms []string    // The search terms to match against the crawled content.
	Results     *[]PageData // A pointer to a slice of PageData where the crawling results will be stored.
	Debug       bool        // A flag indicating whether to enable debug mode for the crawling process.
}

// Logger returns the logger instance associated with the CrawlManager.
// It provides access to the logging functionality for the crawling operations.
func (cm *CrawlManager) Logger() logger.Logger {
	return cm.LoggerField
}

// StartCrawling initiates the crawling process with the given parameters.
// It validates the input parameters, configures the collector, and starts the crawling process.
// It returns an error if the crawling process fails to start.
// Parameters:
// - ctx: The context for the crawling operation.
// - url: The URL to start crawling from.
// - searchTerms: The search terms to match against the crawled content.
// - crawlSiteID: The ID of the site to crawl.
// - maxDepth: The maximum depth to crawl.
// - debug: A flag indicating whether to enable debug mode for the crawling process.
func (cm *CrawlManager) StartCrawling(ctx context.Context, url, searchTerms, crawlSiteID string, maxDepth int, debug bool) error {
	if err := cm.validateParameters(url, searchTerms, crawlSiteID, maxDepth); err != nil {
		return err
	}

	cm.initializeStatsManager()

	host, err := cm.extractHostFromURL(url)
	if err != nil {
		return err
	}

	if err := cm.configureCollector(host, maxDepth); err != nil {
		return err
	}

	options := cm.createCrawlingOptions(crawlSiteID, searchTerms, debug)

	return cm.performCrawling(ctx, url, options)
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

// extractHostFromURL extracts the host from the given URL.
// It uses the GetHostFromURL function to parse the URL and retrieve the host.
// If the URL cannot be parsed, it logs an error and returns an empty string along with the error.
// Parameters:
// - url: The URL from which to extract the host.
// Returns:
// - string: The extracted host from the URL.
// - error: An error if the URL cannot be parsed.
func (cm *CrawlManager) extractHostFromURL(url string) (string, error) {
	host, err := GetHostFromURL(url, cm.Logger())
	if err != nil {
		cm.LoggerField.Error(fmt.Sprintf("Failed to parse URL: url: %v, error: %v", url, err))
		return "", err
	}
	cm.LoggerField.Debug(fmt.Sprintf("Extracted host from URL: %s", host))
	return host, nil
}

// configureCollector configures the crawler's collector with the specified host and maximum depth.
// It attempts to set up the collector for crawling operations.
// If the configuration fails, it logs a fatal error and returns the error.
// Parameters:
// - host: The host to configure the collector for.
// - maxDepth: The maximum depth for the crawling operations.
// Returns:
// - error: An error if the collector configuration fails.
func (cm *CrawlManager) configureCollector(host string, maxDepth int) error {
	err := cm.ConfigureCollector([]string{host}, maxDepth)
	if err != nil {
		cm.LoggerField.Fatal(fmt.Sprintf("Failed to configure collector: %v", err))
		return err
	}
	return nil
}

// createCrawlingOptions creates and returns a CrawlOptions instance based on the provided parameters.
// It splits the search terms and then calls createStartCrawlingOptions to create the CrawlOptions.
// Parameters:
// - crawlSiteID: The ID of the site to crawl.
// - searchTerms: The search terms to match against the crawled content.
// - debug: A flag indicating whether to enable debug mode for the crawling process.
// Returns:
// - *CrawlOptions: A pointer to a CrawlOptions instance configured with the provided parameters.
func (cm *CrawlManager) createCrawlingOptions(crawlSiteID, searchTerms string, debug bool) *CrawlOptions {
	splitSearchTerms := cm.splitSearchTerms(searchTerms)
	return cm.createStartCrawlingOptions(crawlSiteID, splitSearchTerms, debug)
}

// performCrawling executes the crawling process for the specified URL with the given options.
// It performs the crawl, logs crawling statistics, saves the results to Redis, and logs the results.
// Parameters:
// - ctx: The context for the crawling operation.
// - url: The URL to crawl.
// - options: The CrawlOptions containing configuration for the crawling process.
// Returns:
// - error: An error if the crawling process encounters any issues.
func (cm *CrawlManager) performCrawling(ctx context.Context, url string, options *CrawlOptions) error {
	results, err := cm.Crawl(url, options)
	if err != nil {
		return err
	}

	cm.logCrawlingStatistics()

	if err := cm.SaveResultsToRedis(ctx, results, options.CrawlSiteID); err != nil {
		return err
	}

	logResults(cm, results)

	return nil
}

// Crawl starts the crawling process for a given URL with the provided options.
// It logs the URL being crawled, sets up the crawling logic, visits the URL, and returns the crawling results.
// Parameters:
// - url: The URL to start crawling.
// - options: The CrawlOptions containing configuration for the crawling process.
// Returns:
// - []PageData: A slice of PageData representing the crawling results.
// - error: An error if the crawling process encounters any issues.
func (cm *CrawlManager) Crawl(url string, options *CrawlOptions) ([]PageData, error) {
	cm.LoggerField.Debug(fmt.Sprintf("CrawlURL: %s", url))
	err := cm.SetupCrawlingLogic(options)
	if err != nil {
		return nil, err
	}

	err = cm.CrawlURL(url)
	if err != nil {
		return nil, err
	}

	return *options.Results, nil
}

// SetupHTMLParsingHandler sets up the handler for HTML parsing with gocolly, using the provided parameters.
// It configures the collector to handle HTML elements matching the "a[href]" selector by invoking the provided handler function.
// Parameters:
// - handler: A function that takes a *colly.HTMLElement as an argument and performs actions on the element.
// Returns:
// - error: An error if the setup fails.
func (cm *CrawlManager) SetupHTMLParsingHandler(handler func(*colly.HTMLElement)) error {
	cm.Collector.OnHTML("a[href]", handler)
	return nil
}

// SetupErrorEventHandler sets up the HTTP request error handling for the colly collector.
// It configures the collector to handle different types of errors, specifically logging 500 Internal Server Errors without printing the stack trace, and logging other errors normally.
// Parameters:
// - collector: A pointer to the colly.Collector instance for which the error handling is being set up.
func (cm *CrawlManager) SetupErrorEventHandler(collector *colly.Collector) {
	collector.OnError(func(r *colly.Response, err error) {
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

	cm.SetupErrorEventHandler(cm.Collector)

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
	//	cm.trackVisitedPage(url, options)
	cm.Collector.Wait()
	cm.Logger().Info("[CrawlURL] Crawling completed.")
	return nil
}

// HandleVisitError handles the error occurred during the visit of a URL.
// It logs the error and returns it.
// Parameters:
// - url: The URL that encountered an error during the visit.
// - err: The error that occurred during the visit.
// Returns:
// - error: The error that was logged and returned.
func (cm *CrawlManager) HandleVisitError(url string, err error) error {
	cm.LoggerField.Error(fmt.Sprintf("Error visiting URL: url: %s, error: %v", url, err))
	return err
}
