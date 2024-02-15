package cmd_test

import (
	"context"
	"os"
	"testing"

	"github.com/jonesrussell/page-prowler/cmd"
	"github.com/jonesrussell/page-prowler/cmd/mocks"
	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestGetLinksCmd(t *testing.T) {
	// Initialize the command
	getLinksCmd := &cobra.Command{
		Use:   "getlinks",
		Short: "Get the list of links for a given siteid",
		RunE:  cmd.GetLinksCmd.RunE,
	}

	// Set the 'siteid' flag
	getLinksCmd.Flags().StringP("siteid", "s", "", "CrawlSite ID")
	err := viper.BindPFlag("siteid", getLinksCmd.Flags().Lookup("siteid"))
	if err != nil {
		return
	}

	// Enable viper to read from environment variables
	viper.AutomaticEnv()

	// Test cases
	tests := []struct {
		name    string
		env     string
		flag    string
		wantErr bool
	}{
		{
			name:    "Siteid empty, no env, no flag",
			env:     "",
			flag:    "",
			wantErr: true, // Assuming your command returns an error when siteid is not provided
		},
		{
			name:    "Siteid by flag, env not set",
			env:     "",
			flag:    "flag_value",
			wantErr: false,
		},
		{
			name:    "Siteid by flag, env set",
			env:     "env_value",
			flag:    "flag_value",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the 'siteid' environment variable
			err := os.Setenv("SITEID", tt.env)
			if err != nil {
				return
			}

			// Set the 'siteid' flag value
			viper.Set("siteid", tt.flag)

			// Initialize a CrawlManager
			manager := &crawler.CrawlManager{
				Client:         prowlredis.NewMockClient().(*prowlredis.MockClient),
				LoggerField:    mocks.NewMockLogger(),
				MongoDBWrapper: mocks.NewMockMongoDBWrapper(),
			}

			// Create a context with the CrawlManager
			ctx := context.WithValue(context.Background(), common.CrawlManagerKey, manager)

			// Set the context in the command
			getLinksCmd.SetContext(ctx)

			// Execute the command with the 'siteid' flag and environment variable values
			getLinksCmd.SetArgs([]string{"--siteid=" + tt.flag})
			err = getLinksCmd.Execute()

			// Check for error
			if (err != nil) != tt.wantErr {
				t.Errorf("getLinksCmd.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Unset the 'siteid' environment variable
			err = os.Unsetenv("SITEID")
			if err != nil {
				return
			}
		})
	}
}
