// Package stats provides a simple way to track and manipulate statistics related to web crawling.
package stats

import "sync"

// Stats holds counters for various metrics related to web crawling.
type Stats struct {
	TotalLinks      int
	MatchedLinks    int
	NotMatchedLinks int
	TotalPages      int
	Links           []string
	mu              sync.Mutex
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

// IncrementTotalPages increases the TotalPages counter by one.
func (s *Stats) IncrementTotalPages() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalPages++
}

// GetMatchedLinks retrieves the total number of not matched links.
func (s *Stats) GetMatchedLinks() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.MatchedLinks
}

// GetNotMatchedLinks retrieves the total number of not matched links.
func (s *Stats) GetNotMatchedLinks() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.NotMatchedLinks
}

// GetTotalPages retrieves the total number of pages crawled.
func (s *Stats) GetTotalPages() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.TotalPages
}

func (s *Stats) Report() map[string]interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return map[string]interface{}{
		"TotalLinks":      s.TotalLinks,
		"MatchedLinks":    s.MatchedLinks,
		"NotMatchedLinks": s.NotMatchedLinks,
		"TotalPages":      s.TotalPages,
	}
}
