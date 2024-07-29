package crawler

import (
	"context"
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
	Logger() loggo.LoggerInterface
	ProcessMatchingLink(currentURL string, pageData models.PageData, matchingTerms []string) error
	SaveResultsToRedis(ctx context.Context, results []models.PageData, key string) error
	SetOptions(options *CrawlOptions) error
	UpdateStats(options *CrawlOptions, matchingTerms []string)
}

var _ CrawlManagerInterface = &CrawlManager{}

type CrawlManager struct {
	CollectorInstance *CollectorWrapper
	CrawlingMu        *sync.Mutex
	DBManager         dbmanager.DatabaseManagerInterface
	LoggerField       loggo.LoggerInterface
	Options           *CrawlOptions
	Results           *Results
	StatsManager      *StatsManager
	TermMatcher       *termmatcher.TermMatcher
}

func NewCrawlManager(
	loggerField loggo.LoggerInterface,
	dbManager dbmanager.DatabaseManagerInterface,
	collectorInstance *CollectorWrapper,
	options *CrawlOptions,
) *CrawlManager {
	return &CrawlManager{
		LoggerField:       loggerField,
		DBManager:         dbManager,
		CollectorInstance: collectorInstance,
		CrawlingMu:        &sync.Mutex{},
		Options:           options,
		Results:           NewResults(),
		TermMatcher:       termmatcher.NewTermMatcher(loggerField),
	}
}

func (cm *CrawlManager) Crawl() error {
	cm.LoggerField.Info("[Crawl] Starting Crawl function")

	options := cm.GetOptions()

	startURL := options.StartURL

	cm.initializeStatsManager()

	host, err := utils.GetHostFromURL(startURL)
	if err != nil {
		return err
	}

	if err := cm.ConfigureCollector([]string{host}, options.MaxDepth); err != nil {
		return err
	}

	if err := cm.CollectorInstance.Visit(startURL); err != nil {
		return err
	}

	cm.LoggerField.Info("[Crawl] Crawling completed.")

	return nil
}

func (cm *CrawlManager) ConfigureCollector(allowedDomains []string, maxDepth int) error {
	collector := colly.NewCollector(
		colly.Async(false),
		colly.MaxDepth(maxDepth),
	)

	collector.AllowedDomains = allowedDomains

	limitRule := &colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: DefaultParallelism,
		Delay:       DefaultDelay,
	}

	if err := collector.Limit(limitRule); err != nil {
		return err
	}

	collector.AllowURLRevisit = false
	collector.IgnoreRobotsTxt = false

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

	cm.CollectorInstance = &CollectorWrapper{collector}

	return nil
}

func (cm *CrawlManager) AppendResult(pageData models.PageData) {
	if cm.Results == nil || cm.Results.Pages == nil {
		return
	}
	cm.Results.Pages = append(cm.Results.Pages, pageData)
}

func (cm *CrawlManager) GetResults() *Results {
	return cm.Results
}

func (cm *CrawlManager) SaveResultsToRedis(ctx context.Context, results []models.PageData, key string) error {
	return cm.DBManager.SaveResultsToRedis(ctx, results, key)
}

func (cm *CrawlManager) GetDBManager() dbmanager.DatabaseManagerInterface {
	return cm.DBManager
}
