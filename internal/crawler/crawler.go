package crawler

import (
	"context"
	"errors"
	"log"
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
	// Debugging statement
	if ctx.Err() != nil {
		log.Println("Crawl: context error:", ctx.Err())
	} else {
		log.Println("Crawl: context is not done")
	}

	// Debugging statement
	if cs.Collector == nil {
		log.Println("Crawl: cs.Collector is nil")
	} else {
		log.Println("Crawl: cs.Collector is not nil")
	}

	err := cs.setupCrawlingLogic(ctx, options)
	if err != nil {
		return nil, err
	}

	// Debugging statement
	if cs.Collector == nil {
		log.Println("Crawl after setupCrawlingLogic: cs.Collector is nil")
	} else {
		log.Println("Crawl after setupCrawlingLogic: cs.Collector is not nil")
	}

	cs.visitURL(url, options)

	// Debugging statement
	if cs.Collector == nil {
		log.Println("Crawl after visitURL: cs.Collector is nil")
	} else {
		log.Println("Crawl after visitURL: cs.Collector is not nil")
	}

	return cs.handleResults(options), nil
}

// setupHTMLParsingHandler sets up the handler for HTML parsing with gocolly, using the provided parameters.
func (cs *CrawlManager) setupHTMLParsingHandler(ctx context.Context, options *CrawlOptions) error {
	// Debugging statement
	if ctx.Err() != nil {
		log.Println("Crawl: context error:", ctx.Err())
	} else {
		log.Println("Crawl: context is not done")
	}

	// Debugging statement
	if cs.Collector == nil {
		log.Println("setupHTMLParsingHandler: cs.Collector is nil")
	} else {
		log.Println("setupHTMLParsingHandler: cs.Collector is not nil")
	}

	cs.Collector.OnHTML("a[href]", cs.getAnchorElementHandler(ctx, options))

	// Debugging statement
	if cs.Collector == nil {
		log.Println("setupHTMLParsingHandler after OnHTML: cs.Collector is nil")
	} else {
		log.Println("setupHTMLParsingHandler after OnHTML: cs.Collector is not nil")
	}

	return nil
}

func (cs *CrawlManager) getAnchorElementHandler(ctx context.Context, options *CrawlOptions) func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		// Debugging statement
		if ctx.Err() != nil {
			log.Println("Crawl: context error:", ctx.Err())
		} else {
			log.Println("Crawl: context is not done")
		}

		// Debugging statement
		if cs.Collector == nil {
			log.Println("getAnchorElementHandler: cs.Collector is nil")
		} else {
			log.Println("getAnchorElementHandler: cs.Collector is not nil")
		}

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

		// Debugging statement
		if cs.Collector == nil {
			log.Println("getAnchorElementHandler after processing: cs.Collector is nil")
		} else {
			log.Println("getAnchorElementHandler after processing: cs.Collector is not nil")
		}
	}
}

