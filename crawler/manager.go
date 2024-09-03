package crawler

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/gocolly/redisstorage"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/dbmanager"
	"github.com/jonesrussell/page-prowler/internal/termmatcher"
	"github.com/jonesrussell/page-prowler/utils"
)

type CrawlManagerInterface interface {
	Crawl() error
	GetDBManager() dbmanager.DatabaseManagerInterface
	GetLogger() loggo.LoggerInterface
	SetOptions(options *CrawlOptions) error
	UpdateStats(options *CrawlOptions, matchingTerms []string)
}

type CrawlManager struct {
	CollectorInstance *CollectorWrapper
	CrawlingMu        *sync.Mutex
	DBManager         dbmanager.DatabaseManagerInterface
	Logger            loggo.LoggerInterface
	Options           *CrawlOptions
	Results           *Results
	StatsManager      *StatsManager
	Storage           *redisstorage.Storage
	TermMatcher       *termmatcher.TermMatcher
}

var _ CrawlManagerInterface = &CrawlManager{}

func NewCrawlManager(
	logger loggo.LoggerInterface,
	dbManager dbmanager.DatabaseManagerInterface,
	collectorInstance *CollectorWrapper,
	options *CrawlOptions,
	storage *redisstorage.Storage,
	termMatcher *termmatcher.TermMatcher,
) *CrawlManager {
	return &CrawlManager{
		Logger:            logger,
		DBManager:         dbManager,
		CollectorInstance: collectorInstance,
		CrawlingMu:        &sync.Mutex{},
		Options:           options,
		Results:           NewResults(),
		Storage:           storage,
		TermMatcher:       termMatcher,
		StatsManager:      NewStatsManager(),
	}
}

func (cm *CrawlManager) Crawl() error {
	cm.Logger.Info("[Crawl] Starting Crawl function")

	options := cm.GetOptions()
	startURL := options.StartURL

	cm.initializeStatsManager()

	host, err := utils.GetHostFromURL(startURL)
	if err != nil {
		return fmt.Errorf("failed to get host from URL: %w", err)
	}

	if err := cm.configureCollector([]string{host}, options.MaxDepth); err != nil {
		return fmt.Errorf("failed to configure collector: %w", err)
	}

	q, err := cm.createQueue()
	if err != nil {
		return fmt.Errorf("failed to create queue: %w", err)
	}

	if err := q.AddURL(startURL); err != nil {
		return fmt.Errorf("failed to add URL to queue: %w", err)
	}

	if err := q.Run(cm.CollectorInstance.GetCollector()); err != nil {
		return fmt.Errorf("failed to run queue: %w", err)
	}

	defer cm.Storage.Client.Close()

	cm.Logger.Info("[Crawl] Crawling completed.")
	return nil
}

func (cm *CrawlManager) createQueue() (*queue.Queue, error) {
	return queue.New(2, cm.Storage)
}

func (cm *CrawlManager) configureCollector(allowedDomains []string, maxDepth int) error {
	cm.Logger.Debug("[configureCollector]", "maxDepth", maxDepth)

	collector := cm.CollectorInstance.GetCollector()
	collector.AllowedDomains = allowedDomains
	cm.Logger.Info("Allowed domains: ", "whitelist", allowedDomains)

	collector.AllowURLRevisit = false
	collector.Async = false
	collector.IgnoreRobotsTxt = false
	collector.MaxDepth = maxDepth

	limitRule := &colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: DefaultParallelism,
		Delay:       DefaultDelay,
	}

	if err := collector.Limit(limitRule); err != nil {
		return fmt.Errorf("failed to set limit rule: %w", err)
	}

	collector.OnHTML("a[href]", cm.handleLink)
	collector.OnError(cm.handleError)

	return nil
}

func (cm *CrawlManager) handleLink(e *colly.HTMLElement) {
	href, err := cm.getHref(e)
	if err != nil || href == "" {
		return
	}

	matchingTerms, err := cm.TermMatcher.GetMatchingTerms(href, e.Text, cm.Options.SearchTerms)
	if err != nil {
		cm.Logger.Error("Failed to get matching terms", err)
		return
	}

	if len(matchingTerms) > 0 {
		pageData := cm.createPageData(href)
		if err := cm.handleMatchingTerms(cm.Options, e.Request.URL.String(), pageData, matchingTerms); err != nil {
			cm.Logger.Error("Failed to handle matching terms", err)
		}
	}

	if err := e.Request.Visit(href); err != nil {
		cm.Logger.Error("Failed to visit URL", err)
	}
}
func (cm *CrawlManager) handleError(r *colly.Response, err error) {
	cm.Logger.Error("Request failed",
		err,
		"url", r.Request.URL.String(),
		"statusCode", r.StatusCode,
		"body", string(r.Body),
	)
}

func (cm *CrawlManager) GetDBManager() dbmanager.DatabaseManagerInterface {
	return cm.DBManager
}
