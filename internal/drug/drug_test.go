package drug

import (
	"testing"

	"github.com/adrg/strutil/metrics"
)

func TestMatcher_Match(t *testing.T) {
	swg := metrics.NewSmithWatermanGotoh()
	matcher := NewMatcher(swg)

	tests := []struct {
		href     string
		expected bool
	}{
		{"http://example.com/drug", true},
		{"http://example.com/smoke-joint", true},
		{"http://example.com/recipe", false},
		{"http://example.com/", false},
		{"http://example.com/unknown-title", false},
	}

	for _, test := range tests {
		t.Run(test.href, func(t *testing.T) {
			result := matcher.Match(test.href)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}
