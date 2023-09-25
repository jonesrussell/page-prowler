// main_test.go

package main

import (
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
