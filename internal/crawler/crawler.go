package crawler

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/jonesrussell/page-prowler/internal/prowlredis"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/mongodbwrapper"
	"github.com/jonesrussell/page-prowler/internal/stats"
	"github.com/jonesrussell/page-prowler/internal/termmatcher"
)

type CrawlManagerInterface interface {
	StartCrawling(ctx context.Context, url, searchTerms, crawlSiteID string, maxDepth int, debug bool) error
}

// CrawlManager encapsulates shared dependencies for crawler functions.
type CrawlManager struct {
	Logger         logger.Logger
	Client         prowlredis.ClientInterface
	MongoDBWrapper mongodbwrapper.MongoDBInterface
	Collector      *colly.Collector
	CrawlingMu     sync.Mutex
}

// CrawlOptions represents the options for configuring and initiating the crawling logic.
type CrawlOptions struct {
	CrawlSiteID string
	SearchTerms []string
	Results     *[]PageData
	LinkStats   *stats.Stats
	LinkStatsMu sync.Mutex // Mutex for LinkStats
	Debug       bool
}

// NewCrawlManager creates a new instance of CrawlManager.
func NewCrawlManager(
	logger logger.Logger,
	client prowlredis.ClientInterface,
	mongoDBWrapper mongodbwrapper.MongoDBInterface,
) *CrawlManager {
	return &CrawlManager{
		Logger:         logger,
		Client:         client,
		MongoDBWrapper: mongoDBWrapper,
	}
}

func (cs *CrawlManager) crawl(url string, options *CrawlOptions) ([]PageData, error) {
	err := cs.setupCrawlingLogic(options)
	if err != nil {
		return nil, err
	}

	cs.visitURL(url)

	return cs.handleResults(options), nil
}

// setupHTMLParsingHandler sets up the handler for HTML parsing with gocolly, using the provided parameters.
func (cs *CrawlManager) setupHTMLParsingHandler(options *CrawlOptions, handler func(*colly.HTMLElement)) error {
	cs.Collector.OnHTML("a[href]", handler)
	return nil
}

func (cs *CrawlManager) getAnchorElementHandler(options *CrawlOptions) func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		href := e.Request.AbsoluteURL(e.Attr("href"))
		anchorText := e.Text
		options.LinkStatsMu.Lock()
		options.LinkStats.IncrementTotalLinks()
		options.LinkStatsMu.Unlock()
		cs.Logger.Debug("Incremented total links count")
		pageData := PageData{
			URL: href,
		}
		cs.Logger.Debug("Search terms: %v", options.SearchTerms)
		matchingTerms := termmatcher.GetMatchingTerms(href, anchorText, options.SearchTerms)
		if len(matchingTerms) > 0 {
			cs.processMatchingLinkAndUpdateStats(options, href, pageData, matchingTerms)
		} else {
			cs.incrementNonMatchedLinkCount(options)
		}
	}
}

// handleMatchingLinks is responsible for handling the links that match the search criteria during crawling.
func (cs *CrawlManager) handleMatchingLinks(href string) error {
	cs.Logger.Debug("Start handling matching links", "url", href)

	err := cs.visit(href)
	if err != nil {
		cs.Logger.Error("Error visiting URL", "url", href, "error", err)
		return err
	}

	cs.Logger.Debug("End handling matching links", "url", href)
	return nil
}

// setupErrorEventHandler sets up the error handling for the colly collector.
func (cs *CrawlManager) setupErrorEventHandler(collector *colly.Collector) {
	collector.OnError(func(r *colly.Response, err error) {
		statusCode := r.StatusCode
		requestURL := r.Request.URL.String()

		if statusCode != 404 {
			cs.Logger.Error("Request URL failed request_url=", requestURL, "status_code=", statusCode, "error=", err)
		}
	})
}

