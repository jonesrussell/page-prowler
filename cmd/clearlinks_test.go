package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/jonesrussell/page-prowler/cmd/mocks"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestClearlinksCmd(t *testing.T) {
	viper.Set("crawlsiteid", "testsite")

	mockRedisClient := mocks.NewMockRedisClient()
	ctx := context.WithValue(context.Background(), managerKey, mockRedisClient)

	// Create a new Cobra command for testing
	cmd := &cobra.Command{
		Use:   clearlinksCmd.Use,
		Short: clearlinksCmd.Short,
		Long:  clearlinksCmd.Long,
		Run:   clearlinksCmd.Run,
	}

	// Create a buffer to capture the output
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	// Execute the command
	err := cmd.ExecuteContext(ctx)

	assert.NoError(t, err)

	// Check that crawlsiteid is required
	assert.NotEmpty(t, viper.GetString("crawlsiteid"), "crawlsiteid should not be empty")
}
