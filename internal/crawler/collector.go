package crawler

import (
	"log"
	"time"

	"github.com/gocolly/colly"
)

// CollectorInterface defines the interface for the crawling logic.
type CollectorInterface interface {
	Visit(url string) error
	OnRequest(requestFunc func(*colly.Request))
	OnHTML(selector string, htmlFunc func(*colly.HTMLElement))
	OnError(func(r *colly.Response, err error))
	OnScraped(callback func(*colly.Response))
	Wait()
	Limit(limitRule colly.LimitRule) error
	SetAllowedDomains([]string)
	AllowURLRevisit() bool
	SetAllowURLRevisit(allow bool)
	IgnoreRobotsTxt() bool
	SetIgnoreRobotsTxt(bool)

	// GetUnderlyingCollector returns the underlying *colly.Collector instance.
	GetUnderlyingCollector() *colly.Collector
}

// CollectorWrapper is a wrapper around to colly.Collector that implements the CollectorInterface.
type CollectorWrapper struct {
	collector *colly.Collector
}

// Modify NewCollectorWrapper to apply middleware
func NewCollectorWrapper(collector *colly.Collector) *CollectorWrapper {
	wrapper := &CollectorWrapper{collector: collector}
	addUserAgentHeader(wrapper.GetUnderlyingCollector())
	return wrapper
}

// GetUnderlyingCollector implements the CollectorInterface method.
func (cw *CollectorWrapper) GetUnderlyingCollector() *colly.Collector {
	return cw.collector
}

// Enhanced Visit method with logging and timing
func (cw *CollectorWrapper) Visit(URL string) error {
	startTime := time.Now()
	err := cw.collector.Visit(URL)
	log.Printf("Visited %s in %v", URL, time.Since(startTime))
	return err
}

// Example of a middleware function to add a User-Agent header
func addUserAgentHeader(c *colly.Collector) {
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "My Custom User Agent")
	})
}

// OnHTML is a wrapper method that delegates to the underlying colly.Collector.
func (cw *CollectorWrapper) OnHTML(selector string, cb func(*colly.HTMLElement)) {
	cw.collector.OnHTML(selector, cb)
}

func (cw *CollectorWrapper) OnError(callback func(r *colly.Response, err error)) {
	cw.collector.OnError(callback)
}

func (cw *CollectorWrapper) Wait() {
	cw.collector.Wait()
}

func (cw *CollectorWrapper) Limit() error {
	return nil
}

// AllowURLRevisit implements CollectorInterface.
func (cw *CollectorWrapper) AllowURLRevisit() bool {
	return cw.collector.AllowURLRevisit
}

// SetAllowedDomains Implement other methods as needed
func (cw *CollectorWrapper) SetAllowedDomains(domains []string) {
	cw.collector.AllowedDomains = domains
}

// IgnoreRobotsTxt implements CollectorInterface.
func (cw *CollectorWrapper) IgnoreRobotsTxt() bool {
	return cw.collector.IgnoreRobotsTxt
}

// SetAllowURLRevisit SetIgnoreRobotsTxt implements CollectorInterface.
func (cw *CollectorWrapper) SetAllowURLRevisit(allow bool) {
	cw.collector.AllowURLRevisit = allow
}

// SetIgnoreRobotsTxt implements CollectorInterface.
func (cw *CollectorWrapper) SetIgnoreRobotsTxt(ignore bool) {
	cw.collector.IgnoreRobotsTxt = ignore
}
