package tasks

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/mock"
)

// Create a mock of your interface.
type MockAsynqClient struct {
	mock.Mock
}

func (m *MockAsynqClient) Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	args := m.Called(task, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*asynq.TaskInfo), args.Error(1)
}

func TestEnqueueCrawlTask(t *testing.T) {
	client := new(MockAsynqClient)
	payload := &CrawlTaskPayload{
		URL:         "https://example.com",
		SearchTerms: "example",
		CrawlSiteID: "site123",
		MaxDepth:    3,
		Debug:       false,
	}

	// Define the expected behavior for the mock.
	client.On("Enqueue", mock.Anything, mock.Anything).Return(&asynq.TaskInfo{ID: "123"}, nil)

	id, err := EnqueueCrawlTask(client, payload)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if id != "123" {
		t.Errorf("Expected ID to be '123', got %v", id)
	}

	// Assert that the mock method was called.
	client.AssertExpectations(t)
}

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

func TestNewCrawlTaskInvalidPayload(t *testing.T) {
	// Add test cases for invalid payloads
}

func TestEnqueueCrawlTaskError(t *testing.T) {
	client := new(MockAsynqClient)
	payload := &CrawlTaskPayload{
		URL:         "https://example.com",
		SearchTerms: "example",
		CrawlSiteID: "site123",
		MaxDepth:    3,
		Debug:       false,
	}

	// Define the expected behavior for the mock.
	client.On("Enqueue", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("enqueue error"))

	_, err := EnqueueCrawlTask(client, payload)
	if err == nil {
		t.Errorf("Expected error, got none")
	}

	// Assert that the mock method was called.
	client.AssertExpectations(t)
}

func TestNewCrawlTaskCorrectTaskCreation(t *testing.T) {
	payload := &CrawlTaskPayload{
		URL:         "https://example.com",
		SearchTerms: "example",
		CrawlSiteID: "site123",
		MaxDepth:    3,
		Debug:       false,
	}

	task, _ := NewCrawlTask(payload)
	if task.Type() != CrawlTaskType {
		t.Errorf("Expected task type to be '%s', got '%s'", CrawlTaskType, task.Type())
	}
}
