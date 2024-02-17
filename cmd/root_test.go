package cmd_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jonesrussell/page-prowler/cmd"
	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/mongodbwrapper"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/jonesrussell/page-prowler/mocks"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_RootCmd(t *testing.T) {
	assert.Equal(t, cmd.RootCmd.Use, "page-prowler", "RootCmd.Use should be 'page-prowler'")
}

func Test_Execute(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "Execute with no arguments",
			args:    []string{"cmd"},
			wantErr: false,
		},
		{
			name:    "Execute with invalid command",
			args:    []string{"cmd", "invalidCommand"},
			wantErr: true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the command arguments
			os.Args = tt.args
			err := cmd.RootCmd.Execute()

			// Check if Execute returns an error
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_FlagValues(t *testing.T) {
	tests := []struct {
		name      string
		flagName  string
		flagValue string
		want      interface{}
	}{
		{
			name:      "Debug flag set to true",
			flagName:  "debug",
			flagValue: "true",
			want:      true,
		},
		{
			name:      "Siteid flag set to a value",
			flagName:  "siteid",
			flagValue: "exampleSiteId",
			want:      "exampleSiteId",
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the flag value
			cmd.RootCmd.SetArgs([]string{fmt.Sprintf("--%s=%s", tt.flagName, tt.flagValue)})
			err := cmd.RootCmd.Execute()
			assert.NoError(t, err)

			// Bind the flags to viper
			cmd.RootCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
				err := viper.BindPFlag(flag.Name, flag)
				if err != nil {
					return
				}
			})

			// Check if the flag value is correctly set
			assert.Equal(t, tt.want, viper.Get(tt.flagName))
		})
	}
}

func Test_InitializeManager(t *testing.T) {
	tests := []struct {
		name           string
		redisClient    prowlredis.ClientInterface
		appLogger      logger.Logger
		mongoDBWrapper mongodbwrapper.MongoDBInterface
		wantErr        bool
	}{
		{
			name:           "With nil MongoDB wrapper",
			redisClient:    mocks.NewMockClient(),
			appLogger:      mocks.NewMockLogger(),
			mongoDBWrapper: nil,
			wantErr:        true,
		},
		{
			name:           "With nil Redis client",
			redisClient:    nil,
			appLogger:      mocks.NewMockLogger(),
			mongoDBWrapper: mocks.NewMockMongoDBWrapper(),
			wantErr:        true,
		},
		{
			name:           "With nil logger",
			redisClient:    mocks.NewMockClient(),
			appLogger:      nil,
			mongoDBWrapper: mocks.NewMockMongoDBWrapper(),
			wantErr:        true,
		},
		{
			name:           "With valid dependencies",
			redisClient:    mocks.NewMockClient(),
			appLogger:      mocks.NewMockLogger(),
			mongoDBWrapper: mocks.NewMockMongoDBWrapper(),
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cmd.InitializeManager(tt.redisClient, tt.appLogger, tt.mongoDBWrapper)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitializeManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_PersistentPreRunE(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		setup   func() // Function to set up the environment for the test
		wantErr bool
	}{
		{
			name: "Valid command with no arguments",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{},
			},
			setup: func() {
				// Use mocks for Redis client and MongoDB wrapper
				redisClient := mocks.NewMockClient()
				mongoDBWrapper := mocks.NewMockMongoDBWrapper()
				appLogger := mocks.NewMockLogger()

				// Initialize the manager with mocks
				manager, err := cmd.InitializeManager(redisClient, appLogger, mongoDBWrapper)
				if err != nil {
					t.Fatalf("Failed to initialize manager: %v", err)
				}

				// Set the manager to the context
				ctx := context.WithValue(context.Background(), common.CrawlManagerKey, manager)

				// Set the context of the command
				cmd.RootCmd.SetContext(ctx)
			},
			wantErr: false,
		},
		{
			name: "Missing REDIS_HOST environment variable",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{},
			},
			setup: func() {
				// Unset REDIS_HOST to simulate missing environment variable
				os.Unsetenv("REDIS_HOST")
			},
			wantErr: true,
		},
		{
			name: "Missing MONGODB_URI environment variable",
			args: args{
				cmd:  &cobra.Command{},
				args: []string{},
			},
			setup: func() {
				// Unset MONGODB_URI to simulate missing environment variable
				os.Unsetenv("MONGODB_URI")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup() // Set up the environment for the test
			if err := cmd.RootCmd.PersistentPreRunE(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("PersistentPreRunE() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
