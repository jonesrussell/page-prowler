package crawler

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/dbmanager"
	"github.com/jonesrussell/page-prowler/internal/termmatcher"
	"github.com/jonesrussell/page-prowler/models"
	"github.com/jonesrussell/page-prowler/utils"
)

type CrawlManagerInterface interface {
	Crawl() error
	GetDBManager() dbmanager.DatabaseManagerInterface
	GetLogger() loggo.LoggerInterface
	ProcessMatchingLink(currentURL string, pageData models.PageData, matchingTerms []string) error
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
	TermMatcher       *termmatcher.TermMatcher
}

var _ CrawlManagerInterface = &CrawlManager{}

func NewCrawlManager(
	logger loggo.LoggerInterface,
	dbManager dbmanager.DatabaseManagerInterface,
	collectorInstance *CollectorWrapper,
	options *CrawlOptions,
) *CrawlManager {
	return &CrawlManager{
		Logger:            logger,
		DBManager:         dbManager,
		CollectorInstance: collectorInstance,
		CrawlingMu:        &sync.Mutex{},
		Options:           options,
		Results:           NewResults(),
		TermMatcher:       termmatcher.NewTermMatcher(logger),
	}
}

func (cm *CrawlManager) Crawl() error {
	cm.Logger.Info("[Crawl] Starting Crawl function")

	options := cm.GetOptions()

	startURL := options.StartURL

	cm.initializeStatsManager()

	host, err := utils.GetHostFromURL(startURL)
	if err != nil {
		return err
	}

	if err := cm.configureCollector([]string{host}, options.MaxDepth); err != nil {
		return err
	}

	if err := cm.CollectorInstance.Visit(startURL); err != nil {
		return err
	}

	cm.Logger.Info("[Crawl] Crawling completed.")

	return nil
}

func (cm *CrawlManager) configureCollector(allowedDomains []string, maxDepth int) error {
	cm.Logger.Debug("[configureCollector]", "maxDepth", maxDepth)

	// Get the underlying colly.Collector from the CollectorWrapper
	collector := cm.CollectorInstance.GetCollector()

	collector.AllowedDomains = allowedDomains
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
		return err
	}

	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href, err := cm.getHref(e)
		if err != nil {
			return
		}
		if href == "" {
			return
		}

		err = cm.processLink(e, href)
		if err != nil {
			return
		}

		err = cm.CollectorInstance.Visit(href)
		if err != nil {
			return
		}
	})

	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	return nil
}

func (cm *CrawlManager) GetDBManager() dbmanager.DatabaseManagerInterface {
	return cm.DBManager
}
