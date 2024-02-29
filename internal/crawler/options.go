package crawler

func NewCrawlOptions(crawlSiteID string, searchTerms []string, debug bool, results *[]PageData) *CrawlOptions {
	return &CrawlOptions{
		CrawlSiteID: crawlSiteID,
		SearchTerms: searchTerms,
		Results:     results,
		Debug:       debug,
	}
}
