// Package crawler provides the tools and logic needed to perform web crawling and data extraction.
package crawler

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/crawler/internal/crawlresult"
	"github.com/jonesrussell/crawler/internal/rediswrapper"
	"github.com/jonesrussell/crawler/internal/stats"
	"github.com/jonesrussell/crawler/internal/termmatcher"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CrawlerService encapsulates shared dependencies for crawler functions.
type CrawlerService struct {
	Logger       *zap.SugaredLogger
	RedisWrapper *rediswrapper.RedisWrapper
}

// InitializeLogger initializes the logger used in the application.
func InitializeLogger(debug bool) (*zap.SugaredLogger, error) {
	var logger *zap.Logger
	var err error

	if debug {
		// Development logger is more verbose and writes to standard output
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // optional, colorizes the output
		logger, err = config.Build()
	} else {
		// Production logger is less verbose and could be set to log to a file
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	return logger.Sugar(), nil
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
func (cs *CrawlerService) HandleHTMLParsing(
	ctx context.Context,
	crawlSiteID string,
	collector *colly.Collector,
	searchTerms []string,
	linkStats *stats.Stats,
	results *[]crawlresult.PageData,
) error {
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Request.AbsoluteURL(e.Attr("href"))
		linkStats.IncrementTotalLinks()

		pageData := crawlresult.PageData{
			URL: href,
			// Add other fields as necessary
		}

		if termmatcher.Related(href, searchTerms) {
			linkStats.IncrementMatchedLinks()
			if err := cs.handleMatchingLinks(ctx, collector, href, crawlSiteID); err != nil {
				cs.Logger.Error("Error handling matching links", zap.Error(err))
			}
			pageData.MatchingTerms = searchTerms
		} else {
			linkStats.IncrementNotMatchedLinks()
			cs.handleNonMatchingLinks(href)
		}

		*results = append(*results, pageData)
	})
	return nil
}

// handleMatchingLinks is responsible for handling the links that match the search criteria during crawling.
func (cs *CrawlerService) handleMatchingLinks(
	ctx context.Context,
	collector *colly.Collector,
	href string,
	crawlSiteID string,
) error {
	cs.Logger.Info("Found: ", zap.String("url", href))
	if _, err := cs.RedisWrapper.SAdd(ctx, crawlSiteID, href); err != nil {
		cs.Logger.Error("Error adding URL to Redis set", zap.String("set", crawlSiteID), zap.Error(err))
		return err
	}

	err := collector.Visit(href)
	if err != nil {
		if err == colly.ErrAlreadyVisited {
			cs.Logger.Info("URL already visited", zap.String("url", href))
			return nil
		}
		cs.Logger.Error("Error visiting URL", zap.String("url", href), zap.Error(err))
		return err
	}

	return nil
}

// handleNonMatchingLinks logs the occurrence of a non-matching link.
func (cs *CrawlerService) handleNonMatchingLinks(href string) {
	cs.Logger.Info("Non-matching link: ", zap.String("url", href))
}

// handleErrorEvents sets up the error handling for the colly collector.
func (cs *CrawlerService) handleErrorEvents(collector *colly.Collector) {
	collector.OnError(func(r *colly.Response, err error) {
		statusCode := r.StatusCode
		requestURL := r.Request.URL.String()

		if statusCode != 404 {
			cs.Logger.Error("Request URL failed",
				zap.String("request_url", requestURL),
				zap.Int("status_code", statusCode),
				zap.Error(err),
			)
		}
	})
}

// handleRedisOperations manages the Redis operations after crawling a page.
func (cs *CrawlerService) handleRedisOperations(ctx context.Context) error {
	// You need to pass the context and the appropriate key to SMembers
	hrefs, err := cs.RedisWrapper.SMembers(ctx, "yourKeyHere") // Replace "yourKeyHere" with the actual key you're interested in
	if err != nil {
		cs.Logger.Error("Error getting members from Redis", zap.Error(err))
		return err
	}

	for _, href := range hrefs {
		err = cs.RedisWrapper.PublishHref(ctx, "streetcode", href) // Make sure to pass ctx to PublishHref if it requires it
		if err != nil {
			cs.Logger.Error("Error publishing href to Redis", zap.Error(err))
			return err
		}

		if _, err = cs.RedisWrapper.Del(ctx, href); err != nil { // Make sure to pass ctx to Del if it requires it
			cs.Logger.Error("Error deleting from Redis", zap.Error(err))
			return err
		}
	}
	return nil
}

// SetupCrawlingLogic configures and initiates the crawling logic.
func (cs *CrawlerService) SetupCrawlingLogic(
	ctx context.Context,
	crawlSiteID string,
	collector *colly.Collector,
	searchTerms []string,
	results *[]crawlresult.PageData,
) {
	linkStats := stats.NewStats()

	err := cs.HandleHTMLParsing(ctx, crawlSiteID, collector, searchTerms, linkStats, results)
	if err != nil {
		cs.Logger.Error("Error during HTML parsing", zap.Error(err))
		return
	}

	cs.handleErrorEvents(collector)

	collector.OnScraped(func(r *colly.Response) {
		err := cs.handleRedisOperations(ctx)
		if err != nil {
			cs.Logger.Error("Error with Redis operations", zap.Error(err))
			return
		}

		cs.Logger.Info("Finished scraping the page", zap.String("url", r.Request.URL.String()))
		cs.Logger.Info("Total links found", zap.Int("total_links", linkStats.TotalLinks))
		cs.Logger.Info("Matched links", zap.Int("matched_links", linkStats.MatchedLinks))
		cs.Logger.Info("Not matched links", zap.Int("not_matched_links", linkStats.NotMatchedLinks))
		// Here, you would add code to populate the 'results' slice with data
	})

	collector.OnRequest(func(r *colly.Request) {
		cs.Logger.Info("Visiting", zap.String("url", r.URL.String()))
	})
}

// GetHostFromURL extracts the host from a given URL string.
func GetHostFromURL(inputURL string) string {
	u, err := url.Parse(inputURL)
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
	}
	return u.Host
}
