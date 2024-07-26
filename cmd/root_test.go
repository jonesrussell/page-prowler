package cmd_test

import (
	"log"
	"os"
	"testing"

	"github.com/jonesrussell/page-prowler/cmd"
	"github.com/jonesrussell/page-prowler/mocks"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	manager := // initialize your manager here...
	rootCmd := cmd.NewRootCmd(manager)
	assert.Equal(t, rootCmd.Use, "page-prowler", "RootCmd.Use should be 'page-prowler'")
}

func TestExecute(t *testing.T) {
	// Call Execute with no arguments
	os.Args = []string{"cmd"}
	err := cmd.RootCmd.Execute()
	assert.NoError(t, err, "Execute() without arguments should not return an error")
}

func TestFlagValues(t *testing.T) {
	// Set a flag value
	cmd.RootCmd.SetArgs([]string{"--debug=true"})
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
	assert.True(t, viper.GetBool("debug"))
}

func TestExecuteError(t *testing.T) {
	// Provide an invalid command
	cmd.RootCmd.SetArgs([]string{"invalidCommand"})
	err := cmd.RootCmd.Execute()

	// Check if Execute returns an error
	assert.Error(t, err)
}

func TestPersistentFlags(t *testing.T) {
	// Set the flags
	if err := cmd.RootCmd.PersistentFlags().Set("debug", "true"); err != nil {
		log.Fatalf("Error setting debug flag: %v", err)
	}

	// Bind the flags to viper
	cmd.RootCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		err := viper.BindPFlag(flag.Name, flag)
		if err != nil {
			return
		}
	})

	// Check if the flags are correctly set
	assert.True(t, viper.GetBool("debug"))
}

func TestInitializeManager_WithNilRedisClient(t *testing.T) {
	// Initialize the manager with a nil Redis client and a new mock logger
	_, err := cmd.InitializeManager(
		nil,
		mocks.NewMockLogger(),
	)

	// Check if an error was returned
	assert.Error(t, err, "Expected an error when initializing with a nil Redis client")
}

func TestInitializeManager(t *testing.T) {
	// Set the environment variables
	err := os.Setenv("REDIS_HOST", "172.17.0.1")
	if err != nil {
		log.Fatalf("Failed to set environment variable: %v", err)
	}
	err = os.Setenv("REDIS_AUTH", "password")
	if err != nil {
		log.Fatalf("Failed to set environment variable: %v", err)
	}
	err = os.Setenv("REDIS_PORT", "6379")
	if err != nil {
		log.Fatalf("Failed to set environment variable: %v", err)
	}

	// Initialize the manager with a mock Redis client and a new mock logger
	manager, err := cmd.InitializeManager(
		mocks.NewMockClient(),
		mocks.NewMockLogger(),
	)
	if err != nil {
		t.Fatalf("Failed to initialize manager: %v", err)
	}

	// Add assertions
	assert.NotNil(t, manager.Client, "Client should not be nil")
	assert.NotNil(t, manager.Logger, "Logger should not be nil")
}
