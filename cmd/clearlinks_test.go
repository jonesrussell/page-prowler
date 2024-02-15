package cmd

import (
	"context"
	"testing"

	"github.com/jonesrussell/page-prowler/cmd/mocks"
	"github.com/jonesrussell/page-prowler/internal/common"
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestClearlinksCmd_RunE(t *testing.T) {
	ctx := context.Background()
	mockLogger := mocks.NewMockLogger()
	cm := &crawler.CrawlManager{
		LoggerField: mockLogger,
		Client:      prowlredis.NewMockClient(),
	}
	ctx = context.WithValue(ctx, common.CrawlManagerKey, cm)

	cmd := &cobra.Command{
		Use:   "clearlinks",
		Short: "Clear the Redis set for a given siteid",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return ClearlinksCmd.RunE(cmd, []string{})
		},
	}
	cmd.Flags().StringVarP(&Siteid, "siteid", "s", "", "CrawlSite ID")
	cmd.SetArgs([]string{"--siteid=test"})

	err := cmd.ExecuteContext(ctx)
	assert.Nil(t, err)
	if mockClientInstance, ok := cm.Client.(*prowlredis.MockClient); ok {
		assert.True(t, mockClientInstance.WasDelCalled)
	} else {
		t.Fatal("Failed to cast mockClient to *prowlredis.MockClient")
	}
}
