package crawler

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/gocolly/redisstorage"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/dbmanager"
	"github.com/jonesrussell/page-prowler/internal/matcher"
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
	TermMatcher       *termmatcher.TermMatcher // Ensure TermMatcher is included
}

var _ CrawlManagerInterface = &CrawlManager{}

func NewCrawlManager(
	logger loggo.LoggerInterface,
	dbManager dbmanager.DatabaseManagerInterface,
	collectorInstance *CollectorWrapper,
	options *CrawlOptions,
	storage *redisstorage.Storage,
) *CrawlManager {
	return &CrawlManager{
		Logger:            logger,
		DBManager:         dbManager,
		CollectorInstance: collectorInstance,
		CrawlingMu:        &sync.Mutex{},
		Options:           options,
		Results:           NewResults(),
		Storage:           storage,
		TermMatcher:       termmatcher.NewTermMatcher(logger, []matcher.Matcher{}), // Initialize TermMatcher with empty matcher slice
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
		return err
	}

	cm.Logger.Debug("options", "MaxDepth", options.MaxDepth)
	if err := cm.configureCollector([]string{host}, options.MaxDepth); err != nil {
		return err
	}

	// Create a new request queue with the Redis storage backend
	q, err := queue.New(2, cm.Storage)
	if err != nil {
		return fmt.Errorf("failed to create queue: %v", err)
	}

	// Add the start URL to the queue
	err = q.AddURL(startURL)
	if err != nil {
		return fmt.Errorf("failed to add URL to queue: %v", err)
	}

	// Consume requests
	err = q.Run(cm.CollectorInstance.GetCollector())
	if err != nil {
		return fmt.Errorf("failed to run queue: %v", err)
	}

	// close redis client
	defer cm.Storage.Client.Close()

	cm.Logger.Info("[Crawl] Crawling completed.")

	return nil
}

func (cm *CrawlManager) configureCollector(allowedDomains []string, maxDepth int) error {
	cm.Logger.Debug("[configureCollector]", "maxDepth", maxDepth)

	// Get the underlying colly.Collector from the CollectorWrapper
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

		// Use TermMatcher to find matching terms in the URL and anchor text
		matchingTerms := cm.TermMatcher.GetMatchingTerms(href, e.Text, cm.Options.SearchTerms)
		if len(matchingTerms) > 0 {
			pageData := cm.createPageData(href)
			err := cm.handleMatchingTerms(cm.Options, e.Request.URL.String(), pageData, matchingTerms)
			if err != nil {
				return
			}
		}

		err = e.Request.Visit(href)
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
