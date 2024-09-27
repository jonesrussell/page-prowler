package dbmanager

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/stretchr/testify/assert"
)

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
				mockLogger.EXPECT().Debug(gomock.Any()).AnyTimes()
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
