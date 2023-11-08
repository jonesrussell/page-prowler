package stats

import (
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
