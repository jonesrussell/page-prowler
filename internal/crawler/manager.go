package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/jonesrussell/page-prowler/internal/stats"
)

// CrawlManager is the main struct that manages the crawling operations.
// It includes fields for logging, Redis operations, and the Colly collector.
type CrawlManager struct {
	Client            prowlredis.ClientInterface
	CollectorInstance *CollectorWrapper
	CrawlingMu        *sync.Mutex
	LoggerField       logger.Logger
	Options           *CrawlOptions
	StatsManager      *StatsManager
	Results           *Results
}

func (cm *CrawlManager) GetStatsManager() *StatsManager {
	return cm.StatsManager
}

func (cm *CrawlManager) GetCollector() *CollectorWrapper {
	return cm.CollectorInstance
}

// NewCrawlManager creates a new instance of CrawlManager with the provided logger,
// Redis client. It initializes the CrawlingMu mutex.
func NewCrawlManager(
	loggerField logger.Logger,
	client prowlredis.ClientInterface,
	options *CrawlOptions,
) *CrawlManager {
	return &CrawlManager{
		LoggerField:       loggerField,
		Client:            client,
		CollectorInstance: NewCollectorWrapper(colly.NewCollector()), // Initialize the CollectorInstance field
		CrawlingMu:        &sync.Mutex{},
		Options:           options, // Store the provided CrawlOptions
		Results:           NewResults(),
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

func (cm *CrawlManager) Crawl() error {
	startURL := cm.GetOptions().StartURL

	cm.LoggerField.Debug(fmt.Sprintf("[Crawl] Starting crawl for URL: %s", startURL))

	cm.initializeStatsManager()

	host, err := cm.extractHostFromURL(startURL)
	if err != nil {
		return err
	}

	if err := cm.ConfigureCollector([]string{host}, cm.GetOptions().MaxDepth); err != nil {
		return err
	}

	if err := cm.SetupCrawlingLogic(); err != nil {
		return err
	}

	if err := cm.visitWithColly(startURL); err != nil {
		return cm.HandleVisitError(startURL, err)
	}

	cm.CollectorInstance.Wait()

	cm.Logger().Info("[Crawl] Crawling completed.")

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

	if err := cm.CollectorInstance.Limit(); err != nil {
		cm.LoggerField.Error(fmt.Sprintf("Failed to set limit rule: %v", err))
		return err
	}

	// Respect robots.txt
	cm.CollectorInstance.SetAllowURLRevisit(false)
	cm.CollectorInstance.SetIgnoreRobotsTxt(false)

	// Register OnScraped callback
	//cm.CollectorInstance.OnScraped(func(r *colly.Response) {
	//	cm.LoggerField.Debug(fmt.Sprintf("[OnScraped] Page scraped: %s", r.Request.URL.String()))
	//	cm.StatsManager.LinkStatsMu.Lock()
	//	defer cm.StatsManager.LinkStatsMu.Unlock()
	//	cm.StatsManager.LinkStats.IncrementTotalPages()
	//})

	return nil
}

func (cm *CrawlManager) visitWithColly(url string) error {
	// Visit the URL with the Colly collector
	err := cm.CollectorInstance.Visit(url)
	if err != nil {
		return err
	}

	// Wait for the collector to finish its tasks
	cm.CollectorInstance.Wait()

	return nil
}

// AppendResult appends a PageData to the Results.
func (cm *CrawlManager) AppendResult(pageData PageData) {
	if cm.Results == nil || cm.Results.Pages == nil {
		fmt.Println("Warning: Results or Pages is nil")
		return
	}
	cm.Results.Pages = append(cm.Results.Pages, pageData)
}

// GetResults retrieves the Results managed by this CrawlManager.
func (cm *CrawlManager) GetResults() *Results {
	return cm.Results
}

func (cm *CrawlManager) SaveResultsToRedis(ctx context.Context, results []PageData, key string) error {
	cm.LoggerField.Debug(fmt.Sprintf("SaveResultsToRedis: Number of results before processing: %d", len(results)))

	for _, result := range results {
		cm.LoggerField.Debug(fmt.Sprintf("SaveResultsToRedis: Processing result %v", result))

		data, err := json.Marshal(result)
		if err != nil {
			cm.LoggerField.Error(fmt.Sprintf("SaveResultsToRedis: Error occurred during marshalling to JSON: %v", err))
			return err
		}
		str := string(data)
		err = cm.Client.SAdd(ctx, key, str)
		if err != nil {
			cm.LoggerField.Error(fmt.Sprintf("SaveResultsToRedis: Error occurred during saving to Redis: %v", err))
			return err
		}
		cm.LoggerField.Debug("SaveResultsToRedis: Added elements to the set")

		// Debugging: Verify that the result was saved correctly
		isMember, err := cm.Client.SIsMember(ctx, key, str)
		if err != nil {
			cm.LoggerField.Error(fmt.Sprintf("SaveResultsToRedis: Error occurred during checking membership in Redis set: %v", err))
			return err
		}
		if !isMember {
			cm.LoggerField.Error(fmt.Sprintf("SaveResultsToRedis: Result was not saved correctly in Redis set: %v", str))
		} else {
			cm.LoggerField.Debug(fmt.Sprintf("SaveResultsToRedis: Result was saved correctly in Redis set, key: %s, result: %s", key, str))
		}
	}

	cm.LoggerField.Debug(fmt.Sprintf("SaveResultsToRedis: Number of results after processing: %d", len(results)))

	return nil
}
