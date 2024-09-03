package termmatcher

import (
	"reflect"
	"testing"

	"github.com/jonesrussell/loggo"
)

func TestGetMatchingTerms(t *testing.T) {
	logger := loggo.NewMockLogger()
	tm := NewTermMatcher(logger, 0.8)

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
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tm.GetMatchingTerms(tt.href, tt.anchorText, tt.searchTerms)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMatchingTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMatchingTerms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTermMatcher_processContent(t *testing.T) {
	logger := loggo.NewMockLogger()
	tm := NewTermMatcher(logger, 0.8)

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "Test case 1: Process and stem content",
			content: "Running shoes are the best",
			want:    "run shoe best",
		},
		{
			name:    "Test case 2: Remove hyphens and stopwords",
			content: "The quick-brown fox jumps over the lazy dog",
			want:    "quick brown fox jump lazi dog",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tm.processContent(tt.content); got != tt.want {
				t.Errorf("TermMatcher.processContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTermMatcher_CompareTerms(t *testing.T) {
	logger := loggo.NewMockLogger()
	tm := NewTermMatcher(logger, 0.8)

	tests := []struct {
		name       string
		searchTerm string
		content    string
		want       float64
	}{
		{
			name:       "Test case 1: High similarity",
			searchTerm: "running",
			content:    "run",
			want:       1.0,
		},
		{
			name:       "Test case 2: Low similarity",
			searchTerm: "laptop",
			content:    "computer",
			want:       0.31666666666666665,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tm.CompareTerms(tt.searchTerm, tt.content)
			if got != tt.want {
				t.Errorf("TermMatcher.CompareTerms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTermMatcher_findMatchingTerms(t *testing.T) {
	logger := loggo.NewMockLogger()
	tm := NewTermMatcher(logger, 0.8)

	tests := []struct {
		name        string
		content     string
		searchTerms []string
		want        []string
	}{
		{
			name:        "Test case 1: Matching terms",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tm.findMatchingTerms(tt.content, tt.searchTerms)
			if len(got) != len(tt.want) {
				t.Errorf("TermMatcher.findMatchingTerms() = %v, want %v", got, tt.want)
			}
		})
	}
}
