package termmatcher_test

import (
	"fmt"
	"testing"

	"github.com/adrg/strutil/metrics"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/internal/termmatcher"
	"github.com/jonesrussell/page-prowler/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

type LoggerInterface interface {
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
	SetLevel(zapcore.Level)
	AllEntries() []zapcore.Entry
}

func TestGetMatchingTerms(t *testing.T) {
	tests := []struct {
		name        string
		href        string
		anchorText  string
		searchTerms []string
		expected    []string
	}{
		{
			name:        "Test with a URL and anchor text that should match the search terms",
			href:        "https://example.com/privacy-policy",
			anchorText:  "Privacy Policy",
			searchTerms: []string{"privacy", "policy"},
			expected:    []string{"privaci", "polici"},
		},
		// Add more test cases as needed...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := loggo.NewLogger("/path/to/logfile") // Create a new logger
			if err != nil {
				t.Fatalf("Failed to create logger: %v", err)
			}
			assert.Equal(t, tt.expected, termmatcher.GetMatchingTerms(tt.href, tt.anchorText, tt.searchTerms, logger))
		})
	}
}

func TestExtractLastSegmentFromURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"Test URL with path", "https://example.com/path/to/page", "page"},
		{"Test URL without path", "https://example.com", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := termmatcher.ExtractLastSegmentFromURL(tt.url); got != tt.want {
				t.Errorf("ExtractLastSegmentFromURL() = %v, want %v", got, tt.want)
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
			if got := termmatcher.RemoveHyphens(tt.input); got != tt.want {
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
			got := termmatcher.RemoveStopwords(tt.input)
			if got != tt.want {
				t.Errorf("RemoveStopwords() = %v, want %v", got, tt.want)
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
			got := termmatcher.ProcessAndStem(tt.input)
			fmt.Println("Expected: ", tt.want)
			fmt.Println("Actual: ", got)
			if got != tt.want {
				t.Errorf("procesAndStem() = %v, want %v", got, tt.want)
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
			if got := termmatcher.ProcessContent(tt.title); got != tt.want {
				t.Errorf("processTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareAndAppendTerm(t *testing.T) {
	swg := termmatcher.CreateSWG()
	logger := mocks.NewMockLogger()
	var matchingTerms []string

	termmatcher.CompareAndAppendTerm("test", "test", swg, &matchingTerms, logger)

	if len(matchingTerms) != 1 {
		t.Errorf("Expected matchingTerms to have 1 element, got %v", len(matchingTerms))
	}

	if matchingTerms[0] != "test" {
		t.Errorf("Expected first element of matchingTerms to be 'test', got %v", matchingTerms[0])
	}
}

func TestCreateSWG(t *testing.T) {
	swg := termmatcher.CreateSWG()

	if swg.CaseSensitive != false {
		t.Errorf("Expected CaseSensitive to be false, got %v", swg.CaseSensitive)
	}

	if swg.GapPenalty != -0.1 {
		t.Errorf("Expected GapPenalty to be -0.1, got %v", swg.GapPenalty)
	}

	matchMismatch, ok := swg.Substitution.(metrics.MatchMismatch)
	if !ok {
		t.Fatalf("Unexpected type for Substitution: %T", swg.Substitution)
	}
	if matchMismatch.Match != 1 {
		t.Errorf("Expected Substitution.Match to be 1, got %v", matchMismatch.Match)
	}

	if matchMismatch.Mismatch != -0.5 {
		t.Errorf("Expected Substitution.Mismatch to be -0.5, got %v", matchMismatch.Mismatch)
	}
}

func TestFindMatchingTerms(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		searchTerms []string
		expected    []string
	}{
		{
			name:        "Test with matching terms",
			content:     "privacy policy",
			searchTerms: []string{"privacy", "policy"},
			expected:    []string{"privaci", "polici"},
		},
		{
			name:        "Test with no matching terms",
			content:     "unrelated term",
			searchTerms: []string{"privacy", "policy"},
			expected:    []string{},
		},
		{
			name:        "Test with term appearing in content",
			content:     "trump win iowa caucus crucial victori outset republican presidenti campaign trump win iowa caucus crucial victori outset republican presidenti campaign",
			searchTerms: []string{"fight", "prescription", "gang", "drug", "JOINT", "CANNABI", "IMPAIR", "SHOOT", "FIREARM", "MURDER", "COCAIN", "POSSESS"},
			expected:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := mocks.NewMockLogger()
			mockLogger.SetLevel(zapcore.DebugLevel) // Set the logger level to Debug
			fmt.Println("Content: ", tt.content)
			fmt.Println("Search Terms: ", tt.searchTerms)
			actual := termmatcher.FindMatchingTerms(tt.content, tt.searchTerms, mockLogger)
			assert.Equal(t, tt.expected, actual)

			// Print the logs with human-readable similarity scores
			for _, entry := range mockLogger.AllEntries() {
				fmt.Printf("Message: %s, Fields: %v\n", entry.Message, entry.Context)
			}

		})
	}
}

func TestCombineContents(t *testing.T) {
	tests := []struct {
		name     string
		content1 string
		content2 string
		expected string
	}{
		{
			name:     "Test with two strings",
			content1: "Hello",
			content2: "World",
			expected: "Hello World",
		},
		{
			name:     "Test with empty second string",
			content1: "Hello",
			content2: "",
			expected: "Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, termmatcher.CombineContents(tt.content1, tt.content2))
		})
	}
}

func TestConvertToLowercase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Test with uppercase letters", "HELLO WORLD", "hello world"},
		{"Test with lowercase letters", "hello world", "hello world"},
		{"Test with mixed case letters", "HeLlO WoRlD", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := termmatcher.ConvertToLowercase(tt.input); got != tt.want {
				t.Errorf("ConvertToLowercase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStemContent(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Update these expected values based on the actual behavior of your stemmer
		{"Test with multiple words", "running tests", "run test"},
		{"Test with single word", "test", "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := termmatcher.StemContent(tt.input); got != tt.want {
				t.Errorf("StemContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareTerms(t *testing.T) {
	tests := []struct {
		name       string
		searchTerm string
		content    string
		want       float64
	}{
		{"Test with similar terms", "test", "testing", 1},
		{"Test with dissimilar terms", "apple", "banana", 0.2},
	}

	swg := metrics.NewSmithWatermanGotoh()
	swg.CaseSensitive = false
	swg.GapPenalty = -0.1
	swg.Substitution = metrics.MatchMismatch{
		Match:    1,
		Mismatch: -0.5,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := mocks.NewMockLogger()
			mockLogger.SetLevel(zapcore.DebugLevel) // Set the logger level to Debug

			if got := termmatcher.CompareTerms(tt.searchTerm, tt.content, swg, mockLogger); got != tt.want {
				t.Errorf("CompareTerms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func generateLargeNumberOfUniqueSearchTerms(n int) []string {
	terms := make([]string, n)
	for i := 0; i < n; i++ {
		terms[i] = fmt.Sprintf("term%d", i)
	}
	return terms
}

func Test_CombineContents(t *testing.T) {
	type args struct {
		content1 string
		content2 string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test case 1",
			args: args{
				content1: "Hello",
				content2: "World",
			},
			want: "Hello World",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := termmatcher.CombineContents(tt.args.content1, tt.args.content2); got != tt.want {
				t.Errorf("CombineContents() = %v, want %v", got, tt.want)
			}
		})
	}
}
