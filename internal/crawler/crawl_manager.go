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
		CrawlingMu:     &sync.Mutex{},
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

func (cm *CrawlManager) ConfigureCollector(allowedDomains []string, maxDepth int) error {
	cm.Collector = colly.NewCollector(
		colly.Async(false),
		colly.MaxDepth(maxDepth),
	)

	cm.LoggerField.Debug("Allowed Domains", "domains", allowedDomains)
	cm.Collector.AllowedDomains = allowedDomains

	limitRule := cm.createLimitRule()
	if err := cm.Collector.Limit(limitRule); err != nil {
		cm.Logger().Errorf("Failed to set limit rule: %v", err)
		return err
	}

	// Respect robots.txt
	cm.Collector.AllowURLRevisit = false
	cm.Collector.IgnoreRobotsTxt = false

	// Register OnScraped callback
	cm.Collector.OnScraped(func(r *colly.Response) {
		cm.Logger().Debug("[OnScraped] Page scraped", "url", r.Request.URL)
		cm.StatsManager.LinkStatsMu.Lock()
		defer cm.StatsManager.LinkStatsMu.Unlock()
		cm.StatsManager.LinkStats.IncrementTotalPages()
	})

	return nil
}

func (cm *CrawlManager) logCrawlingStatistics() {
	report := cm.StatsManager.LinkStats.Report()
	cm.Logger().Info("Crawling statistics",
		"TotalLinks", report["TotalLinks"],
		"MatchedLinks", report["MatchedLinks"],
		"NotMatchedLinks", report["NotMatchedLinks"],
		"TotalPages", report["TotalPages"],
	)
}

func (cm *CrawlManager) visitWithColly(url string) error {
	cm.LoggerField.Debug("[visitWithColly] Visiting URL with Colly", "url", url)

	err := cm.Collector.Visit(url)
	if err != nil {
		switch {
		case errors.Is(err, colly.ErrAlreadyVisited):
			cm.LoggerField.Debug("[visitWithColly] URL already visited", "url", url)
		case errors.Is(err, colly.ErrForbiddenDomain):
			cm.LoggerField.Debug("[visitWithColly] Forbidden domain - Skipping visit", "url", url)
		default:
			cm.LoggerField.Debug("[visitWithColly] Error visiting URL", "url", url, "error", err)
		}
		return nil
	}

	cm.LoggerField.Debug("[visitWithColly] Successfully visited URL with Colly", "url", url)
	return nil
}
