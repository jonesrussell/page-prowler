package stats

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStats(t *testing.T) {
	// Create a new Stats instance
	stats := NewStats()

	// Check that all counters are initially zero
	assert.Equal(t, 0, stats.TotalLinks)
	assert.Equal(t, 0, stats.MatchedLinks)
	assert.Equal(t, 0, stats.NotMatchedLinks)

	// Increment TotalLinks and check the counter
	stats.IncrementTotalLinks()
	assert.Equal(t, 1, stats.TotalLinks)

	// Increment MatchedLinks and check the counter
	stats.IncrementMatchedLinks()
	assert.Equal(t, 1, stats.MatchedLinks)

	// Increment NotMatchedLinks and check the counter
	stats.IncrementNotMatchedLinks()
	assert.Equal(t, 1, stats.NotMatchedLinks)
}

func TestNewStats(t *testing.T) {
	tests := []struct {
		name string
		want *Stats
	}{
		{
			name: "Test NewStats",
			want: &Stats{
				TotalLinks:      0,
				MatchedLinks:    0,
				NotMatchedLinks: 0,
				Links:           nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStats(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStats() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStats_IncrementTotalLinks(t *testing.T) {
	tests := []struct {
		name       string
		TotalLinks int
	}{
		{
			name:       "Test IncrementTotalLinks",
			TotalLinks: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStats()
			s.TotalLinks = tt.TotalLinks
			s.IncrementTotalLinks()
			assert.Equal(t, 1, s.TotalLinks)
		})
	}
}

func TestStats_IncrementMatchedLinks(t *testing.T) {
	tests := []struct {
		name            string
		TotalLinks      int
		MatchedLinks    int
		NotMatchedLinks int
		Links           []string
	}{
		{
			name:            "Test IncrementMatchedLinks",
			TotalLinks:      0,
			MatchedLinks:    0,
			NotMatchedLinks: 0,
			Links:           nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStats()
			s.TotalLinks = tt.TotalLinks
			s.MatchedLinks = tt.MatchedLinks
			s.NotMatchedLinks = tt.NotMatchedLinks
			s.Links = tt.Links
			s.IncrementMatchedLinks()
			assert.Equal(t, 1, s.MatchedLinks)
		})
	}
}

func TestStats_IncrementNotMatchedLinks(t *testing.T) {
	tests := []struct {
		name            string
		TotalLinks      int
		MatchedLinks    int
		NotMatchedLinks int
		Links           []string
	}{
		{
			name:            "Test IncrementNotMatchedLinks",
			TotalLinks:      0,
			MatchedLinks:    0,
			NotMatchedLinks: 0,
			Links:           nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStats()
			s.TotalLinks = tt.TotalLinks
			s.MatchedLinks = tt.MatchedLinks
			s.NotMatchedLinks = tt.NotMatchedLinks
			s.Links = tt.Links
			s.IncrementNotMatchedLinks()
			assert.Equal(t, 1, s.NotMatchedLinks)
		})
	}
}

func TestGetMatchedLinks(t *testing.T) {
	// Create a new Stats instance
	stats := NewStats()

	// Increment the NotMatchedLinks counter
	stats.IncrementMatchedLinks()

	// Check the result
	if stats.GetMatchedLinks() != 1 {
		t.Errorf("Expected matched links count to be 1, got %d", stats.GetMatchedLinks())
	}
}

func TestGetNotMatchedLinks(t *testing.T) {
	// Create a new Stats instance
	stats := NewStats()

	// Increment the NotMatchedLinks counter
	stats.IncrementNotMatchedLinks()

	// Check the result
	if stats.GetNotMatchedLinks() != 1 {
		t.Errorf("Expected not matched links count to be 1, got %d", stats.GetNotMatchedLinks())
	}
}