// handleMatchingLinks is responsible for handling the links that match the search criteria during crawling.
func (cs *CrawlManager) handleMatchingLinks(
	ctx context.Context,
	options *CrawlOptions,
	href string,
) error {
	// Debugging statement
	if ctx.Err() != nil {
		log.Println("Crawl: context error:", ctx.Err())
	} else {
		log.Println("Crawl: context is not done")
	}

	cs.Logger.Info("Start handling matching links", "url", href)

	err := cs.Collector.Visit(href)
	if err != nil {
		if errors.Is(err, colly.ErrAlreadyVisited) {
			cs.Logger.Info("URL already visited", "url", href)
		} else if errors.Is(err, colly.ErrForbiddenDomain) {
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
	// Debugging statement
	if ctx.Err() != nil {
		log.Println("Crawl: context error:", ctx.Err())
	} else {
		log.Println("Crawl: context is not done")
	}

	// Debugging statement
	if cs.Collector == nil {
		log.Println("setupCrawlingLogic: cs.Collector is nil")
	} else {
		log.Println("setupCrawlingLogic: cs.Collector is not nil")
	}

	err := cs.setupHTMLParsingHandler(ctx, options)
	if err != nil {
		return err
	}

	cs.setupErrorEventHandler(cs.Collector)
	cs.handleRequestEvents(options)

	// Debugging statement
	if cs.Collector == nil {
		log.Println("setupCrawlingLogic after handleRequestEvents: cs.Collector is nil")
	} else {
		log.Println("setupCrawlingLogic after handleRequestEvents: cs.Collector is not nil")
	}

	return nil
}

func (cs *CrawlManager) visitURL(url string, options *CrawlOptions) {
	// Debugging statement
	if cs.Collector == nil {
		log.Println("visitURL: cs.Collector is nil")
	} else {
		log.Println("visitURL: cs.Collector is not nil")
	}

	if err := cs.Collector.Visit(url); err != nil {
		cs.Logger.Error("Error visiting URL", "url", url, "error", err)
	}
	cs.Collector.Wait()
	cs.Logger.Info("Crawling completed.")

	// Debugging statement
	if cs.Collector == nil {
		log.Println("visitURL after Wait: cs.Collector is nil")
	} else {
		log.Println("visitURL after Wait: cs.Collector is not nil")
	}
}

func (cs *CrawlManager) handleResults(options *CrawlOptions) []PageData {
	// Debugging statement
	if cs.Collector == nil {
		log.Println("handleResults: cs.Collector is nil")
	} else {
		log.Println("handleResults: cs.Collector is not nil")
	}

	results := *options.Results

	// Debugging statement
	if cs.Collector == nil {
		log.Println("handleResults after getting results: cs.Collector is nil")
	} else {
		log.Println("handleResults after getting results: cs.Collector is not nil")
	}

	return results
}

func (cs *CrawlManager) processMatchingLinkAndUpdateStats(ctx context.Context, options *CrawlOptions, href string, pageData PageData, matchingTerms []string) {
	// Debugging statement
	if ctx.Err() != nil {
		log.Println("Crawl: context error:", ctx.Err())
	} else {
		log.Println("Crawl: context is not done")
	}

	// Debugging statement
	if cs.Collector == nil {
		log.Println("processMatchingLinkAndUpdateStats: cs.Collector is nil")
	} else {
		log.Println("processMatchingLinkAndUpdateStats: cs.Collector is not nil")
	}

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

	// Debugging statement
	if cs.Collector == nil {
		log.Println("processMatchingLinkAndUpdateStats after appending results: cs.Collector is nil")
	} else {
		log.Println("processMatchingLinkAndUpdateStats after appending results: cs.Collector is not nil")
	}
}

func (cs *CrawlManager) incrementNonMatchedLinkCount(options *CrawlOptions, href string) {
	// Debugging statement
	if cs.Collector == nil {
		log.Println("incrementNonMatchedLinkCount: cs.Collector is nil")
	} else {
		log.Println("incrementNonMatchedLinkCount: cs.Collector is not nil")
	}

	options.LinkStatsMu.Lock()
	options.LinkStats.IncrementNotMatchedLinks()
	options.LinkStatsMu.Unlock()
	cs.Logger.Info("Incremented not matched links count")
	cs.handleNonMatchingLinks(href)

	// Debugging statement
	if cs.Collector == nil {
		log.Println("incrementNonMatchedLinkCount after handleNonMatchingLinks: cs.Collector is nil")
	} else {
		log.Println("incrementNonMatchedLinkCount after handleNonMatchingLinks: cs.Collector is not nil")
	}
}

func (cs *CrawlManager) handleRequestEvents(options *CrawlOptions) {
	// Create a local copy of Collector
	collector := cs.Collector

	// Debugging statement
	if collector == nil {
		log.Println("handleRequestEvents: collector is nil")
	} else {
		log.Println("handleRequestEvents: collector is not nil")
	}

	collector.OnRequest(func(r *colly.Request) {
		// Debugging statement
		if collector == nil {
			log.Println("handleRequestEvents OnRequest callback: collector is nil")
		} else {
			log.Println("handleRequestEvents OnRequest callback: collector is not nil")
		}

		cs.Logger.Info("Start OnRequest callback", "url", r.URL.String())
		cs.Logger.Info("Visiting URL", "url", r.URL.String())
		cs.Logger.Info("End OnRequest callback", "url", r.URL.String())
	})
}

// ConfigureCollector initializes a new gocolly collector with the specified domains and depth.
func (cs *CrawlManager) ConfigureCollector(allowedDomains []string, maxDepth int) error {
	collector := colly.NewCollector(
		colly.Async(false),
		colly.MaxDepth(maxDepth),
	)

	collector.AllowedDomains = allowedDomains

	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       3000 * time.Millisecond,
	})

	cs.Collector = collector

	// Respect robots.txt
	cs.Collector.AllowURLRevisit = false
	cs.Collector.IgnoreRobotsTxt = false

	return nil
}

// StartCrawling starts the crawling process.
func StartCrawling(ctx context.Context, url, searchTerms, crawlSiteID string, maxDepth int, debug bool, crawlerService *CrawlManager, server *CrawlServer) error {
	// Debugging statement
	if ctx.Err() != nil {
		log.Println("Crawl: context error:", ctx.Err())
	} else {
		log.Println("Crawl: context is not done")
	}

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

	// Debugging statement
	log.Println("StartCrawling: collector is not nil")

	var results []PageData

	options := CrawlOptions{
		CrawlSiteID: crawlSiteID,
		SearchTerms: splitSearchTerms,
		Results:     &results,
		LinkStats:   stats.NewStats(),
		Debug:       debug,
	}

	// Debugging statement
	if crawlerService.Collector == nil {
		log.Println("StartCrawling before Crawl: cs.Collector is nil")
	} else {
		log.Println("StartCrawling before Crawl: cs.Collector is not nil")
	}

	results, err = crawlerService.Crawl(ctx, url, &options)
	if err != nil {
		return err
	}

	// Debugging statement
	if crawlerService.Collector == nil {
		log.Println("StartCrawling after Crawl: cs.Collector is nil")
	} else {
		log.Println("StartCrawling after Crawl: cs.Collector is not nil")
	}

	err = server.SaveResultsToRedis(ctx, results, crawlSiteID)
	if err != nil {
		return err
	}

	printResults(crawlerService, results)

	return nil
}
