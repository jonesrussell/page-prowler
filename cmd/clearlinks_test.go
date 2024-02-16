package cmd_test

import (
	"context"
	"testing"

	"github.com/jonesrussell/page-prowler/cmd"
	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/jonesrussell/page-prowler/mocks"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestClearlinksCmd_RunE(t *testing.T) {
	ctx := context.Background()
	mockLogger := mocks.NewMockLogger()
	cm := &crawler.CrawlManager{
		LoggerField: mockLogger,
		Client:      prowlredis.NewMockClient(),
	}
	ctx = context.WithValue(ctx, common.CrawlManagerKey, cm)

	testCmd := &cobra.Command{
		Use:   "clearlinks",
		Short: "Clear the Redis set for a given siteid",
		RunE: func(testCmd *cobra.Command, _ []string) error {
			return cmd.ClearlinksCmd.RunE(testCmd, []string{})
		},
	}
	testCmd.Flags().StringVarP(&cmd.Siteid, "siteid", "s", "", "CrawlSite ID")
	testCmd.SetArgs([]string{"--siteid=test"})

	err := testCmd.ExecuteContext(ctx)
	assert.Nil(t, err)
	if mockClientInstance, ok := cm.Client.(*prowlredis.MockClient); ok {
		assert.True(t, mockClientInstance.WasDelCalled)
	} else {
		t.Fatal("Failed to cast mockClient to *prowlredis.MockClient")
	}
}

func Test_ClearlinksFunc(t *testing.T) {
	type args struct {
		testCmd *cobra.Command
		in1     []string
	}
	tests := []struct {
		name    string
		args    args
		setup   func()
		wantErr bool
	}{
		{
			name: "Missing SiteID",
			args: args{
				testCmd: &cobra.Command{},
				in1:     []string{},
			},
			setup: func() {
				cmd.Siteid = ""
			},
			wantErr: true,
		},
		{
			name: "Redis Error",
			args: args{
				testCmd: &cobra.Command{},
				in1:     []string{},
			},
			setup: func() {
				cmd.Siteid = "error-site-id"
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup() // Set up the test environment
			if err := cmd.ClearlinksMain(tt.args.testCmd, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("ClearlinksFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
