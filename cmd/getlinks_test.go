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
		Short: "Get the list of links for a given crawlsiteid",
		RunE:  cmd.GetLinksCmd.RunE,
	}

	// Set the 'crawlsiteid' flag
	getLinksCmd.Flags().StringP("crawlsiteid", "s", "", "CrawlSite ID")
	viper.BindPFlag("crawlsiteid", getLinksCmd.Flags().Lookup("crawlsiteid"))

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
			name:    "Crawlsiteid empty, no env, no flag",
			env:     "",
			flag:    "",
			wantErr: true, // Assuming your command returns an error when crawlsiteid is not provided
		},
		{
			name:    "Crawlsiteid by flag, env not set",
			env:     "",
			flag:    "flag_value",
			wantErr: false,
		},
		{
			name:    "Crawlsiteid by flag, env set",
			env:     "env_value",
			flag:    "flag_value",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the 'crawlsiteid' environment variable
			os.Setenv("CRAWLSITEID", tt.env)

			// Set the 'crawlsiteid' flag value
			viper.Set("crawlsiteid", tt.flag)

			// Initialize a CrawlManager
			manager := &crawler.CrawlManager{
				Client:         prowlredis.NewMockClient().(*prowlredis.MockClient),
				Logger:         mocks.NewMockLogger(),
				MongoDBWrapper: mocks.NewMockMongoDBWrapper(),
			}

			// Create a context with the CrawlManager
			ctx := context.WithValue(context.Background(), common.CrawlManagerKey, manager)

			// Set the context in the command
			getLinksCmd.SetContext(ctx)

			// Execute the command with the 'crawlsiteid' flag and environment variable values
			getLinksCmd.SetArgs([]string{"--crawlsiteid=" + tt.flag})
			err := getLinksCmd.Execute()

			// Unset the 'crawlsiteid' environment variable
			os.Unsetenv("CRAWLSITEID")

			// Check for error
			if (err != nil) != tt.wantErr {
				t.Errorf("getLinksCmd.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
