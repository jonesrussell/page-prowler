// Package crawlresult defines the data structures used for storing
// results of a web crawl operation.
package crawlresult

import "time"

// Metadata contains metadata information of a crawled web page.
type Metadata struct {
	Description string   `json:"description"` // The description of the web page
	Keywords    []string `json:"keywords"`    // The keywords associated with the web page
}

// Content holds the main content extracted from a crawled web page.
type Content struct {
	Title string `json:"title"` // The title of the web page
	Body  string `json:"body"`  // The body content of the web page
}

// PageData encapsulates all relevant information gathered from a crawled web page.
type PageData struct {
	URL           string    `json:"url"`             // The URL of the web page
	CrawlTime     time.Time `json:"crawl_time"`      // The timestamp when the crawl was performed
	StatusCode    int       `json:"status_code"`     // The HTTP status code received for the web page
	Metadata      *Metadata `json:"metadata"`        // Metadata associated with the web page
	Content       *Content  `json:"content"`         // Main content extracted from the web page
	Links         []string  `json:"links"`           // The hyperlinks found on the web page
	SearchTerms   []string  `json:"search_terms"`    // The search terms used during the crawl
	MatchingTerms []string  `json:"matching_terms"`  // The terms that matched the search criteria
	Error         string    `json:"error,omitempty"` // Any error encountered during crawling of this page
}
