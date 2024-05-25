package crawler

import "time"

// CrawlOptions represents the configuration for a crawl.
type CrawlOptions struct {
	CrawlSiteID           string
	Debug                 bool
	DelayBetweenRequests  time.Duration
	MaxConcurrentRequests int
	MaxDepth              int
	SearchTerms           []string
	StartURL              string
}

// Method to set options
func (cm *CrawlManager) SetOptions(options *CrawlOptions) error {
	cm.Options = options
	return nil
}

// Method to get options
func (cm *CrawlManager) GetOptions() *CrawlOptions {
	return cm.Options
}