// setupCrawlingLogic configures and initiates the crawling logic.
func (cs *CrawlManager) setupCrawlingLogic(options *CrawlOptions) error {
	err := cs.setupHTMLParsingHandler(options, cs.getAnchorElementHandler(options))
	if err != nil {
		return err
	}

	cs.setupErrorEventHandler(cs.Collector)

	return nil
}

func (cs *CrawlManager) visitURL(url string) {
	err := cs.visit(url)
	if err != nil {
		cs.Logger.Error("Error visiting URL", "url", url, "error", err)
	}

	cs.Collector.Wait()
	cs.Logger.Info("Crawling completed.")
}

func (cs *CrawlManager) handleResults(options *CrawlOptions) []PageData {
	results := *options.Results

	return results
}

func (cs *CrawlManager) processMatchingLinkAndUpdateStats(options *CrawlOptions, href string, pageData PageData, matchingTerms []string) {
	options.LinkStatsMu.Lock()
	options.LinkStats.IncrementMatchedLinks()
	options.LinkStatsMu.Unlock()
	cs.Logger.Debug("Incremented matched links count")
	if err := cs.handleMatchingLinks(href); err != nil {
		cs.Logger.Error("Error handling matching links", "error", err)
	}
	pageData.MatchingTerms = matchingTerms
	options.LinkStatsMu.Lock()
	*options.Results = append(*options.Results, pageData)
	options.LinkStatsMu.Unlock()
}

func (cs *CrawlManager) incrementNonMatchedLinkCount(options *CrawlOptions) {
	options.LinkStatsMu.Lock()
	options.LinkStats.IncrementNotMatchedLinks()
	options.LinkStatsMu.Unlock()
	cs.Logger.Debug("Incremented not matched links count")
}

// ConfigureCollector initializes a new gocolly collector with the specified domains and depth.
func (cs *CrawlManager) ConfigureCollector(allowedDomains []string, maxDepth int) error {
	collector := colly.NewCollector(
		colly.Async(false),
		colly.MaxDepth(maxDepth),
	)

	collector.AllowedDomains = allowedDomains

	limitRule := &colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       3000 * time.Millisecond,
	}

	if err := collector.Limit(limitRule); err != nil {
		cs.Logger.Errorf("Failed to set limit rule: %v", err)
		return err
	}

	cs.Collector = collector

	// Respect robots.txt
	cs.Collector.AllowURLRevisit = false
	cs.Collector.IgnoreRobotsTxt = false

	return nil
}

// StartCrawling starts the crawling process.
func (cs *CrawlManager) StartCrawling(ctx context.Context, url, searchTerms, crawlSiteID string, maxDepth int, debug bool) error {
	cs.CrawlingMu.Lock()
	defer cs.CrawlingMu.Unlock()

	splitSearchTerms := strings.Split(searchTerms, ",")
	host, err := GetHostFromURL(url, cs.Logger)
	if err != nil {
		cs.Logger.Error("Failed to parse URL", "url", url, "error", err)
		return err
	}

	err = cs.ConfigureCollector([]string{host}, maxDepth)
	if err != nil {
		cs.Logger.Fatal("Failed to configure collector", "error", err)
		return err
	}

	var results []PageData

	options := CrawlOptions{
		CrawlSiteID: crawlSiteID,
		SearchTerms: splitSearchTerms,
		Results:     &results,
		LinkStats:   stats.NewStats(),
		Debug:       debug,
	}

	results, err = cs.crawl(url, &options)
	if err != nil {
		return err
	}

	err = cs.SaveResultsToRedis(ctx, results, crawlSiteID)
	if err != nil {
		return err
	}

	printResults(cs, results)

	return nil
}

func (cs *CrawlManager) visit(url string) error {
	err := cs.Collector.Visit(url)
	if err != nil {
		if errors.Is(err, colly.ErrAlreadyVisited) {
			cs.Logger.Debug("URL already visited", "url", url)
		} else if errors.Is(err, colly.ErrForbiddenDomain) {
			cs.Logger.Debug("Forbidden domain - Skipping visit", "url", url)
		} else {
			return err
		}
	}
	return nil
}
