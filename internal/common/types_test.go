package common_test

import (
	"testing"

	"github.com/jonesrussell/page-prowler/internal/common"
)

func TestCrawlManagerKey(t *testing.T) {
	if common.CrawlManagerKey == nil {
		t.Errorf("Expected CrawlManagerKey to be assigned, but it was nil")
	}
}
