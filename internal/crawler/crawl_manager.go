package crawler

import (
	"errors"
	"sync"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/mongodbwrapper"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/jonesrussell/page-prowler/internal/stats"
)

type CrawlManager struct {
	LoggerField    logger.Logger
	Client         prowlredis.ClientInterface
	MongoDBWrapper mongodbwrapper.MongoDBInterface
	Collector      *colly.Collector
	CrawlingMu     *sync.Mutex
	StatsManager   *StatsManager
}

// NewCrawlManager creates a new instance of CrawlManager.
func NewCrawlManager(
	loggerField logger.Logger,
	client prowlredis.ClientInterface,
	mongoDBWrapper mongodbwrapper.MongoDBInterface,
) *CrawlManager {
	return &CrawlManager{
		LoggerField:    loggerField,
		Client:         client,
		MongoDBWrapper: mongoDBWrapper,
	}
}

type StatsManager struct {
	LinkStats   *stats.Stats
	LinkStatsMu sync.RWMutex
}

// NewStatsManager creates a new StatsManager with initialized fields.
func NewStatsManager() *StatsManager {
	return &StatsManager{
		LinkStats:   &stats.Stats{},
		LinkStatsMu: sync.RWMutex{},
	}
}

func (cs *CrawlManager) ConfigureCollector(allowedDomains []string, maxDepth int) error {
	cs.Collector = colly.NewCollector(
		colly.Async(false),
		colly.MaxDepth(maxDepth),
	)

	cs.LoggerField.Debug("Allowed Domains", "domains", allowedDomains)
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
	cs.LoggerField.Debug("[visitWithColly] Visiting URL with Colly", "url", url)

	err := cs.Collector.Visit(url)
	if err != nil {
		switch {
		case errors.Is(err, colly.ErrAlreadyVisited):
			cs.LoggerField.Debug("[visitWithColly] URL already visited", "url", url)
		case errors.Is(err, colly.ErrForbiddenDomain):
			cs.LoggerField.Debug("[visitWithColly] Forbidden domain - Skipping visit", "url", url)
		default:
			cs.LoggerField.Debug("[visitWithColly] Error visiting URL", "url", url, "error", err)
		}
		return nil
	}

	cs.LoggerField.Debug("[visitWithColly] Successfully visited URL with Colly", "url", url)
	return nil
}
