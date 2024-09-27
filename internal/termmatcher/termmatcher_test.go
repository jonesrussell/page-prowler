package termmatcher

import (
	"reflect"
	"strings"
	"testing"

	"github.com/adrg/strutil/metrics"
	"github.com/golang/mock/gomock"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/internal/matcher"
)

type fields struct {
	logger loggo.LoggerInterface
	swg    *metrics.SmithWatermanGotoh
}

type args struct {
	content string
}

// MockMatcher is a mock implementation of the Matcher interface for testing.
type MockMatcher struct{}

// Match implements the Matcher interface.
func (mm *MockMatcher) Match(content string, pattern string) (bool, error) {
	if content == "" || pattern == "" {
		return false, nil
	}
	return strings.Contains(content, pattern), nil
}

func TestNewTermMatcher(t *testing.T) {
	logger := loggo.NewMockLogger(gomock.NewController(t))
	mockMatcher := &MockMatcher{}
	mockMatchers := []matcher.Matcher{mockMatcher}

	tm := NewTermMatcher(logger, mockMatchers)

	if tm == nil {
		t.Error("Expected TermMatcher to be initialized, got nil")
	}
}

func TestGetMatchingTerms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		mockMatchers []matcher.Matcher
		href         string
		anchorText   string
		searchTerms  []string
		want         []string
	}{
		{
			name:         "Test case 1",
			mockMatchers: []matcher.Matcher{&MockMatcher{}},
			href:         "https://example.com/test",
			anchorText:   "Example Anchor Text",
			searchTerms:  []string{"example", "test"},
			want:         []string{"test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := loggo.NewMockLoggerInterface(ctrl)
			mockLogger.EXPECT().Debug(gomock.Any()).AnyTimes()

			tm := NewTermMatcher(mockLogger, tt.mockMatchers)
			got := tm.GetMatchingTerms(tt.href, tt.anchorText, tt.searchTerms)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMatchingTerms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTermMatcher_compareAndAppendTerm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := loggo.NewMockLogger(ctrl)
	swg := metrics.NewSmithWatermanGotoh()

	type args struct {
		searchTerm string
		content    string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Test case 1: Exact match",
			fields: fields{
				logger: logger,
				swg:    swg,
			},
			args: args{
				searchTerm: "hello",
				content:    "hello world",
			},
			want: true,
		},
		{
			name: "Test case 2: No match",
			fields: fields{
				logger: logger,
				swg:    swg,
			},
			args: args{
				searchTerm: "goodbye",
				content:    "hello world",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := loggo.NewMockLoggerInterface(ctrl)
			mockLogger.EXPECT().Debug(gomock.Any()).AnyTimes()

			tm := &TermMatcher{
				logger: mockLogger,
				swg:    tt.fields.swg,
			}
			if got := tm.compareAndAppendTerm(tt.args.searchTerm, tt.args.content); got != tt.want {
				t.Errorf("TermMatcher.compareAndAppendTerm() = %v, want %v", got, tt.want)
			}
		})
	}
}
