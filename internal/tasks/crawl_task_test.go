package tasks

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewCrawlTask(t *testing.T) {
	type args struct {
		payload *CrawlTaskPayload
	}
	tests := []struct {
		name    string
		args    args
		want    *CrawlTaskPayload
		wantErr bool
	}{
		{
			name: "valid payload",
			args: args{
				payload: &CrawlTaskPayload{
					URL:         "http://example.com",
					SearchTerms: "example",
					CrawlSiteID: "site123",
					MaxDepth:    3,
					Debug:       false,
				},
			},
			want: &CrawlTaskPayload{
				URL:         "http://example.com",
				SearchTerms: "example",
				CrawlSiteID: "site123",
				MaxDepth:    3,
				Debug:       false,
			},
			wantErr: false,
		},
		{
			name: "invalid payload",
			args: args{
				payload: &CrawlTaskPayload{
					URL:         "",
					SearchTerms: "example",
					CrawlSiteID: "site123",
					MaxDepth:    3,
					Debug:       false,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTask, err := NewCrawlTask(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCrawlTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTask != nil {
				var gotPayload CrawlTaskPayload
				if err := json.Unmarshal(gotTask.Payload(), &gotPayload); err != nil {
					t.Errorf("Failed to unmarshal payload: %v", err)
					return
				}
				if !reflect.DeepEqual(&gotPayload, tt.want) {
					t.Errorf("NewCrawlTask() = %v, want %v", &gotPayload, tt.want)
				}
			}
		})
	}
}
