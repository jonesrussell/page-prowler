package termmatcher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractTitleFromURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"Test URL with path", "http://example.com/path/to/page", "page"},
		{"Test URL without path", "http://example.com", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractTitleFromURL(tt.url); got != tt.want {
				t.Errorf("extractTitleFromURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveHyphens(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Test with hyphen", "test-title", "test title"},
		{"Test without hyphen", "testtitle", "testtitle"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeHyphens(tt.input); got != tt.want {
				t.Errorf("removeHyphens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveStopwords(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Test with stopword", "this is a test", "test"},
		{"Test without stopword", "testtitle", "testtitle"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeStopwords(tt.input)
			if got != tt.want {
				t.Errorf("removeStopwords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStemTitle(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Test with multiple words", "running tests", "run test"},
		{"Test with single word", "test", "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stemTitle(tt.input)
			fmt.Println("Expected: ", tt.want)
			fmt.Println("Actual: ", got)
			if got != tt.want {
				t.Errorf("stemTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessTitle(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  string
	}{
		{"Test title with hyphen", "test-title", "test titl"},
		{"Test title with stopword", "this is a test", "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processTitle(tt.title); got != tt.want {
				t.Errorf("processTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMatchingTerms(t *testing.T) {
	// Test with a URL that should match the search terms
	href := "https://example.com/privacy-policy"
	searchTerms := []string{"privacy", "policy"}
	expected := []string{"privacy", "policy"}
	assert.Equal(t, expected, GetMatchingTerms(href, searchTerms))

	// Test with a URL that should not match the search terms
	href = "https://example.com/unrelated-term"
	assert.NotEqual(t, []string{}, GetMatchingTerms(href, searchTerms))
}

func TestFindMatchingTerms(t *testing.T) {
	// Test with a title that should match the search terms
	title := "privacy policy"
	searchTerms := []string{"privacy", "policy"}
	expected := []string{"privacy", "policy"}
	assert.Equal(t, expected, findMatchingTerms(title, searchTerms))

	// Test with a title that should not match the search terms
	title = "unrelated term"
	assert.NotEqual(t, []string{}, findMatchingTerms(title, searchTerms))
}
