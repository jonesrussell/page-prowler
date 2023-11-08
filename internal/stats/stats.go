// Package stats provides a simple way to track and manipulate statistics related to web crawling.
package stats

// Stats holds counters for various metrics related to web crawling.
type Stats struct {
	TotalLinks      int // TotalLinks is the total number of links processed.
	MatchedLinks    int // MatchedLinks is the count of links that matched the search criteria.
	NotMatchedLinks int // NotMatchedLinks is the count of links that did not match the search criteria.
}

// NewStats creates and returns a new Stats instance with all counters set to zero.
func NewStats() *Stats {
	return &Stats{}
}

// IncrementTotalLinks increases the TotalLinks counter by one.
func (s *Stats) IncrementTotalLinks() {
	s.TotalLinks++
}

// IncrementMatchedLinks increases the MatchedLinks counter by one.
func (s *Stats) IncrementMatchedLinks() {
	s.MatchedLinks++
}

// IncrementNotMatchedLinks increases the NotMatchedLinks counter by one.
func (s *Stats) IncrementNotMatchedLinks() {
	s.NotMatchedLinks++
}
