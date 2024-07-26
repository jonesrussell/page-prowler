package crawler_test

import (
	"sync"
	"testing"

	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/jonesrussell/page-prowler/mocks"
)

func TestUtils_GetHostFromURL(t *testing.T) {
	type args struct {
		inputURL  string
		appLogger loggo.Logger
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid URL",
			args: args{
				inputURL:  "https://www.example.com",
				appLogger: mocks.NewMockLogger(),
			},
			want:    "www.example.com",
			wantErr: false,
		},
		{
			name: "Invalid URL",
			args: args{
				inputURL:  "not a valid url",
				appLogger: mocks.NewMockLogger(),
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := crawler.GetHostFromURL(tt.args.inputURL, tt.args.appLogger)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetHostFromURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrawlManager_ProcessMatchingLinkAndUpdateStats(t *testing.T) {
	type fields struct {
		LoggerField  loggo.Logger
		Client       prowlredis.ClientInterface
		Collector    *crawler.CollectorWrapper
		CrawlingMu   *sync.Mutex
		StatsManager *crawler.StatsManager
	}

	type args struct {
		options       *crawler.CrawlOptions
		href          string
		pageData      crawler.PageData
		matchingTerms []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test Case 1",
			fields: fields{
				LoggerField:  mocks.NewMockLogger(),
				Client:       mocks.NewMockClient(),
				Collector:    &crawler.CollectorWrapper{},
				CrawlingMu:   &sync.Mutex{},
				StatsManager: crawler.NewStatsManager(),
			},
			args: args{
				href:          "https://example.com",
				pageData:      crawler.PageData{},
				matchingTerms: []string{"example"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &crawler.CrawlManager{
				LoggerField: tt.fields.LoggerField,
				Client:      tt.fields.Client,

				CollectorInstance: tt.fields.Collector,
				CrawlingMu:        tt.fields.CrawlingMu,
				StatsManager:      tt.fields.StatsManager,
			}

			cs.ProcessMatchingLink(tt.args.href, tt.args.pageData, tt.args.matchingTerms)
			cs.UpdateStats(tt.args.options, tt.args.matchingTerms)
		})
	}
}
