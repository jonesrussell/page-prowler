package crawler

import (
	"testing"
)

func TestNewCrawlOptions(t *testing.T) {
	debug := true
	var results []PageData

	co := NewCrawlOptions(debug, &results)

	if co.Debug != debug {
		t.Errorf("Expected Debug to be %v, got %v", debug, co.Debug)
	}
}
