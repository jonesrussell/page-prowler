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

var ErrVisit = errors.New("error visiting URL")

//go:generate mockery --name=CrawlManagerInterface
type CrawlManagerInterface interface {
	Crawl(url string, options *CrawlOptions) ([]PageData, error)
	SetupHTMLParsingHandler(handler func(*colly.HTMLElement)) error
	SetupErrorEventHandler(collector *colly.Collector)
	SetupCrawlingLogic(*CrawlOptions) error
	CrawlURL(url string, options *CrawlOptions) error
	HandleVisitError(url string, err error) error
	LogError(message string, keysAndValues ...interface{})
	Logger() logger.Logger
	StartCrawling(ctx context.Context, url string, searchterms string, siteid string, maxdepth int, debug bool) error
	ProcessMatchingLinkAndUpdateStats(*CrawlOptions, string, PageData, []string)
}

var _ CrawlManagerInterface = &CrawlManager{}

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

func (cs *CrawlManager) StartCrawling(ctx context.Context, url, searchTerms, crawlSiteID string, maxDepth int, debug bool) error {
	if url == "" || searchTerms == "" || crawlSiteID == "" || maxDepth <= 0 {
		return errors.New("invalid parameters")
	}

	// Initialize LinkStats...
	cs.StatsManager = &StatsManager{
		LinkStats:   &stats.Stats{},
		LinkStatsMu: sync.RWMutex{},
	}
	cs.CrawlingMu.Lock()
	defer cs.CrawlingMu.Unlock()

	host, err := GetHostFromURL(url, cs.Logger())
	if err != nil {
		cs.LoggerField.Error("Failed to parse URL", "url", url, "error", err)
		return err
	}

	cs.LoggerField.Debug("Extracted host from URL", "host", host)

	err = cs.ConfigureCollector([]string{host}, maxDepth)
	if err != nil {
		cs.Logger().Fatal("Failed to configure collector", "error", err)
		return err
	}

	splitSearchTerms := cs.splitSearchTerms(searchTerms)
	options := cs.createStartCrawlingOptions(crawlSiteID, splitSearchTerms, debug)

	results, err := cs.Crawl(url, options)
	if err != nil {
		return err
	}

	cs.logCrawlingStatistics(options)

	err = cs.SaveResultsToRedis(ctx, results, crawlSiteID)
	if err != nil {
		return err
	}

	logResults(cs, results)

	return nil
}

// crawl starts the crawling process for a given URL with the provided options.
func (cs *CrawlManager) Crawl(url string, options *CrawlOptions) ([]PageData, error) {
	cs.LoggerField.Debug("CrawlURL", "url", url)
	err := cs.SetupCrawlingLogic(options)
	if err != nil {
		return nil, err
	}

	err = cs.CrawlURL(url, options)
	if err != nil {
		return nil, err
	}

	return *options.Results, nil
}

// SetupHTMLParsingHandler sets up the handler for HTML parsing with gocolly, using the provided parameters.
func (cs *CrawlManager) SetupHTMLParsingHandler(handler func(*colly.HTMLElement)) error {
	cs.Collector.OnHTML("a[href]", handler)
	return nil
}

// SetupErrorEventHandler sets up the HTTP request error handling for the colly collector.
func (cs *CrawlManager) SetupErrorEventHandler(collector *colly.Collector) {
	collector.OnError(func(r *colly.Response, err error) {
		statusCode := r.StatusCode
		requestURL := r.Request.URL.String()

		if statusCode == 500 {
			// Handle 500 Internal Server Error without printing the stack trace
			cs.LoggerField.Debug("[SetupErrorEventHandler] Internal Server Error",
				"request_url", requestURL,
				"status_code", fmt.Sprintf("%d", statusCode))
		} else if statusCode != 404 {
			// Handle other errors normally
			cs.LoggerField.Debug("[SetupErrorEventHandler] Request URL failed",
				"request_url", requestURL,
				"status_code", fmt.Sprintf("%d", statusCode))
		}
	})
}

// SetupCrawlingLogic configures and initiates the crawling logic.
func (cs *CrawlManager) SetupCrawlingLogic(options *CrawlOptions) error {
	err := cs.SetupHTMLParsingHandler(cs.getAnchorElementHandler(options))
	if err != nil {
		return cs.handleSetupError(err)
	}

	cs.SetupErrorEventHandler(cs.Collector)

	return nil
}

// CrawlURL visits the given URL and performs the crawling operation.
func (cs *CrawlManager) CrawlURL(url string, options *CrawlOptions) error {
	cs.Logger().Debug("[CrawlURL] Visiting URL", "url", url)
	err := cs.visitWithColly(url)
	if err != nil {
		return cs.HandleVisitError(url, err)
	}
	//	cs.trackVisitedPage(url, options)
	cs.Collector.Wait()
	cs.Logger().Info("[CrawlURL] Crawling completed.")
	return nil
}

// HandleVisitError handles the error occurred during the visit of a URL.
func (cs *CrawlManager) HandleVisitError(url string, err error) error {
	cs.LogError("Error visiting URL", "url", url, "error", err)
	return err
}

// LogError logs the error message along with the provided key-value pairs.
func (cs *CrawlManager) LogError(message string, keysAndValues ...interface{}) {
	cs.LoggerField.Error(message, keysAndValues...)
}
