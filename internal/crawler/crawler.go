package crawler

import (
	"sync"
	"time"

	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/internal/stats"
)

const (
	DefaultParallelism = 2
	DefaultDelay       = 3000 * time.Millisecond
)

func (cm *CrawlManager) GetLogger() loggo.LoggerInterface {
	return cm.Logger
}

func (cm *CrawlManager) initializeStatsManager() {
	cm.StatsManager = &StatsManager{
		LinkStats:   &stats.Stats{},
		LinkStatsMu: sync.RWMutex{},
	}
	cm.CrawlingMu.Lock()
	defer cm.CrawlingMu.Unlock()
}
