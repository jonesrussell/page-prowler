package crawler

import (
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// CollectorInterface defines the interface for the crawling logic.
type CollectorInterface interface {
	AllowURLRevisit() bool
	GetCollector() *colly.Collector
	IgnoreRobotsTxt() bool
	SetAllowedDomains([]string)
	SetAllowURLRevisit(allow bool)
	SetIgnoreRobotsTxt(bool)
	Visit(url string) error
}

// CollectorWrapper is a wrapper around to colly.Collector that implements the CollectorInterface.
type CollectorWrapper struct {
	collector *colly.Collector
}

// Modify NewCollectorWrapper to apply middleware
func NewCollectorWrapper(collector *colly.Collector) *CollectorWrapper {
	// Set a timeout
	collector.WithTransport(&http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second, // Set the timeout
		}).DialContext,
	})

	// Enable asynchronous operation
	collector.Async = true

	// Add OnResponse callback
	collector.OnResponse(func(r *colly.Response) {
		contentType := r.Headers.Get("Content-Type")
		log.Printf("Content-Type is %s", contentType)
		if !strings.Contains(contentType, "text/html") {
			log.Printf("Skipping non-HTML URL: %s", r.Request.URL)
			return
		}
	})

	wrapper := &CollectorWrapper{collector: collector}
	addUserAgentHeader(wrapper.GetCollector())

	return wrapper
}

// GetCollector implements the CollectorInterface method.
func (cw *CollectorWrapper) GetCollector() *colly.Collector {
	log.Println("Getting the underlying collector")
	return cw.collector
}

// Enhanced Visit method with logging and timing
func (cw *CollectorWrapper) Visit(URL string) error {
	start := time.Now()
	log.Printf("Starting to visit: %s", URL)
	err := cw.collector.Visit(URL)
	elapsed := time.Since(start)
	log.Printf("Visited: %s in %s", URL, elapsed)
	return err
}

// Example of a middleware function to add a User-Agent header
func addUserAgentHeader(c *colly.Collector) {
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
		log.Printf("Visiting: %s", r.URL.String())
	})
}

// SetAllowedDomains Implement other methods as needed
func (cw *CollectorWrapper) SetAllowedDomains(domains []string) {
	log.Printf("Setting allowed domains: %v", domains)
	cw.collector.AllowedDomains = domains
}

// IgnoreRobotsTxt implements CollectorInterface.
func (cw *CollectorWrapper) IgnoreRobotsTxt() bool {
	log.Println("Checking if robots.txt is ignored")
	return cw.collector.IgnoreRobotsTxt
}

// SetAllowURLRevisit SetIgnoreRobotsTxt implements CollectorInterface.
func (cw *CollectorWrapper) SetAllowURLRevisit(allow bool) {
	log.Printf("Setting allow URL revisit to: %v", allow)
	cw.collector.AllowURLRevisit = allow
}

// SetIgnoreRobotsTxt implements CollectorInterface.
func (cw *CollectorWrapper) SetIgnoreRobotsTxt(ignore bool) {
	log.Printf("Setting ignore robots.txt to: %v", ignore)
	cw.collector.IgnoreRobotsTxt = ignore
}
