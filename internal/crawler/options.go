package crawler

// CrawlOptions represents the options for configuring and initiating the crawling logic.
// It includes parameters that control the crawling process, such as the site ID to crawl, search terms to match, a pointer to store results, and a debug flag.
type CrawlOptions struct {
	CrawlSiteID string      // The ID of the site to crawl.
	SearchTerms []string    // The search terms to match against the crawled content.
	Results     *[]PageData // A pointer to a slice of PageData where the crawling results will be stored.
	Debug       bool        // A flag indicating whether to enable debug mode for the crawling process.
}

// NewCrawlOptions creates a new CrawlOptions instance with the given parameters.
// Parameters:
// - crawlSiteID: The ID of the site to crawl.
// - searchTerms: The search terms to match against the crawled content.
// - debug: A flag indicating whether to enable debug mode for the crawling process.
// - results: A pointer to a slice of PageData where the crawling results will be stored.
// Returns:
// - *CrawlOptions: A pointer to a new CrawlOptions instance configured with the provided parameters.
func NewCrawlOptions(crawlSiteID string, searchTerms []string, debug bool, results *[]PageData) *CrawlOptions {
	return &CrawlOptions{
		CrawlSiteID: crawlSiteID,
		SearchTerms: searchTerms,
		Results:     results,
		Debug:       debug,
	}
}
