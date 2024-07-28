package crawler

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/internal/stats"
)

const (
	DefaultParallelism = 2
	DefaultDelay       = 3000 * time.Millisecond
)

type LoggerDebugger struct {
	loggo.LoggerInterface
	debug.Debugger
}

var _ loggo.LoggerInterface = &LoggerDebugger{}
var _ debug.Debugger = &LoggerDebugger{}

func (cm *CrawlManager) Logger() loggo.LoggerInterface {
	return cm.LoggerField
}

func (cm *CrawlManager) initializeStatsManager() {
	cm.StatsManager = &StatsManager{
		LinkStats:   &stats.Stats{},
		LinkStatsMu: sync.RWMutex{},
	}
	cm.CrawlingMu.Lock()
	defer cm.CrawlingMu.Unlock()
}

func (cm *CrawlManager) SetupErrorEventHandler() {
	cm.CollectorInstance.OnError(func(r *colly.Response, err error) {
		statusCode := r.StatusCode
		requestURL := r.Request.URL.String()

		if statusCode == 500 {
			cm.LoggerField.Debug(fmt.Sprintf("[SetupErrorEventHandler] Internal Server Error request_url: %s, status_code: %d, error: %v", requestURL, statusCode, err))
		} else if statusCode != 404 {
			cm.LoggerField.Debug(fmt.Sprintf("[SetupErrorEventHandler] Request URL failed request_url: %s, status_code: %d, error: %v", requestURL, statusCode, err))
		}
	})
}

func (cm *CrawlManager) CrawlURL(url string) error {
	// Check if CrawlManager or LoggerField is nil
	if cm == nil || cm.LoggerField == nil {
		fmt.Println("Error: CrawlManager or LoggerField is nil")
		return errors.New("CrawlManager or LoggerField is nil")
	}

	if cm.LoggerField.IsDebugEnabled() {
		cm.LoggerField.Debug(fmt.Sprintf("[CrawlURL] Debug is enabled. Visiting URL: %v", url))
	}

	err := cm.visitWithColly(url)
	if err != nil {
		cm.LoggerField.Error(fmt.Sprintf("[CrawlURL] Error visiting URL: %v", url), err)
		return cm.HandleVisitError(url, err)
	}

	cm.CollectorInstance.Wait()

	cm.Logger().Info("[CrawlURL] Crawling completed.")
	return nil
}
