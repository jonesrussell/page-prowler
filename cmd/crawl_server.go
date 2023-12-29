package cmd

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/stats"
	"github.com/jonesrussell/page-prowler/redis"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// CrawlServer represents the server that handles the crawling process.
type CrawlServer struct {
	CrawlManager *crawler.CrawlManager
}

// PostArticlesStart starts the article posting process.
func (s *CrawlServer) PostArticlesStart(ctx echo.Context) error {
	var req PostArticlesStartJSONBody
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	// Ensure the URL is not empty
	if *req.URL == "" {
		return errors.New("URL cannot be empty")
	}

	// Initialize your Redis client here
	var redisClient redis.Datastore
	if testing.Testing() {
		redisClient = &mockRedisClient{}
	} else {
		var err error
		redisClient, err = redis.NewClient(viper.GetString("REDIS_HOST"), viper.GetString("REDIS_AUTH"), viper.GetString("REDIS_PORT"))
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	// Initialize your crawler service here
	crawlerService, err := initializeManager(ctx.Request().Context(), *req.Debug, redisClient)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to initialize CrawlManager or its Client"})
	}

	// Ensure the URL is correctly formatted
	url := strings.TrimSpace(*req.URL)
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	err = StartCrawling(ctx.Request().Context(), url, *req.SearchTerms, *req.CrawlSiteID, *req.MaxDepth, *req.Debug, crawlerService, s)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Crawling started successfully"})
}

// GetPing handles the ping request.
func (s *CrawlServer) GetPing(ctx echo.Context) error {
	// Implement your logic here
	return ctx.String(http.StatusOK, "Pong")
}

// StartCrawling starts the crawling process.
func StartCrawling(ctx context.Context, url, searchTerms, crawlSiteID string, maxDepth int, debug bool, crawlerService *crawler.CrawlManager, server *CrawlServer) error {
	splitSearchTerms := strings.Split(searchTerms, ",")
	host, err := getHostFromURL(url, crawlerService.Logger)
	if err != nil {
		crawlerService.Logger.Error("Failed to parse URL", "url", url, "error", err)
		return err
	}

	collector := configureCollector([]string{host}, maxDepth)
	if collector == nil {
		crawlerService.Logger.Fatal("Failed to configure collector")
		return errors.New("failed to configure collector")
	}

	var results []crawler.PageData

	options := crawler.CrawlOptions{
		CrawlSiteID: crawlSiteID,
		Collector:   collector,
		SearchTerms: splitSearchTerms,
		Results:     &results,
		LinkStats:   stats.NewStats(),
		Debug:       debug,
	}

	results, err = crawlerService.Crawl(ctx, url, &options)
	if err != nil {
		return err
	}

	err = server.saveResultsToRedis(ctx, results, crawlSiteID)
	if err != nil {
		return err
	}

	printResults(crawlerService, results)

	return nil
}

// getHostFromURL extracts the host from a given URL string.
func getHostFromURL(inputURL string, appLogger logger.Logger) (string, error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		appLogger.Fatal("Failed to parse URL", "url", inputURL, "error", err)
		return "", err // return an empty string and the error
	}

	return u.Host, nil // return the host and nil for the error
}

// configureCollector initializes a new gocolly collector with the specified domains and depth.
func configureCollector(allowedDomains []string, maxDepth int) *colly.Collector {
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
