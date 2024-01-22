// Package stats provides a simple way to track and manipulate statistics related to web crawling.
package stats

import "sync"

type fields struct {
	TotalLinks      int
	MatchedLinks    int
	NotMatchedLinks int
	Links           []string
}

// Stats holds counters for various metrics related to web crawling.
type Stats struct {
	fields
	mu sync.Mutex
}

// NewStats creates and returns a new Stats instance with all counters set to zero.
func NewStats() *Stats {
	return &Stats{}
}

// IncrementTotalLinks increases the TotalLinks counter by one.
func (s *Stats) IncrementTotalLinks() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalLinks++
}

// IncrementMatchedLinks increases the MatchedLinks counter by one.
func (s *Stats) IncrementMatchedLinks() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.MatchedLinks++
}

// IncrementNotMatchedLinks increases the NotMatchedLinks counter by one.
func (s *Stats) IncrementNotMatchedLinks() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.NotMatchedLinks++
}

func (s *Stats) Report() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return map[string]int{
		"TotalLinks":      s.TotalLinks,
		"MatchedLinks":    s.MatchedLinks,
		"NotMatchedLinks": s.NotMatchedLinks,
	}
}
