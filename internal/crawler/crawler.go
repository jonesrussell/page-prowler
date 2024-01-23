package crawler

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jonesrussell/page-prowler/internal/prowlredis"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/mongodbwrapper"
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
	SetupCrawlingLogic(options *CrawlOptions) error
	CrawlURL(url string, options *CrawlOptions) error
	HandleVisitError(url string, err error) error
	LogError(message string, keysAndValues ...interface{})
	Logger() logger.Logger
	StartCrawling(ctx context.Context, url string, searchterms string, crawlsiteid string, maxdepth int, debug bool) error
	GetMatchedLinkProcessor() MatchedLinkProcessor
	ProcessMatchingLinkAndUpdateStats(options *CrawlOptions, href string, pageData PageData, matchingTerms []string)
}

// CrawlManager encapsulates shared dependencies for crawler functions.
type CrawlManager struct {
	LoggerField          logger.Logger
	Client               prowlredis.ClientInterface
	MongoDBWrapper       mongodbwrapper.MongoDBInterface
	Collector            *colly.Collector
	CrawlingMu           sync.Mutex
	VisitedPages         map[string]bool
	MatchedLinkProcessor MatchedLinkProcessor
}

func (cm *CrawlManager) Debug(msg string, keysAndValues ...interface{}) {
	cm.LoggerField.Debug(msg, keysAndValues...)
}

func (cm *CrawlManager) Info(msg string, keysAndValues ...interface{}) {
	cm.LoggerField.Info(msg, keysAndValues...)
}

func (cm *CrawlManager) Error(msg string, keysAndValues ...interface{}) {
	cm.LoggerField.Error(msg, keysAndValues...)
}

func (cm *CrawlManager) Errorf(msg string, keysAndValues ...interface{}) {
	cm.LoggerField.Errorf(msg, keysAndValues...)
}

func (cm *CrawlManager) Fatal(msg string, keysAndValues ...interface{}) {
	cm.LoggerField.Fatal(msg, keysAndValues...)
}

var _ CrawlManagerInterface = &CrawlManager{}

// CrawlOptions represents the options for configuring and initiating the crawling logic.
type CrawlOptions struct {
	CrawlSiteID string
	SearchTerms []string
	Results     *[]PageData
	LinkStats   *stats.Stats
	LinkStatsMu sync.Mutex // Mutex for LinkStats
	Debug       bool
}

func (cm *CrawlManager) Logger() logger.Logger {
	return cm.LoggerField
}

func (cm *CrawlManager) GetMatchedLinkProcessor() MatchedLinkProcessor {
	return cm.MatchedLinkProcessor
}

// crawl starts the crawling process for a given URL with the provided options.
func (cs *CrawlManager) Crawl(url string, options *CrawlOptions) ([]PageData, error) {
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
			cs.Debug("[SetupErrorEventHandler] Internal Server Error",
				"request_url", requestURL,
				"status_code", fmt.Sprintf("%d", statusCode))
		} else if statusCode != 404 {
			// Handle other errors normally
			cs.Debug("[SetupErrorEventHandler] Request URL failed",
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
	cs.Logger().Debug("Visiting URL", "url", url)
	err := cs.visitWithColly(url)
	if err != nil {
		return cs.HandleVisitError(url, err)
	}
	cs.trackVisitedPage(url, options)
	cs.Collector.Wait()
	cs.Logger().Info("Crawling completed.")
	return nil
}

// HandleVisitError handles the error occurred during the visit of a URL.
func (cs *CrawlManager) HandleVisitError(url string, err error) error {
	cs.LogError("Error visiting URL", "url", url, "error", err)
	return err
}

// LogError logs the error message along with the provided key-value pairs.
func (cs *CrawlManager) LogError(message string, keysAndValues ...interface{}) {
	cs.Error(message, keysAndValues...)
}
