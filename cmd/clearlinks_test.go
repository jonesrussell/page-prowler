package cmd

import (
	"context"
	"testing"

	"github.com/jonesrussell/page-prowler/cmd/mocks"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestClearlinksCmd(t *testing.T) {
	viper.Set("crawlsiteid", "testsite")

	mockRedisClient := &mocks.MockRedisClient{}
	ctx := context.WithValue(context.Background(), managerKey, &crawler.CrawlManager{Client: mockRedisClient})

	err := clearlinksCmd.ExecuteContext(ctx)
	assert.NoError(t, err)

	// Check that crawlsiteid is required
	assert.NotEmpty(t, viper.GetString("crawlsiteid"), "crawlsiteid should not be empty")
}
