// stats.go
package stats

type Stats struct {
	TotalLinks      int
	MatchedLinks    int
	NotMatchedLinks int
	// Add new stats here as needed
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) IncrementTotalLinks() {
	s.TotalLinks++
}

func (s *Stats) IncrementMatchedLinks() {
	s.MatchedLinks++
}

func (s *Stats) IncrementNotMatchedLinks() {
	s.NotMatchedLinks++
}

// Add more methods as needed to manipulate and access your stats
