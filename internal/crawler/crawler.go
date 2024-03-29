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
	DefaultParallelism = 2
	DefaultDelay       = 3000 * time.Millisecond
)

//go:generate mockery --name=CrawlManagerInterface
type CrawlManagerInterface interface {
	Crawl(url string, options *CrawlOptions) ([]PageData, error)
	SetupHTMLParsingHandler(handler func(*colly.HTMLElement)) error
	SetupErrorEventHandler(collector *colly.Collector)
	SetupCrawlingLogic(*CrawlOptions) error
	CrawlURL(url string) error
	HandleVisitError(url string, err error) error
	Logger() logger.Logger
	StartCrawling(ctx context.Context, url string, searchterms string, siteid string, maxdepth int, debug bool) error
	ProcessMatchingLinkAndUpdateStats(*CrawlOptions, string, PageData, []string)
}

var _ CrawlManagerInterface = &CrawlManager{
	LoggerField:    nil,
	Client:         nil,
	MongoDBWrapper: nil,
	Collector:      &colly.Collector{},
	CrawlingMu:     &sync.Mutex{},
	StatsManager:   &StatsManager{},
}

// CrawlOptions represents the options for configuring and initiating the crawling logic.
type CrawlOptions struct {
	CrawlSiteID string
	SearchTerms []string
	Results     *[]PageData
	Debug       bool
}

func (cm *CrawlManager) Logger() logger.Logger {
	return cm.LoggerField
}

func (cm *CrawlManager) StartCrawling(ctx context.Context, url, searchTerms, crawlSiteID string, maxDepth int, debug bool) error {
	if url == "" || searchTerms == "" || crawlSiteID == "" || maxDepth <= 0 {
		return errors.New("invalid parameters")
	}

	// Initialize LinkStats...
	cm.StatsManager = &StatsManager{
		LinkStats:   &stats.Stats{},
		LinkStatsMu: sync.RWMutex{},
	}
	cm.CrawlingMu.Lock()
	defer cm.CrawlingMu.Unlock()

	host, err := GetHostFromURL(url, cm.Logger())
	if err != nil {
		cm.LoggerField.Error(fmt.Sprintf("Failed to parse URL: url: %v, error: %v", url, err))
		return err
	}

	cm.LoggerField.Debug(fmt.Sprintf("Extracted host from URL: %s", host))

	err = cm.ConfigureCollector([]string{host}, maxDepth)
	if err != nil {
		cm.LoggerField.Fatal(fmt.Sprintf("Failed to configure collector: %v", err))
		return err
	}

	splitSearchTerms := cm.splitSearchTerms(searchTerms)
	options := cm.createStartCrawlingOptions(crawlSiteID, splitSearchTerms, debug)

	results, err := cm.Crawl(url, options)
	if err != nil {
		return err
	}

	cm.logCrawlingStatistics()

	err = cm.SaveResultsToRedis(ctx, results, crawlSiteID)
	if err != nil {
		return err
	}

	logResults(cm, results)

	return nil
}

// Crawl starts the crawling process for a given URL with the provided options.
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
func (cm *CrawlManager) SetupHTMLParsingHandler(handler func(*colly.HTMLElement)) error {
	cm.Collector.OnHTML("a[href]", handler)
	return nil
}

// SetupErrorEventHandler sets up the HTTP request error handling for the colly collector.
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
func (cm *CrawlManager) SetupCrawlingLogic(options *CrawlOptions) error {
	err := cm.SetupHTMLParsingHandler(cm.GetAnchorElementHandler(options))
	if err != nil {
		return cm.handleSetupError(err)
	}

	cm.SetupErrorEventHandler(cm.Collector)

	return nil
}

// CrawlURL visits the given URL and performs the crawling operation.
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
func (cm *CrawlManager) HandleVisitError(url string, err error) error {
	cm.LoggerField.Error(fmt.Sprintf("Error visiting URL: url: %s, error: %v", url, err))
	return err
}

