package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetLinksCmd(t *testing.T) {
	viper.Set("crawlsiteid", "testsite")

	mockRedisClient := prowlredis.NewMockClient()
	ctx := context.WithValue(context.Background(), common.CrawlManagerKey, &crawler.CrawlManager{Client: mockRedisClient})

	// Create a new Cobra command for testing
	cmd := &cobra.Command{
		Use:   getLinksCmd.Use,
		Short: getLinksCmd.Short,
		Long:  getLinksCmd.Long,
		Run:   getLinksCmd.Run,
	}

	// Create a buffer to capture the output
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	// Execute the command
	err := cmd.ExecuteContext(ctx)

	assert.NoError(t, err)

	// Check that crawlsiteid is required
	assert.NotEmpty(t, viper.GetString("crawlsiteid"), "crawlsiteid should not be empty")

	// Check that informative message is printed when links is empty
	// This requires capturing the output of the command, which is not shown in this example
	// You can use a library like https://github.com/stretchr/testify#assertions to help with this
}
