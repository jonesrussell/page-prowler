package termmatcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractTitleFromURL(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"https://example.com/article-title", "article-title"},
		{"https://example.com", ""},
		{"", ""},
	}

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			title := extractTitleFromURL(test.url)
			assert.Equal(t, test.expected, title)
		})
	}
}

func TestProcessTitle(t *testing.T) {
	tests := []struct {
		title    string
		expected string
	}{
		{"This is a - test title", "TEST TITL"},
		{"Some - example - title", "EXAMP TITL"},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			processedTitle := processTitle(test.title)
			assert.Equal(t, test.expected, processedTitle)
		})
	}
}

func TestMatchSearchTerms(t *testing.T) {
	tests := []struct {
		title       string
		searchTerms []string
		expected    bool
	}{
		{"This is a test title", []string{"test"}, true},
		{"Example title", []string{"example", "test"}, true},
		{"Another title", []string{"word"}, false},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			match := matchSearchTerms(test.title, test.searchTerms)
			assert.Equal(t, test.expected, match)
		})
	}
}
