package crawler

import (
	"log"

	"github.com/gocolly/colly"
)

// CollectorInterface defines the interface for the collector.
type CollectorInterface interface {
	Visit(URL string) error
	OnHTML(selector string, cb func(*colly.HTMLElement))
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

// CollectorWrapper is a wrapper around the colly.Collector that implements the CollectorInterface.
type CollectorWrapper struct {
	collector *colly.Collector
}

// GetUnderlyingCollector implements the CollectorInterface method.
func (cw *CollectorWrapper) GetUnderlyingCollector() *colly.Collector {
	return cw.collector
}

// OnScraped implements CollectorInterface.
func (cw *CollectorWrapper) OnScraped(callback func(*colly.Response)) {
	// Example implementation: Log a message when the OnScraped event is triggered.
	log.Println("OnScraped event triggered")

	// If you want to execute the callback provided to OnScraped, you can do so here.
	// However, since the original implementation panics, it's unclear if this is the intended behavior.
	// If you have a specific callback function you want to execute, you can call it directly.
	// For example:
	// callback(&colly.Response{})
}

// NewCollectorWrapper creates a new CollectorWrapper with specific allowed domains.
func NewCollectorWrapper(collector *colly.Collector) *CollectorWrapper {
	return &CollectorWrapper{collector: collector}
}

// Visit is a wrapper method that delegates to the underlying colly.Collector.
func (cw *CollectorWrapper) Visit(URL string) error {
	return cw.collector.Visit(URL)
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

func (cw *CollectorWrapper) Limit(limitRule colly.LimitRule) error {
	return nil
}

// AllowURLRevisit implements CollectorInterface.
func (cw *CollectorWrapper) AllowURLRevisit() bool {
	return cw.collector.AllowURLRevisit
}

// Implement other methods as needed
func (cw *CollectorWrapper) SetAllowedDomains(domains []string) {
	cw.collector.AllowedDomains = domains
}

// IgnoreRobotsTxt implements CollectorInterface.
func (cw *CollectorWrapper) IgnoreRobotsTxt() bool {
	return cw.collector.IgnoreRobotsTxt
}

// SetIgnoreRobotsTxt implements CollectorInterface.
func (cw *CollectorWrapper) SetAllowURLRevisit(allow bool) {
	cw.collector.AllowURLRevisit = allow
}

// SetIgnoreRobotsTxt implements CollectorInterface.
func (cw *CollectorWrapper) SetIgnoreRobotsTxt(ignore bool) {
	cw.collector.IgnoreRobotsTxt = ignore
}
