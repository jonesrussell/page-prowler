package crawler

import (
	"context"
	"sync"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/mongodbwrapper"
	"github.com/jonesrussell/page-prowler/internal/stats"
	"github.com/jonesrussell/page-prowler/internal/termmatcher"
	"github.com/jonesrussell/page-prowler/redis"
)

// CrawlManager encapsulates shared dependencies for crawler functions.
type CrawlManager struct {
	Logger         logger.Logger
	Client         redis.Datastore
	MongoDBWrapper *mongodbwrapper.MongoDBWrapper
}

// CrawlOptions represents the options for configuring and initiating the crawling logic.
type CrawlOptions struct {
	CrawlSiteID string
	Collector   *colly.Collector
	SearchTerms []string
	Results     *[]PageData
	LinkStats   *stats.Stats
	LinkStatsMu sync.Mutex // Mutex for LinkStats
	Debug       bool
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
	options.Collector.OnHTML("a[href]", cs.getAnchorElementHandler(ctx, options))
	return nil
}

func (cs *CrawlManager) getAnchorElementHandler(ctx context.Context, options *CrawlOptions) func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		href := e.Request.AbsoluteURL(e.Attr("href"))
		options.LinkStatsMu.Lock()
		options.LinkStats.IncrementTotalLinks()
		options.LinkStatsMu.Unlock()
		cs.Logger.Info("Incremented total links count")
		pageData := PageData{
			URL: href,
		}
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
	cs.Logger.Info("Start handling matching links", "url", href)

	err := options.Collector.Visit(href)
	if err != nil {
		if err == colly.ErrAlreadyVisited {
			cs.Logger.Info("URL already visited", "url", href)
		} else if err == colly.ErrForbiddenDomain {
			cs.Logger.Info("Forbidden domain - Skipping visit", "url", href)
		} else {
			cs.Logger.Error("Error visiting URL", "url", href, "error", err)
			cs.Logger.Info("End handling matching links", "url", href)
			return err
		}
	}

	cs.Logger.Info("End handling matching links", "url", href)
	return nil
}

// handleNonMatchingLinks logs the occurrence of a non-matching link.
func (cs *CrawlManager) handleNonMatchingLinks(href string) {
	cs.Logger.Info("Non-matching link", "url", href)
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
	cs.setupErrorEventHandler(options.Collector)
	cs.handleRequestEvents(options)
	return nil
}

func (cs *CrawlManager) visitURL(url string, options *CrawlOptions) {
	if err := options.Collector.Visit(url); err != nil {
		cs.Logger.Error("Error visiting URL", "url", url, "error", err)
	}
	options.Collector.Wait()
	cs.Logger.Info("Crawling completed.")
}

func (cs *CrawlManager) handleResults(options *CrawlOptions) []PageData {
	return *options.Results
}

func (cs *CrawlManager) processMatchingLinkAndUpdateStats(ctx context.Context, options *CrawlOptions, href string, pageData PageData, matchingTerms []string) {
	options.LinkStatsMu.Lock()
	options.LinkStats.IncrementMatchedLinks()
	options.LinkStatsMu.Unlock()
	cs.Logger.Info("Incremented matched links count")
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
	cs.Logger.Info("Incremented not matched links count")
	cs.handleNonMatchingLinks(href)
}

func (cs *CrawlManager) handleRequestEvents(options *CrawlOptions) {
	options.Collector.OnRequest(func(r *colly.Request) {
		cs.Logger.Info("Start OnRequest callback", "url", r.URL.String())
		cs.Logger.Info("Visiting URL", "url", r.URL.String())
		cs.Logger.Info("End OnRequest callback", "url", r.URL.String())
	})
}
