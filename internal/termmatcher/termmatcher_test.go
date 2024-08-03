package termmatcher

import (
	"reflect"
	"testing"

	"github.com/adrg/strutil/metrics"
	"github.com/jonesrussell/loggo"
)

type fields struct {
	logger loggo.LoggerInterface
	swg    *metrics.SmithWatermanGotoh
}

type args struct {
	content string
}

func TestGetMatchingTerms(t *testing.T) {
	logger := loggo.NewMockLogger()
	tm := NewTermMatcher(logger)

	tests := []struct {
		name        string
		href        string
		anchorText  string
		searchTerms []string
		want        []string
	}{
		{
			name:        "Test case 1",
			href:        "https://example.com/test",
			anchorText:  "Example Anchor Text",
			searchTerms: []string{"example", "test"},
			want:        []string{"examp", "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tm.GetMatchingTerms(tt.href, tt.anchorText, tt.searchTerms); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMatchingTerms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTermMatcher_stemContent(t *testing.T) {
	logger := loggo.NewMockLogger()
	swg := metrics.NewSmithWatermanGotoh()

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test case 1: Stemming single word",
			fields: fields{
				logger: logger,
				swg:    swg,
			},
			args: args{
				content: "running",
			},
			want: "run",
		},
		{
			name: "Test case 2: Stemming multiple words",
			fields: fields{
				logger: logger,
				swg:    swg,
			},
			args: args{
				content: "jumps jumped jumping",
			},
			want: "jump jum jum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TermMatcher{
				logger: tt.fields.logger,
				swg:    tt.fields.swg,
			}
			if got := tm.stemContent(tt.args.content); got != tt.want {
				t.Errorf("TermMatcher.stemContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTermMatcher_convertToLowercase(t *testing.T) {
	logger := loggo.NewMockLogger()
	swg := metrics.NewSmithWatermanGotoh()

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test case 1: Lowercase single word",
			fields: fields{
				logger: logger,
				swg:    swg,
			},
			args: args{
				content: "Hello",
			},
			want: "hello",
		},
		{
			name: "Test case 2: Lowercase sentence",
			fields: fields{
				logger: logger,
				swg:    swg,
			},
			args: args{
				content: "Hello World!",
			},
			want: "hello world!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TermMatcher{
				logger: tt.fields.logger,
				swg:    tt.fields.swg,
			}
			if got := tm.convertToLowercase(tt.args.content); got != tt.want {
				t.Errorf("TermMatcher.convertToLowercase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTermMatcher_compareAndAppendTerm(t *testing.T) {
	logger := loggo.NewMockLogger()
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
			tm := &TermMatcher{
				logger: tt.fields.logger,
				swg:    tt.fields.swg,
			}
			if got := tm.compareAndAppendTerm(tt.args.searchTerm, tt.args.content); got != tt.want {
				t.Errorf("TermMatcher.compareAndAppendTerm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTermMatcher_CompareTerms(t *testing.T) {
	logger := loggo.NewMockLogger()
	swg := metrics.NewSmithWatermanGotoh()

	type args struct {
		searchTerm string
		content    string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "Test case 1: High similarity",
			fields: fields{
				logger: logger,
				swg:    swg,
			},
			args: args{
				searchTerm: "hello world",
				content:    "hello world",
			},
			want: 1.0,
		},
		{
			name: "Test case 2: No similarity",
			fields: fields{
				logger: logger,
				swg:    swg,
			},
			args: args{
				searchTerm: "hello",
				content:    "world",
			},
			want: 0.2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TermMatcher{
				logger: tt.fields.logger,
				swg:    tt.fields.swg,
			}
			if got := tm.CompareTerms(tt.args.searchTerm, tt.args.content); got != tt.want {
				t.Errorf("TermMatcher.CompareTerms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTermMatcher_combineContents(t *testing.T) {
	logger := loggo.NewMockLogger()
	swg := metrics.NewSmithWatermanGotoh()

	type args struct {
		content1 string
		content2 string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test case 1: Combine two non-empty contents",
			fields: fields{
				logger: logger,
				swg:    swg,
			},
			args: args{
				content1: "hello",
				content2: "world",
			},
			want: "hello world",
		},
		{
			name: "Test case 2: Combine content with empty string",
			fields: fields{
				logger: logger,
				swg:    swg,
			},
			args: args{
				content1: "hello",
				content2: "",
			},
			want: "hello",
		},
		// TODO: Add more test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := &TermMatcher{
				logger: tt.fields.logger,
				swg:    tt.fields.swg,
			}
			if got := tm.combineContents(tt.args.content1, tt.args.content2); got != tt.want {
				t.Errorf("TermMatcher.combineContents() = %v, want %v", got, tt.want)
			}
		})
	}
}
