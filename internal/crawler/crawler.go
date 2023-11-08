// Package crawler provides the tools and logic needed to perform web crawling and data extraction.
package crawler

import (
	"context"
	"net/url"
	"time"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/crawlresult"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/rediswrapper"
	"github.com/jonesrussell/page-prowler/internal/stats"
	"github.com/jonesrussell/page-prowler/internal/termmatcher"
)

// CrawlManager encapsulates shared dependencies for crawler functions.
type CrawlManager struct {
	Logger       logger.Logger
	RedisWrapper *rediswrapper.RedisWrapper
}

// CrawlOptions represents the options for configuring and initiating the crawling logic.
type CrawlOptions struct {
	CrawlSiteID string
	Collector   *colly.Collector
	SearchTerms []string
	Results     *[]crawlresult.PageData
	LinkStats   *stats.Stats
	Debug       bool
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

	return collector
}

// HandleHTMLParsing sets up the handler for HTML parsing with gocolly, using the provided parameters.
func (cs *CrawlManager) HandleHTMLParsing(ctx context.Context, options CrawlOptions) error {
	options.Collector.OnHTML("a[href]", cs.handleAnchorElement(ctx, options))
	return nil
}

func (cs *CrawlManager) handleAnchorElement(ctx context.Context, options CrawlOptions) func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		href := e.Request.AbsoluteURL(e.Attr("href"))
		options.LinkStats.IncrementTotalLinks()

		pageData := crawlresult.PageData{
			URL: href,
			// Add other fields as necessary
		}

		if termmatcher.Related(href, options.SearchTerms) {
			options.LinkStats.IncrementMatchedLinks()
			if err := cs.handleMatchingLinks(ctx, options, href); err != nil {
				cs.Logger.Error("Error handling matching links", "error", err)
			}
			pageData.MatchingTerms = options.SearchTerms
		} else {
			options.LinkStats.IncrementNotMatchedLinks()
			cs.handleNonMatchingLinks(href)
		}

		*options.Results = append(*options.Results, pageData)
	}
}

// handleMatchingLinks is responsible for handling the links that match the search criteria during crawling.
func (cs *CrawlManager) handleMatchingLinks(
	ctx context.Context,
	options CrawlOptions,
	href string,
) error {
	cs.Logger.Info("Found URL", "url", href)
	if _, err := cs.RedisWrapper.SAdd(ctx, options.CrawlSiteID, href); err != nil {
		cs.Logger.Error("Error adding URL to Redis set", "set", options.CrawlSiteID, "error", err)
		return err
	}

	if options.Debug {
		cs.Logger.Debug("Added URL to Redis set", "set", options.CrawlSiteID, "url", href)
	}

	err := options.Collector.Visit(href)
	if err != nil {
		if err == colly.ErrAlreadyVisited {
			cs.Logger.Info("URL already visited", "url", href)
		} else if err == colly.ErrForbiddenDomain {
			cs.Logger.Info("Forbidden domain - Skipping visit", "url", href)
		} else {
			cs.Logger.Error("Error visiting URL", "url", href, "error", err)
			return err
		}
	}

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

// handleRedisOperations manages the Redis operations after crawling a page.
func (cs *CrawlManager) handleRedisOperations(ctx context.Context) error {
	// Retrieve the members of the set from Redis.
	hrefs, err := cs.RedisWrapper.SMembers(ctx, "yourKeyHere") // Replace "yourKeyHere" with the actual key you're interested in.
	if err != nil {
		cs.Logger.Error("Error getting members from Redis error=", err)
		return err
	}

	// Iterate over the set members and publish each href to the specified channel.
	for _, href := range hrefs {
		err = cs.RedisWrapper.PublishHref(ctx, "streetcode", href)
		if err != nil {
			cs.Logger.Error("Error publishing href to Redis href=", href, "error=", err)
			return err
		}

		// Delete the href from Redis now that it's been published.
		if _, err = cs.RedisWrapper.Del(ctx, href); err != nil {
			cs.Logger.Error("Error deleting href from Redis href=", href, "error=", err)
			return err
		}
	}

	// If no errors occurred, return nil to indicate success.
	return nil
}

// SetupCrawlingLogic configures and initiates the crawling logic.
func (cs *CrawlManager) SetupCrawlingLogic(ctx context.Context, options CrawlOptions) {
	if options.Debug {
		cs.Logger.Debug("Setting up crawling logic...")
	}

	if err := cs.HandleHTMLParsing(ctx, options); err != nil {
		cs.Logger.Error("Error during HTML parsing", "error", err)
		return
	}

	cs.handleErrorEvents(options.Collector)

	options.Collector.OnScraped(func(r *colly.Response) {
		if err := cs.handleRedisOperations(ctx); err != nil {
			cs.Logger.Error("Error with Redis operations", "error", err)
			return
		}

		cs.Logger.Info("Finished scraping the page", "url", r.Request.URL.String())
		cs.Logger.Info("Total links found", "total_links", options.LinkStats.TotalLinks)
		cs.Logger.Info("Matched links", "matched_links", options.LinkStats.MatchedLinks)
		cs.Logger.Info("Not matched links", "not_matched_links", options.LinkStats.NotMatchedLinks)
		// Here, you would add code to populate the 'results' slice with data
	})

	options.Collector.OnRequest(func(r *colly.Request) {
		cs.Logger.Info("Visiting URL", "url", r.URL.String())
	})
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
