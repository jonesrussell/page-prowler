// main_test.go

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommandLineArguments(t *testing.T) {
	args := []string{"./crawler", "https://www.example.com", "test-group"}
	crawlURL, group, err := parseCommandLineArguments(args)

	// Assertions
	assert.Equal(t, "https://www.example.com", crawlURL, "Expected crawlURL to match")
	assert.Equal(t, "test-group", group, "Expected group to match")
	assert.NoError(t, err, "Expected no error")
}

func TestParseCommandLineArgumentsInvalid(t *testing.T) {
	args := []string{"./crawler"}
	crawlURL, group, err := parseCommandLineArguments(args)

	// Assertions
	assert.Equal(t, "", crawlURL, "Expected crawlURL to be empty")
	assert.Equal(t, "", group, "Expected group to be empty")
	assert.Error(t, err, "Expected an error")
}

func TestCreateLogger(t *testing.T) {
	logger := createLogger()

	// Check if the logger is not nil
	assert.NotNil(t, logger, "Expected logger not to be nil")

	// Log a message
	logger.Info("Test log message")

	// No need to check for a return value from logger.Info
}

func TestCreateRedisClient(t *testing.T) {
	// Set up environment variables for testing
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("REDIS_AUTH", "")

	redisClient := createRedisClient()

	// Assertions
	assert.NotNil(t, redisClient, "Expected redisClient not to be nil")
	// You can add more assertions based on your requirements
}

func TestConfigureCollector(t *testing.T) {
	collector := configureCollector()

	// Assertions
	assert.NotNil(t, collector, "Expected collector not to be nil")
	assert.True(t, collector.Async, "Expected collector to be asynchronous")
	assert.Equal(t, 3, collector.MaxDepth, "Expected MaxDepth to be 3")
	// You can add more assertions based on your requirements
}

func TestMain(m *testing.M) {
	// Run tests and exit with the result
	os.Exit(m.Run())
}
