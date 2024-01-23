package crawler

import "github.com/jonesrussell/page-prowler/internal/stats"

func NewCrawlOptions(crawlSiteID string, searchTerms []string, debug bool, results *[]PageData, linkStats *stats.Stats) *CrawlOptions {
	return &CrawlOptions{
		CrawlSiteID: crawlSiteID,
		SearchTerms: searchTerms,
		Results:     results,
		LinkStats:   linkStats,
		Debug:       debug,
	}
}
