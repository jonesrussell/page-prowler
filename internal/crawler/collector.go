package crawler

import (
	"log"
	"time"

	"github.com/gocolly/colly"
)

// ConfigureCollector initializes a new gocolly collector with the specified domains and depth.
func ConfigureCollector(allowedDomains []string, maxDepth int) (*colly.Collector, error) {
	collector := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(maxDepth),
	)

	collector.AllowedDomains = allowedDomains

	err := collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       3000 * time.Millisecond,
	})
	if err != nil {
		log.Println("Error setting limit rule:", err)
		return nil, err
	}

	// Debugging statement
	log.Println("ConfigureCollector: returning a non-nil collector")

	// Respect robots.txt
	collector.AllowURLRevisit = false
	collector.IgnoreRobotsTxt = false

	return collector, nil
}
