package crawler

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/mongodbwrapper"
	"github.com/jonesrussell/page-prowler/internal/stats"
	"github.com/jonesrussell/page-prowler/internal/termmatcher"
	"github.com/jonesrussell/page-prowler/redis"
)

type echoContextKey string

const (
	echoManagerKey echoContextKey = "manager"
)

// CrawlManager encapsulates shared dependencies for crawler functions.
type CrawlManager struct {
	Logger         logger.Logger
	Client         redis.ClientInterface
	MongoDBWrapper mongodbwrapper.MongoDBInterface
	Collector      *colly.Collector
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

func NewCrawlManager(logger logger.Logger, client redis.ClientInterface, mongoDBWrapper mongodbwrapper.MongoDBInterface) *CrawlManager {
	return &CrawlManager{
		Logger:         logger,
		Client:         client,
		MongoDBWrapper: mongoDBWrapper,
	}
}

func (cs *CrawlManager) Crawl(ctx context.Context, url string, options *CrawlOptions) ([]PageData, error) {
	err := cs.setupCrawlingLogic(ctx, options)
	if err != nil {
		return nil, err
	}

	cs.visitURL(url, options)

	return cs.handleResults(options), nil
}

// setupHTMLParsingHandler sets up the handler for HTML parsing with gocolly, using the provided parameters.
func (cs *CrawlManager) setupHTMLParsingHandler(ctx context.Context, options *CrawlOptions) error {
	cs.Collector.OnHTML("a[href]", cs.getAnchorElementHandler(ctx, options))

	return nil
}

func (cs *CrawlManager) getAnchorElementHandler(ctx context.Context, options *CrawlOptions) func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		href := e.Request.AbsoluteURL(e.Attr("href"))
		options.LinkStatsMu.Lock()
		options.LinkStats.IncrementTotalLinks()
		options.LinkStatsMu.Unlock()
		cs.Logger.Debug("Incremented total links count")
		pageData := PageData{
			URL: href,
		}
		cs.Logger.Debug("Search terms: %v", options.SearchTerms)
		matchingTerms := termmatcher.GetMatchingTerms(href, options.SearchTerms)
		if len(matchingTerms) > 0 {
			cs.processMatchingLinkAndUpdateStats(ctx, options, href, pageData, matchingTerms)
		} else {
			cs.incrementNonMatchedLinkCount(options, href)
		}
	}
}

// handleMatchingLinks is responsible for handling the links that match the search criteria during crawling.
func (cs *CrawlManager) handleMatchingLinks(
	ctx context.Context,
	options *CrawlOptions,
	href string,
) error {
	cs.Logger.Debug("Start handling matching links", "url", href)

	err := cs.Collector.Visit(href)
	if err != nil {
		if errors.Is(err, colly.ErrAlreadyVisited) {
			cs.Logger.Debug("URL already visited", "url", href)
		} else if errors.Is(err, colly.ErrForbiddenDomain) {
			cs.Logger.Debug("Forbidden domain - Skipping visit", "url", href)
		} else {
			cs.Logger.Error("Error visiting URL", "url", href, "error", err)
			cs.Logger.Debug("End handling matching links", "url", href)
			return err
		}
	}

	cs.Logger.Debug("End handling matching links", "url", href)
	return nil
}

// handleNonMatchingLinks logs the occurrence of a non-matching link.
func (cs *CrawlManager) handleNonMatchingLinks(href string) {
	cs.Logger.Debug("Non-matching link", "url", href)
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
func (cs *CrawlManager) setupCrawlingLogic(ctx context.Context, options *CrawlOptions) error {
	err := cs.setupHTMLParsingHandler(ctx, options)
	if err != nil {
		return err
	}

	cs.setupErrorEventHandler(cs.Collector)
	cs.handleRequestEvents()

	return nil
}

func (cs *CrawlManager) visitURL(url string, options *CrawlOptions) {
	if err := cs.Collector.Visit(url); err != nil {
		cs.Logger.Error("Error visiting URL", "url", url, "error", err)
	}
	cs.Collector.Wait()
	cs.Logger.Info("Crawling completed.")
}

func (cs *CrawlManager) handleResults(options *CrawlOptions) []PageData {
	results := *options.Results

	return results
}

func (cs *CrawlManager) processMatchingLinkAndUpdateStats(ctx context.Context, options *CrawlOptions, href string, pageData PageData, matchingTerms []string) {
	options.LinkStatsMu.Lock()
	options.LinkStats.IncrementMatchedLinks()
	options.LinkStatsMu.Unlock()
	cs.Logger.Debug("Incremented matched links count")
	if err := cs.handleMatchingLinks(ctx, options, href); err != nil {
		cs.Logger.Error("Error handling matching links", "error", err)
	}
	pageData.MatchingTerms = matchingTerms
	options.LinkStatsMu.Lock()
	*options.Results = append(*options.Results, pageData)
	options.LinkStatsMu.Unlock()
}

func (cs *CrawlManager) incrementNonMatchedLinkCount(options *CrawlOptions, href string) {
	options.LinkStatsMu.Lock()
	options.LinkStats.IncrementNotMatchedLinks()
	options.LinkStatsMu.Unlock()
	cs.Logger.Debug("Incremented not matched links count")
	cs.handleNonMatchingLinks(href)
}

func (cs *CrawlManager) handleRequestEvents() {
	cs.Collector.OnRequest(func(r *colly.Request) {
		cs.Logger.Debug("Start OnRequest callback", "url", r.URL.String())
		cs.Logger.Debug("Visiting URL", "url", r.URL.String())
		cs.Logger.Debug("End OnRequest callback", "url", r.URL.String())
	})
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
func StartCrawling(ctx context.Context, url, searchTerms, crawlSiteID string, maxDepth int, debug bool, crawlerService *CrawlManager) error {
	splitSearchTerms := strings.Split(searchTerms, ",")
	host, err := GetHostFromURL(url, crawlerService.Logger)
	if err != nil {
		crawlerService.Logger.Error("Failed to parse URL", "url", url, "error", err)
		return err
	}

	err = crawlerService.ConfigureCollector([]string{host}, maxDepth)
	if err != nil {
		crawlerService.Logger.Fatal("Failed to configure collector", "error", err)
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

	results, err = crawlerService.Crawl(ctx, url, &options)
	if err != nil {
		return err
	}

	err = crawlerService.SaveResultsToRedis(ctx, results, crawlSiteID)
	if err != nil {
		return err
	}

	printResults(crawlerService, results)

	return nil
}
