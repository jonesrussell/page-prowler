package mining

import (
	"testing"
)

func TestMatcher_Match(t *testing.T) {
	matcher := NewMatcher()

	tests := []struct {
		href     string
		expected bool
	}{
		{"http://example.com/mining-news", true},
		{"http://example.com/gold-prices", true},
		{"http://example.com/silver-market", true},
		{"http://example.com/technology", false},
		{"http://example.com/coal-extraction", true},
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
