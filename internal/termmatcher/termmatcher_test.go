package termmatcher

import (
	"math"
	"reflect"
	"testing"

	"github.com/jonesrussell/loggo"
	"github.com/stretchr/testify/assert"
)

type MockContentProcessor struct{}

func (m *MockContentProcessor) Process(content string) string {
	return content
}

func (m *MockContentProcessor) Stem(content string) string {
	return content
}

func TestGetMatchingTerms(t *testing.T) {
	mockLogger := loggo.NewMockLogger()
	processor := &MockContentProcessor{}
	tm := NewTermMatcher(mockLogger, 0.5, processor) // Increase the threshold to 0.5

	tests := []struct {
		name        string
		href        string
		anchorText  string
		searchTerms []string
		want        []string
		wantErr     bool
	}{
		{
			name:        "Test case 1: Matching terms",
			href:        "https://example.com/running-shoes",
			anchorText:  "Best Running Shoes",
			searchTerms: []string{"run", "shoe"},
			want:        []string{"run", "shoe"},
			wantErr:     false,
		},
		{
			name:        "Test case 2: No matching terms",
			href:        "https://example.com/laptops",
			anchorText:  "Best Laptops",
			searchTerms: []string{"run", "shoe"},
			want:        []string{},
			wantErr:     false,
		},
		{
			name:        "Test case 3: Short combined content",
			href:        "https://example.com/a",
			anchorText:  "A",
			searchTerms: []string{"a"},
			want:        []string{},
			wantErr:     false,
		},
		{
			name:        "Test case 4: Phrase matching",
			href:        "https://example.com/gang-activities",
			anchorText:  "Recent Gang Activities",
			searchTerms: []string{"gang activities"},
			want:        []string{"gang activities"},
			wantErr:     false,
		},
		{
			name:        "Test case 5: Similarity matching",
			href:        "https://example.com/organized-crime",
			anchorText:  "Organized Criminal Activities",
			searchTerms: []string{"organized crime"},
			want:        []string{"organized crime"},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger.DebugMessages = nil // Clear debug messages before each test
			got, err := tm.GetMatchingTerms(tt.href, tt.anchorText, tt.searchTerms)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMatchingTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMatchingTerms() = %v, want %v", got, tt.want)
				t.Logf("Debug logs:")
				for _, log := range mockLogger.DebugMessages {
					t.Logf(log)
				}
			}
		})
	}
}

