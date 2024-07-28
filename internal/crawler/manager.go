package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/jonesrussell/page-prowler/internal/stats"
)

type CrawlManagerInterface interface {
	Crawl() error
	SetupErrorEventHandler()
	CrawlURL(url string) error
	HandleVisitError(url string, err error) error
	Logger() loggo.LoggerInterface
	ProcessMatchingLink(currentURL string, pageData PageData, matchingTerms []string)
	UpdateStats(options *CrawlOptions, matchingTerms []string)
	SetOptions(options *CrawlOptions) error
	Client() prowlredis.ClientInterface
}

// CrawlManager is the main struct that manages the crawling operations.
// It includes fields for logging, MongoDB operations, and the Colly collector.
type CrawlManager struct {
	client            prowlredis.ClientInterface
	CollectorInstance *CollectorWrapper
	CrawlingMu        *sync.Mutex
	LoggerField       loggo.LoggerInterface
	Options           *CrawlOptions
	StatsManager      *StatsManager
	Results           *Results
}

var _ CrawlManagerInterface = &CrawlManager{}

func (cm *CrawlManager) Client() prowlredis.ClientInterface {
	return cm.client
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
	loggerField loggo.LoggerInterface,
	client prowlredis.ClientInterface,
	collectorInstance *CollectorWrapper,
	options *CrawlOptions,
) *CrawlManager {
	return &CrawlManager{
		LoggerField:       loggerField,
		client:            client,
		CollectorInstance: collectorInstance,
		CrawlingMu:        &sync.Mutex{},
		Options:           options,
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
	cm.Logger().Info("[Crawl] Starting Crawl function")

	// Get and print options
	options := cm.GetOptions()
	cm.Logger().Debug(fmt.Sprintf("[Crawl] Options: %+v", options))

	startURL := options.StartURL

	cm.Logger().Debug(fmt.Sprintf("[Crawl] Starting crawl for URL: %s", startURL))

	cm.initializeStatsManager()

	host, err := cm.extractHostFromURL(startURL)
	if err != nil {
		return err
	}

	if err := cm.ConfigureCollector([]string{host}, options.MaxDepth); err != nil {
		return err
	}

	start := time.Now()
	if err := cm.visitWithColly(startURL); err != nil {
		elapsed := time.Since(start)
		cm.Logger().Error(fmt.Sprintf("Error visiting URL: %s, elapsed time: %s", startURL, elapsed), err)
		return cm.HandleVisitError(startURL, err)
	}
	elapsed := time.Since(start)
	cm.Logger().Info(fmt.Sprintf("Visited URL: %s, elapsed time: %s", startURL, elapsed))

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
	cm.Logger().Error(fmt.Sprintf("Error visiting URL: url: %s", url), err, nil) // Add nil as the third argument
	return err
}

func (cm *CrawlManager) ConfigureCollector(allowedDomains []string, maxDepth int) error {
	collector := colly.NewCollector(
		colly.Async(false),
		colly.MaxDepth(maxDepth),
	)

	cm.Logger().Debug(fmt.Sprintf("Allowed Domains: %v", allowedDomains))
	collector.AllowedDomains = allowedDomains

	if err := collector.Limit(&colly.LimitRule{}); err != nil {
		cm.Logger().Error(fmt.Sprintf("Failed to set limit rule: %v", err), nil)
		return err
	}

	// Respect robots.txt
	collector.AllowURLRevisit = false
	collector.IgnoreRobotsTxt = false

	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := cm.getHref(e)
		if href == "" {
			return // Return to indicate no error occurred
		}

		cm.processLink(e, href)
		err := cm.visitWithColly(href)
		if err != nil {
			cm.LoggerField.Debug(fmt.Sprintf("[GetAnchorElementHandler] Error visiting URL: %s, Error: %v", href, err))
			// Handle the error here
		}
	})

	cm.CollectorInstance = &CollectorWrapper{collector}

	return nil
}

func (cm *CrawlManager) visitWithColly(url string) error {
	// Visit the URL with the Colly collector
	start := time.Now()
	err := cm.CollectorInstance.Visit(url)
	elapsed := time.Since(start)
	if err != nil {
		cm.Logger().Error(fmt.Sprintf("Error visiting URL: %s, elapsed time: %s", url, elapsed), err)
		return err
	}

	// Log a debug message
	cm.Logger().Debug(fmt.Sprintf("Visited URL: %s, elapsed time: %s", url, elapsed))

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
	cm.Logger().Debug(fmt.Sprintf("SaveResultsToRedis: Number of results before processing: %d", len(results)))

	client := cm.Client() // Call the Client method to get the client

	for _, result := range results {
		cm.Logger().Debug(fmt.Sprintf("SaveResultsToRedis: Processing result %v", result))

		data, err := json.Marshal(result)
		if err != nil {
			cm.Logger().Error(fmt.Sprintf("SaveResultsToRedis: Error occurred during marshalling to JSON: %v", err), nil)
			return err
		}
		str := string(data)

		err = client.SAdd(ctx, key, str) // Call the SAdd method on the client
		if err != nil {
			cm.Logger().Error(fmt.Sprintf("SaveResultsToRedis: Error occurred during saving to Redis: %v", err), err, nil) // Add nil as the third argument
			return err
		}

		cm.Logger().Debug("SaveResultsToRedis: Added elements to the set")

		// Debugging: Verify that the result was saved correctly
		isMember, err := client.SIsMember(ctx, key, str) // Call the SIsMember method on the client
		if err != nil {
			cm.Logger().Error(fmt.Sprintf("SaveResultsToRedis: Error occurred during checking membership in Redis set: %v", err), err, nil) // Add nil as the third argument
			return err
		}

		if !isMember {
			cm.Logger().Error(fmt.Sprintf("SaveResultsToRedis: Result was not saved correctly in Redis set: %v", str), nil)
		} else {
			cm.Logger().Debug(fmt.Sprintf("SaveResultsToRedis: Result was saved correctly in Redis set, key: %s, result: %s", key, str))
		}
	}

	cm.Logger().Debug(fmt.Sprintf("SaveResultsToRedis: Number of results after processing: %d", len(results)))

	return nil
}
