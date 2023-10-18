package crawlResult

import "time"

type Metadata struct {
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
}

type Content struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type PageData struct {
	URL           string    `json:"url"`
	CrawlTime     time.Time `json:"crawl_time"`
	StatusCode    int       `json:"status_code"`
	Metadata      *Metadata `json:"metadata"`
	Content       *Content  `json:"content"`
	Links         []string  `json:"links"`
	SearchTerms   []string  `json:"search_terms"`
	MatchingTerms []string  `json:"matching_terms"`
	Error         string    `json:"error,omitempty"` // omitempty means it won't appear in the JSON if it's empty
}
