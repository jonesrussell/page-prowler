package cmd_test

import (
	"context"
	"testing"

	"github.com/jonesrussell/page-prowler/cmd"
	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/mocks"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Test_GetlinksCmd(t *testing.T) {
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
				viper.Set("Siteid", "")
			},
			wantErr: true,
		},
		{
			name: "Valid SiteID",
			args: args{
				testCmd: &cobra.Command{},
				in1:     []string{},
			},
			setup: func() {
				viper.Set("Siteid", "foobar")
			},
			wantErr: false,
		},
		// @TODO test returned links
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup() // Set up the test environment

			// Create a mock CrawlManager with a redisClient
			mockManager := &crawler.CrawlManager{
				Client:      mocks.NewMockClient(),
				LoggerField: mocks.NewMockLogger(),
			}

			// Add the mock CrawlManager to the context
			ctx := context.WithValue(context.Background(), common.CrawlManagerKey, mockManager)

			// Create a new command with the context
			testCmd := &cobra.Command{}
			testCmd.SetContext(ctx)

			// Call ClearlinksMain with the command that has the context
			if err := cmd.GetLinksCmd.RunE(testCmd, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("GetlinksCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
