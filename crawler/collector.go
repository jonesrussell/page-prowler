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
	GetCollector() *colly.Collector
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
