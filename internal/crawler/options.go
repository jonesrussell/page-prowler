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
