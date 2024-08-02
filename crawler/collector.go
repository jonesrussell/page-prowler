package crawler

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/loggo"
)

// CollectorInterface defines the interface for the crawling logic.
type CollectorInterface interface {
	GetCollector() *colly.Collector
	Visit(url string) error
}

// CollectorWrapper is a wrapper around to colly.Collector that implements the CollectorInterface.
type CollectorWrapper struct {
	collector            *colly.Collector
	Logger               loggo.LoggerInterface
	DisallowedURLFilters []*regexp.Regexp
}

var _ CollectorInterface = &CollectorWrapper{}

func NewCollectorWrapper(collector *colly.Collector, logger loggo.LoggerInterface, disallowedURLFilters []*regexp.Regexp) *CollectorWrapper { // Add disallowedURLFilters as a parameter
	// Set a timeout
	collector.WithTransport(&http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second, // Set the timeout
		}).DialContext,
	})

	// Add OnResponse callback
	collector.OnResponse(func(r *colly.Response) {
		contentType := r.Headers.Get("Content-Type")
		logger.Debug(fmt.Sprintf("Content-Type is %s", contentType))
		if !strings.Contains(contentType, "text/html") {
			logger.Debug(fmt.Sprintf("Skipping non-HTML URL: %s", r.Request.URL))
			return
		}
	})

	wrapper := &CollectorWrapper{
		collector:            collector,
		Logger:               logger,
		DisallowedURLFilters: disallowedURLFilters,
	}
	addUserAgentHeader(wrapper.GetCollector(), logger)

	return wrapper
}

// GetCollector implements the CollectorInterface method.
func (cw *CollectorWrapper) GetCollector() *colly.Collector {
	cw.Logger.Debug("Getting the underlying collector")
	return cw.collector
}

// Enhanced Visit method with logging and timing
func (cw *CollectorWrapper) Visit(URL string) error {
	// Check if the URL matches any of the disallowed URL filters
	for _, re := range cw.DisallowedURLFilters {
		if re.MatchString(URL) {
			cw.Logger.Debug(fmt.Sprintf("Skipping disallowed URL: %s", URL))
			return nil // Skip this URL
		}
	}

	start := time.Now()
	cw.Logger.Debug(fmt.Sprintf("Starting to visit: %s", URL)) // Use cw.Logger instead of log
	err := cw.collector.Visit(URL)
	elapsed := time.Since(start)
	cw.Logger.Debug(fmt.Sprintf("Visited: %s in %s", URL, elapsed)) // Use cw.Logger instead of log
	return err
}

// Middleware function to add a User-Agent header
func addUserAgentHeader(c *colly.Collector, logger loggo.LoggerInterface) { // Add logger as a parameter
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
		logger.Debug(fmt.Sprintf("Visiting: %s", r.URL.String())) // Use logger instead of log
	})
}
