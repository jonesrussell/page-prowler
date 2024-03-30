package crawler

import (
	"context"
	"fmt"
	"sync"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/mongodbwrapper"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/jonesrussell/page-prowler/internal/stats"
)

// CrawlManager is the main struct that manages the crawling operations.
// It includes fields for logging, MongoDB operations, and the Colly collector.
type CrawlManager struct {
	LoggerField       logger.Logger
	Client            prowlredis.ClientInterface
	MongoDBWrapper    mongodbwrapper.MongoDBInterface
	CollectorInstance *CollectorWrapper
	CrawlingMu        *sync.Mutex
	StatsManager      *StatsManager
}

func (cm *CrawlManager) GetStatsManager() *StatsManager {
	return cm.StatsManager
}

func (cm *CrawlManager) GetCollector() *CollectorWrapper {
	return cm.CollectorInstance
}

// NewCrawlManager creates a new instance of CrawlManager with the provided logger,
// Redis client, and MongoDB wrapper. It initializes the CrawlingMu mutex.
func NewCrawlManager(
	loggerField logger.Logger,
	client prowlredis.ClientInterface,
	mongoDBWrapper mongodbwrapper.MongoDBInterface,
) *CrawlManager {
	return &CrawlManager{
		LoggerField:       loggerField,
		Client:            client,
		MongoDBWrapper:    mongoDBWrapper,
		CollectorInstance: NewCollectorWrapper(colly.NewCollector()), // Initialize the CollectorInstance field
		CrawlingMu:        &sync.Mutex{},
	}
}

// StatsManager is a struct that manages crawling statistics.
// It includes fields for link statistics and a mutex for thread safety.
type StatsManager struct {
	LinkStats   *stats.Stats
	LinkStatsMu sync.RWMutex
}

// NewStatsManager creates a new StatsManager with initialized fields.
func NewStatsManager() *StatsManager {
	return &StatsManager{
		LinkStats:   &stats.Stats{},
		LinkStatsMu: sync.RWMutex{},
	}
}

// Add a Collector method to the CrawlManager struct.
func (cm *CrawlManager) Collector(c *colly.Collector) *CollectorWrapper {
	return cm.CollectorInstance
}

// Adjust the variable declaration to use the Collector method.
var _ CrawlManagerInterface = &CrawlManager{
	LoggerField:       nil,
	Client:            nil,
	MongoDBWrapper:    nil,
	CollectorInstance: NewCollectorWrapper(colly.NewCollector()),
	CrawlingMu:        &sync.Mutex{},
}

func (cm *CrawlManager) Crawl(ctx context.Context, url string, maxDepth int, debug bool) ([]PageData, error) {
	cm.LoggerField.Debug(fmt.Sprintf("[Crawl] Starting crawl for URL: %s", url))

	/*if err := cm.validateParameters(url, maxDepth); err != nil {
		return nil, err
	}*/

	cm.initializeStatsManager()

	host, err := cm.extractHostFromURL(url)
	if err != nil {
		return nil, err
	}

	if err := cm.ConfigureCollector([]string{host}, maxDepth); err != nil {
		return nil, err
	}

	var results []PageData
	options := NewCrawlOptions(debug, &results)
	if err := cm.SetupCrawlingLogic(options); err != nil {
		return nil, err
	}

	if err := cm.visitWithColly(url); err != nil {
		return nil, cm.HandleVisitError(url, err)
	}

	cm.CollectorInstance.Wait()
	cm.Logger().Info("[Crawl] Crawling completed.")
	return *options.Results, nil
}

func (cm *CrawlManager) Search(ctx context.Context, url string, searchTerms, crawlSiteID string, maxDepth int, debug bool) ([]PageData, error) {
	cm.LoggerField.Debug(fmt.Sprintf("[Crawl] Starting crawl for URL: %s", url))

	if err := cm.validateParameters(url, searchTerms, crawlSiteID, maxDepth); err != nil {
		return nil, err
	}

	cm.initializeStatsManager()

	host, err := cm.extractHostFromURL(url)
	if err != nil {
		return nil, err
	}

	if err := cm.ConfigureCollector([]string{host}, maxDepth); err != nil {
		return nil, err
	}

	var results []PageData
	options := NewCrawlOptions(debug, &results)
	if err := cm.SetupCrawlingLogic(options); err != nil {
		return nil, err
	}

	if err := cm.visitWithColly(url); err != nil {
		return nil, cm.HandleVisitError(url, err)
	}

	cm.CollectorInstance.Wait()
	cm.Logger().Info("[Crawl] Crawling completed.")
	return *options.Results, nil
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

// ConfigureCollector sets up the Colly collector with the specified allowed domains and maximum depth.
// It also configures the collector to log debug information, respect robots.txt, and register an OnScraped callback.
// Parameters:
// - allowedDomains: A slice of strings representing the allowed domains for crawling.
// - maxDepth: The maximum depth to crawl.
// Returns:
// - error: An error if the collector configuration fails.
func (cm *CrawlManager) ConfigureCollector(allowedDomains []string, maxDepth int) error {
	cm.CollectorInstance = &CollectorWrapper{
		colly.NewCollector(
			colly.Async(false),
			colly.MaxDepth(maxDepth),
			colly.Debugger(cm.LoggerField),
		),
	}

	cm.LoggerField.Debug(fmt.Sprintf("Allowed Domains: %v", allowedDomains))
	cm.CollectorInstance.SetAllowedDomains(allowedDomains)

	limitRule := cm.createLimitRule()
	if err := cm.CollectorInstance.Limit(limitRule); err != nil {
		cm.LoggerField.Error(fmt.Sprintf("Failed to set limit rule: %v", err))
		return err
	}

	// Respect robots.txt
	cm.CollectorInstance.SetAllowURLRevisit(false)
	cm.CollectorInstance.SetIgnoreRobotsTxt(false)

	// Register OnScraped callback
	cm.CollectorInstance.OnScraped(func(r *colly.Response) {
		cm.LoggerField.Debug(fmt.Sprintf("[OnScraped] Page scraped: %s", r.Request.URL.String()))
		cm.StatsManager.LinkStatsMu.Lock()
		defer cm.StatsManager.LinkStatsMu.Unlock()
		cm.StatsManager.LinkStats.IncrementTotalPages()
	})

	return nil
}

func (cm *CrawlManager) visitWithColly(url string) error {
	// Assuming you have a method to set up the Colly collector
	var results []PageData
	options := NewCrawlOptions(false, &results)
	err := cm.SetupCrawlingLogic(options)
	if err != nil {
		return err
	}

	// Visit the URL with the Colly collector
	err = cm.CollectorInstance.Visit(url)
	if err != nil {
		return err
	}

	return nil
}
