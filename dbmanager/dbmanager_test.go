package dbmanager

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/jonesrussell/page-prowler/models"
	"github.com/stretchr/testify/assert"
)

func TestSaveResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := prowlredis.NewMockClientInterface(ctrl)
	mockLogger := loggo.NewMockLoggerInterface(ctrl)
	redisManager := NewRedisManager(mockClient, mockLogger)
	ctx := context.TODO()

	tests := []struct {
		name    string
		results []models.PageData
		key     string
		setup   func()
		wantErr bool
	}{
		{
			name: "successful save",
			results: []models.PageData{
				{Title: "Title1", URL: "http://example.com/1"},
			},
			key: "key1",
			setup: func() {
				mockLogger.EXPECT().Debug("Redis", "key", "key1").Times(1)
				mockLogger.EXPECT().Debug("Redis", "results", gomock.Any()).Times(1)
				mockClient.EXPECT().SAdd(ctx, "key1", gomock.Any()).Return(nil).Times(1)
			},
			wantErr: false,
		},
		{
			name: "marshal error",
			results: []models.PageData{
				{Title: "Title1", URL: string([]byte{0xff})}, // invalid UTF-8
			},
			key: "key1",
			setup: func() {
				mockLogger.EXPECT().Debug("Redis", "key", "key1").Times(1)
				mockLogger.EXPECT().Debug("Redis", "results", gomock.Any()).Times(1)
			},
			wantErr: true,
		},
		{
			name: "client SAdd error",
			results: []models.PageData{
				{Title: "Title1", URL: "http://example.com/1"},
			},
			key: "key1",
			setup: func() {
				mockLogger.EXPECT().Debug("Redis", "key", "key1").Times(1)
				mockLogger.EXPECT().Debug("Redis", "results", gomock.Any()).Times(1)
				mockClient.EXPECT().SAdd(ctx, "key1", gomock.Any()).Return(errors.New("SAdd error")).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			err := redisManager.SaveResults(ctx, tt.results, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveResults() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClearRedisSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := prowlredis.NewMockClientInterface(ctrl)
	mockLogger := loggo.NewMockLoggerInterface(ctrl)
	redisManager := NewRedisManager(mockClient, mockLogger)
	ctx := context.TODO()

	tests := []struct {
		name    string
		key     string
		setup   func()
		wantErr bool
	}{
		{
			name: "successful delete",
			key:  "key1",
			setup: func() {
				mockClient.EXPECT().Del(ctx, "key1").Return(nil).Times(1)
			},
			wantErr: false,
		},
		{
			name: "client Del error",
			key:  "key1",
			setup: func() {
				mockClient.EXPECT().Del(ctx, "key1").Return(errors.New("Del error")).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			err := redisManager.ClearRedisSet(ctx, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClearRedisSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetLinksFromRedis(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := prowlredis.NewMockClientInterface(ctrl)
	mockLogger := loggo.NewMockLoggerInterface(ctrl)
	redisManager := NewRedisManager(mockClient, mockLogger)
	ctx := context.TODO()

	tests := []struct {
		name    string
		key     string
		setup   func()
		want    []string
		wantErr bool
	}{
		{
			name: "successful get",
			key:  "key1",
			setup: func() {
				mockClient.EXPECT().SMembers(ctx, "key1").Return([]string{"link1", "link2"}, nil).Times(1)
			},
			want:    []string{"link1", "link2"},
			wantErr: false,
		},
		{
			name: "client SMembers error",
			key:  "key1",
			setup: func() {
				mockClient.EXPECT().SMembers(ctx, "key1").Return(nil, errors.New("SMembers error")).Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			got, err := redisManager.GetLinksFromRedis(ctx, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLinksFromRedis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}