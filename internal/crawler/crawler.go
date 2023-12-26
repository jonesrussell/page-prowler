package crawler

import (
	"context"
	"net/url"
	"sync"
	"time"

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
	cs.SetupCrawlingLogic(ctx, options)

	cs.Logger.Info("Crawler started...")
	if err := options.Collector.Visit(url); err != nil {
		cs.Logger.Error("Error visiting URL", "url", url, "error", err)
		return nil, err
	}

	options.Collector.Wait()

	cs.Logger.Info("Crawling completed.")

	return *options.Results, nil
}

// ConfigureCollector initializes a new gocolly collector with the specified domains and depth.
func ConfigureCollector(allowedDomains []string, maxDepth int) *colly.Collector {
	collector := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(maxDepth),
	)

	collector.AllowedDomains = allowedDomains

	err := collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       3000 * time.Millisecond,
	})
	if err != nil {
		return nil
	}

	// Respect robots.txt
	collector.AllowURLRevisit = false
	collector.IgnoreRobotsTxt = false

	return collector
}

// HandleHTMLParsing sets up the handler for HTML parsing with gocolly, using the provided parameters.
func (cs *CrawlManager) HandleHTMLParsing(ctx context.Context, options *CrawlOptions) error {
	options.Collector.OnHTML("a[href]", cs.handleAnchorElement(ctx, options))
	return nil
}

func (cs *CrawlManager) handleAnchorElement(ctx context.Context, options *CrawlOptions) func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		href := e.Request.AbsoluteURL(e.Attr("href"))

		options.LinkStatsMu.Lock()
		options.LinkStats.IncrementTotalLinks()
		options.LinkStatsMu.Unlock()

		cs.Logger.Info("Incremented total links count")

		pageData := PageData{
			URL: href,
			// Add other fields as necessary
		}

		if termmatcher.Related(href, options.SearchTerms) {
			options.LinkStatsMu.Lock()
			options.LinkStats.IncrementMatchedLinks()
			options.LinkStatsMu.Unlock()

			cs.Logger.Info("Incremented matched links count")

			if err := cs.handleMatchingLinks(ctx, options, href); err != nil {
				cs.Logger.Error("Error handling matching links", "error", err)
			}
			pageData.MatchingTerms = options.SearchTerms
		} else {
			options.LinkStatsMu.Lock()
			options.LinkStats.IncrementNotMatchedLinks()
			options.LinkStatsMu.Unlock()

			cs.Logger.Info("Incremented not matched links count")

			cs.handleNonMatchingLinks(href)
		}

		options.LinkStatsMu.Lock()
		*options.Results = append(*options.Results, pageData)
		options.LinkStatsMu.Unlock()
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

// handleErrorEvents sets up the error handling for the colly collector.
func (cs *CrawlManager) handleErrorEvents(collector *colly.Collector) {
	collector.OnError(func(r *colly.Response, err error) {
		statusCode := r.StatusCode
		requestURL := r.Request.URL.String()

		if statusCode != 404 {
			cs.Logger.Error("Request URL failed request_url=", requestURL, "status_code=", statusCode, "error=", err)
		}
	})
}

// SetupCrawlingLogic configures and initiates the crawling logic.
func (cs *CrawlManager) SetupCrawlingLogic(ctx context.Context, options *CrawlOptions) {
	cs.Logger.Info("Start SetupCrawlingLogic")
	if options.Debug {
		cs.Logger.Debug("Setting up crawling logic...")
	}

	if err := cs.HandleHTMLParsing(ctx, options); err != nil {
		cs.Logger.Error("Error during HTML parsing", "error", err)
		return
	}

	cs.handleErrorEvents(options.Collector)

	options.Collector.OnScraped(func(r *colly.Response) {
		cs.Logger.Info("Start OnScraped callback", "url", r.Request.URL.String())
		cs.Logger.Info("Finished scraping the page", "url", r.Request.URL.String())

		options.LinkStatsMu.Lock()
		cs.Logger.Info("Total links found", "total_links", options.LinkStats.TotalLinks)
		cs.Logger.Info("Matched links", "matched_links", options.LinkStats.MatchedLinks)
		cs.Logger.Info("Not matched links", "not_matched_links", options.LinkStats.NotMatchedLinks)
		options.LinkStatsMu.Unlock()

		// Here, you would add code to populate the 'results' slice with data
		cs.Logger.Info("End OnScraped callback", "url", r.Request.URL.String())
	})

	options.Collector.OnRequest(func(r *colly.Request) {
		cs.Logger.Info("Start OnRequest callback", "url", r.URL.String())
		cs.Logger.Info("Visiting URL", "url", r.URL.String())
		cs.Logger.Info("End OnRequest callback", "url", r.URL.String())
	})
	cs.Logger.Info("End SetupCrawlingLogic")
}

// GetHostFromURL extracts the host from a given URL string.
func GetHostFromURL(inputURL string, log logger.Logger) (string, error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		log.Fatal("Failed to parse URL", "url", inputURL, "error", err)
		return "", err // return an empty string and the error
	}
	return u.Host, nil // return the host and nil for the error
}
