package crawler

import (
	"context"
	"errors"
	"sync"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/mongodbwrapper"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/jonesrussell/page-prowler/internal/stats"
)

// NewCrawlManager creates a new instance of CrawlManager.
func NewCrawlManager(
	logger logger.Logger,
	client prowlredis.ClientInterface,
	mongoDBWrapper mongodbwrapper.MongoDBInterface,
) *CrawlManager {
	cm := &CrawlManager{
		LoggerField:    logger,
		Client:         client,
		MongoDBWrapper: mongoDBWrapper,
		Collector:      colly.NewCollector(),
		VisitedPages:   make(map[string]bool),
	}
	cm.MatchedLinkProcessor = &ConcreteMatchedLinkProcessor{CrawlManager: cm}
	return cm
}

type StatsManager struct {
	LinkStats   *stats.Stats
	LinkStatsMu sync.RWMutex
}

func (cs *CrawlManager) StartCrawling(ctx context.Context, url, searchTerms, crawlSiteID string, maxDepth int, debug bool) error {
	// Initialize LinkStats...
	cs.StatsManager = &StatsManager{
		LinkStats:   &stats.Stats{},
		LinkStatsMu: sync.RWMutex{},
	}
	cs.CrawlingMu.Lock()
	defer cs.CrawlingMu.Unlock()

	host, err := GetHostFromURL(url, cs.Logger())
	if err != nil {
		cs.Error("Failed to parse URL", "url", url, "error", err)
		return err
	}

	cs.Debug("Extracted host from URL", "host", host)

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

func (cs *CrawlManager) ConfigureCollector(allowedDomains []string, maxDepth int) error {
	cs.Collector = colly.NewCollector(
		colly.Async(false),
		colly.MaxDepth(maxDepth),
	)

	cs.Debug("Allowed Domains", "domains", allowedDomains)
	cs.Collector.AllowedDomains = allowedDomains

	limitRule := cs.createLimitRule()
	if err := cs.Collector.Limit(limitRule); err != nil {
		cs.Logger().Errorf("Failed to set limit rule: %v", err)
		return err
	}

	// Respect robots.txt
	cs.Collector.AllowURLRevisit = false
	cs.Collector.IgnoreRobotsTxt = false

	// Register OnScraped callback
	cs.Collector.OnScraped(func(r *colly.Response) {
		cs.Logger().Debug("[OnScraped] Page scraped", "url", r.Request.URL)
		cs.StatsManager.LinkStatsMu.Lock()
		defer cs.StatsManager.LinkStatsMu.Unlock()
		cs.StatsManager.LinkStats.IncrementTotalPages()
	})

	return nil
}

func (cs *CrawlManager) logCrawlingStatistics(options *CrawlOptions) {
	report := cs.StatsManager.LinkStats.Report()
	cs.Logger().Info("Crawling statistics",
		"TotalLinks", report["TotalLinks"],
		"MatchedLinks", report["MatchedLinks"],
		"NotMatchedLinks", report["NotMatchedLinks"],
		"TotalPages", report["TotalPages"],
	)
}

func (cs *CrawlManager) visitWithColly(url string) error {
	cs.Debug("[visitWithColly] Visiting URL with Colly", "url", url)

	err := cs.Collector.Visit(url)
	if err != nil {
		switch {
		case errors.Is(err, colly.ErrAlreadyVisited):
			cs.Debug("[visitWithColly] URL already visited", "url", url)
		case errors.Is(err, colly.ErrForbiddenDomain):
			cs.Debug("[visitWithColly] Forbidden domain - Skipping visit", "url", url)
		default:
			cs.Debug("[visitWithColly] Error visiting URL", "url", url, "error", err)
		}
		return nil
	}

	cs.Debug("[visitWithColly] Successfully visited URL with Colly", "url", url)
	return nil
}
