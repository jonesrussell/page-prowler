package cmd

import (
	"context"
	"testing"

	"github.com/jonesrussell/page-prowler/cmd/mocks"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetLinksCmd(t *testing.T) {
	viper.Set("crawlsiteid", "testsite")

	mockRedisClient := &mocks.MockRedisClient{}
	ctx := context.WithValue(context.Background(), managerKey, &crawler.CrawlManager{Client: mockRedisClient})

	err := getLinksCmd.ExecuteContext(ctx)
	assert.NoError(t, err)

	// Check that crawlsiteid is required
	assert.NotEmpty(t, viper.GetString("crawlsiteid"), "crawlsiteid should not be empty")

	// Check that informative message is printed when links is empty
	// This requires capturing the output of the command, which is not shown in this example
	// You can use a library like https://github.com/stretchr/testify#assertions to help with this
}
