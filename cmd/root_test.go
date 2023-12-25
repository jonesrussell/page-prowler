package cmd

import (
	"log"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	assert.Equal(t, rootCmd.Use, "page-prowler", "rootCmd.Use should be 'page-prowler'")
}

func TestExecute(t *testing.T) {
	// Call Execute with no arguments
	os.Args = []string{"cmd"}
	err := Execute()
	assert.NoError(t, err, "Execute() without arguments should not return an error")
}

func TestInitConfig(t *testing.T) {
	initConfig()

	// Check if the environment variables are correctly set
	assert.Equal(t, viper.GetString("REDIS_HOST"), os.Getenv("REDIS_HOST"), "REDIS_HOST is not correctly set")
	assert.Equal(t, viper.GetString("REDIS_PORT"), os.Getenv("REDIS_PORT"), "REDIS_PORT is not correctly set")
	assert.Equal(t, viper.GetString("REDIS_AUTH"), os.Getenv("REDIS_AUTH"), "REDIS_AUTH is not correctly set")
}

func TestFlagValues(t *testing.T) {
	// Set a flag value
	rootCmd.SetArgs([]string{"--debug=true"})
	err := rootCmd.Execute()
	assert.NoError(t, err)

	// Check if the flag value is correctly set
	assert.True(t, viper.GetBool("debug"))
}

func TestEnvironmentVariables(t *testing.T) {
	// Set an environment variable
	os.Setenv("REDIS_HOST", "localhost")
	initConfig()

	// Check if the environment variable is correctly set
	assert.Equal(t, "localhost", viper.GetString("REDIS_HOST"))
}

func TestExecuteError(t *testing.T) {
	// Provide an invalid command
	rootCmd.SetArgs([]string{"invalidCommand"})
	err := rootCmd.Execute()

	// Check if Execute returns an error
	assert.Error(t, err)
}

func TestPersistentFlags(t *testing.T) {
	// Set the flags
	if err := rootCmd.PersistentFlags().Set("crawlsiteid", "123"); err != nil {
		log.Fatalf("Error setting crawlsiteid flag: %v", err)
	}

	if err := rootCmd.PersistentFlags().Set("debug", "true"); err != nil {
		log.Fatalf("Error setting debug flag: %v", err)
	}

	// Check if the flags are correctly set
	assert.Equal(t, "123", viper.GetString("crawlsiteid"))
	assert.True(t, viper.GetBool("debug"))
}

func TestInitConfigError(t *testing.T) {
	// Temporarily replace the config file with a non-existent file
	origConfigFile := viper.ConfigFileUsed()
	defer func() { viper.SetConfigFile(origConfigFile) }()
	viper.SetConfigFile("non_existent_file")

	// Call initConfig and check if it doesn't panic
	assert.NotPanics(t, func() { initConfig() }, "initConfig should not panic if the config file does not exist")
}
