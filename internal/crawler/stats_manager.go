package crawler

import (
	"sync"

	"github.com/jonesrussell/page-prowler/internal/stats"
)

// StatsManager is a struct that manages crawling statistics.
// It includes fields for link statistics and a mutex for thread safety.
type StatsManager struct {
	LinkStats   *stats.Stats
	LinkStatsMu sync.RWMutex
}

// NewStatsManager creates a new StatsManager with initialized fields.
func NewStatsManager() *StatsManager {
	return &StatsManager{
		LinkStats:   &stats.Stats{},
		LinkStatsMu: sync.RWMutex{},
	}
}
