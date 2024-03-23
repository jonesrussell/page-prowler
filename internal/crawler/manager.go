package crawler

import (
	"errors"
	"fmt"
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
		colly.Debugger(cm.LoggerField),
	)

	cm.LoggerField.Debug(fmt.Sprintf("Allowed Domains: %v", allowedDomains))
	cm.Collector.AllowedDomains = allowedDomains

	limitRule := cm.createLimitRule()
	if err := cm.Collector.Limit(limitRule); err != nil {
		cm.LoggerField.Error(fmt.Sprintf("Failed to set limit rule: %v", err))
		return err
	}

	// Respect robots.txt
	cm.Collector.AllowURLRevisit = false
	cm.Collector.IgnoreRobotsTxt = false

	// Register OnScraped callback
	cm.Collector.OnScraped(func(r *colly.Response) {
		cm.LoggerField.Debug(fmt.Sprintf("[OnScraped] Page scraped: %s", r.Request.URL.String()))
		cm.StatsManager.LinkStatsMu.Lock()
		defer cm.StatsManager.LinkStatsMu.Unlock()
		cm.StatsManager.LinkStats.IncrementTotalPages()
	})

	return nil
}

func (cm *CrawlManager) logCrawlingStatistics() {
	report := cm.StatsManager.LinkStats.Report()
	infoMessage := fmt.Sprintf("Crawling statistics: TotalLinks=%v, MatchedLinks=%v, NotMatchedLinks=%v, TotalPages=%v",
		report["TotalLinks"], report["MatchedLinks"], report["NotMatchedLinks"], report["TotalPages"])
	cm.LoggerField.Info(infoMessage)
}

func (cm *CrawlManager) visitWithColly(url string) error {
	cm.LoggerField.Debug(fmt.Sprintf("[visitWithColly] Visiting URL with Colly: %v", url))

	err := cm.Collector.Visit(url)
  if err != nil {
    switch {
    case errors.Is(err, colly.ErrAlreadyVisited):
      errorMessage := fmt.Sprintf("[visitWithColly] URL already visited: %v", url)
      cm.LoggerField.Debug(errorMessage)
    case errors.Is(err, colly.ErrForbiddenDomain):
      errorMessage := fmt.Sprintf("[visitWithColly] Forbidden domain - Skipping visit: %v", url)
      cm.LoggerField.Debug(errorMessage)
    default:
      errorMessage := fmt.Sprintf("[visitWithColly] Error visiting URL: url=%v, error=%v", url, err)
      cm.LoggerField.Error(errorMessage)
    }
    return nil
  }

	successMessage := fmt.Sprintf("[visitWithColly] Successfully visited URL: %v", url)
  cm.LoggerField.Debug(successMessage)
	return nil
}
