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

func TestGetMatchingTerms(t *testing.T) {
	// Test with a URL that should match the search terms
	href := "https://example.com/privacy-policy"
	searchTerms := []string{"privacy", "policy"}
	expected := []string{"PRIVACI", "POLICI"}
	assert.Equal(t, expected, GetMatchingTerms(href, searchTerms))

	// Test with a URL that should not match the search terms
	href = "https://example.com/unrelated-term"
	assert.Equal(t, []string{}, GetMatchingTerms(href, searchTerms))
}

func TestFindMatchingTerms(t *testing.T) {
	// Test with a title that should match the search terms
	title := "privacy policy"
	searchTerms := []string{"privacy", "policy"}
	expected := []string{"PRIVACI", "POLICI"}
	assert.Equal(t, expected, findMatchingTerms(title, searchTerms))

	// Test with a title that should not match the search terms
	title = "unrelated term"
	assert.Equal(t, []string{}, findMatchingTerms(title, searchTerms))
}