func TestTermMatcher_CompareTerms(t *testing.T) {
	logger := loggo.NewMockLogger()
	processor := &MockContentProcessor{}
	tm := NewTermMatcher(logger, 0.6, processor)

	tests := []struct {
		name       string
		searchTerm string
		content    string
		want       float64
	}{
		{
			name:       "High similarity",
			searchTerm: "running",
			content:    "run",
			want:       0.77,
		},
		{
			name:       "Low similarity",
			searchTerm: "laptop",
			content:    "computer",
			want:       0.0,
		},
		{
			name:       "Exact match",
			searchTerm: "book",
			content:    "book",
			want:       1.0,
		},
		{
			name:       "Case insensitive",
			searchTerm: "HELLO",
			content:    "hello",
			want:       1.0,
		},
		{
			name:       "Empty strings",
			searchTerm: "",
			content:    "",
			want:       0.0,
		},
		{
			name:       "Just below threshold",
			searchTerm: "organize",
			content:    "organ",
			want:       0.0, // Assuming threshold is 0.6
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tm.CompareTerms(tt.searchTerm, tt.content)
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("TermMatcher.CompareTerms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTermMatcher_findMatchingTerms(t *testing.T) {
	logger := loggo.NewMockLogger()
	processor := &MockContentProcessor{}
	tm := NewTermMatcher(logger, 0.8, processor) // Increase threshold to 0.8

	tests := []struct {
		name        string
		content     string
		searchTerms []string
		want        []string
	}{
		{
			name:        "Test case 1: Exact matching terms",
			content:     "running shoes for athletes",
			searchTerms: []string{"run", "shoe", "athlete"},
			want:        []string{"run", "shoe", "athlete"},
		},
		{
			name:        "Test case 2: No matching terms",
			content:     "laptop computers for sale",
			searchTerms: []string{"run", "shoe", "athlete"},
			want:        []string{},
		},
		{
			name:        "Test case 3: Phrase matching",
			content:     "gang activities in the city",
			searchTerms: []string{"gang activities", "city crime"},
			want:        []string{"gang activities"},
		},
		{
			name:        "Test case 4: Similarity matching",
			content:     "organized criminal activities",
			searchTerms: []string{"organized crime", "gang activities"},
			want:        []string{"organized crime"},
		},
		{
			name:        "Test case 5: Similarity matching with threshold",
			content:     "organize criminal activities",
			searchTerms: []string{"organized crime", "gang activities"},
			want:        []string{"organized crime"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tm.findMatchingTerms(tt.content, tt.searchTerms)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TermMatcher.findMatchingTerms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTermMatcher_flattenSearchTerms(t *testing.T) {
	logger := loggo.NewMockLogger()
	processor := &MockContentProcessor{}
	tm := NewTermMatcher(logger, 0.8, processor)

	tests := []struct {
		name        string
		searchTerms []string
		want        []string
	}{
		{
			name:        "Test case 1: Single terms",
			searchTerms: []string{"run", "shoe", "athlete"},
			want:        []string{"run", "shoe", "athlete"},
		},
		{
			name:        "Test case 2: Comma-separated terms",
			searchTerms: []string{"run,shoe", "athlete,sport"},
			want:        []string{"run", "shoe", "athlete", "sport"},
		},
		{
			name:        "Test case 3: Mixed single and comma-separated terms",
			searchTerms: []string{"run", "shoe,athlete", "sport"},
			want:        []string{"run", "shoe", "athlete", "sport"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tm.flattenSearchTerms(tt.searchTerms)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTermMatcher_combineContents(t *testing.T) {
	logger := loggo.NewMockLogger()
	processor := &MockContentProcessor{}
	tm := NewTermMatcher(logger, 0.8, processor)

	tests := []struct {
		name     string
		content1 string
		content2 string
		want     string
	}{
		{
			name:     "Test case 1: Both contents non-empty",
			content1: "running",
			content2: "shoes",
			want:     "running shoes",
		},
		{
			name:     "Test case 2: Second content empty",
			content1: "running",
			content2: "",
			want:     "running",
		},
		{
			name:     "Test case 3: First content empty",
			content1: "",
			content2: "shoes",
			want:     " shoes",
		},
		{
			name:     "Test case 4: Both contents empty",
			content1: "",
			content2: "",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tm.combineContents(tt.content1, tt.content2)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTermMatcher_removeDuplicates(t *testing.T) {
	logger := loggo.NewMockLogger()
	processor := &MockContentProcessor{}
	tm := NewTermMatcher(logger, 0.8, processor)

	tests := []struct {
		name  string
		terms []string
		want  []string
	}{
		{
			name:  "No duplicates",
			terms: []string{"apple", "banana", "cherry"},
			want:  []string{"apple", "banana", "cherry"},
		},
		{
			name:  "With duplicates",
			terms: []string{"apple", "banana", "apple", "cherry", "banana"},
			want:  []string{"apple", "banana", "cherry"},
		},
		{
			name:  "All duplicates",
			terms: []string{"apple", "apple", "apple"},
			want:  []string{"apple"},
		},
		{
			name:  "Empty slice",
			terms: []string{},
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tm.removeDuplicates(tt.terms)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestTermMatcher_compareSingleTerm(t *testing.T) {
	logger := loggo.NewMockLogger()
	processor := &MockContentProcessor{}
	tm := NewTermMatcher(logger, 0.6, processor)

	tests := []struct {
		name  string
		term  string
		words []string
		want  []string
	}{
		{
			name:  "Matching term",
			term:  "run",
			words: []string{"running", "shoes", "for", "athletes"},
			want:  []string{"run"},
		},
		{
			name:  "No matching term",
			term:  "book",
			words: []string{"running", "shoes", "for", "athletes"},
			want:  []string{},
		},
		{
			name:  "Multiple similar words",
			term:  "run",
			words: []string{"running", "runner", "ran", "jog"},
			want:  []string{"run"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tm.compareSingleTerm(tt.term, tt.words)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTermMatcher_compareMultiTerm(t *testing.T) {
	logger := loggo.NewMockLogger()
	processor := &MockContentProcessor{}
	tm := NewTermMatcher(logger, 0.6, processor)

	tests := []struct {
		name  string
		term  string
		words []string
		want  []string
	}{
		{
			name:  "Exact match",
			term:  "running shoes",
			words: []string{"best", "running", "shoes", "for", "athletes"},
			want:  []string{"running shoes"},
		},
		{
			name:  "No match",
			term:  "tennis racket",
			words: []string{"best", "running", "shoes", "for", "athletes"},
			want:  []string{},
		},
		{
			name:  "Partial match",
			term:  "running gear",
			words: []string{"best", "running", "shoes", "and", "gear"},
			want:  []string{"running gear"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tm.compareMultiTerm(tt.term, tt.words)
			assert.Equal(t, tt.want, got)
		})
	}
}
