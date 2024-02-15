package consumer

import (
	"context"
	"testing"
	"time"

	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveAndUnmarshalLinks(t *testing.T) {
	ctx := context.Background()
	manager := &crawler.CrawlManager{
		Client: prowlredis.NewMockClient(), // Use the mock client
	}
	siteid := "testsite"

	// Add some mock data to the client
	link1 := `{"url": "https://example.com/1", "matching_terms": ["term1", "term2"]}`
	link2 := `{"url": "https://example.com/2", "matching_terms": ["term3", "term4"]}`
	err := manager.Client.SAdd(ctx, siteid, link1, link2)
	if err != nil {
		return
	}

	links, err := RetrieveAndUnmarshalLinks(ctx, manager, siteid)
	assert.NoError(t, err)
	assert.NotNil(t, links)
}

func TestRetrieveAndUnmarshalLinksEmptySet(t *testing.T) {
	ctx := context.Background()
	manager := &crawler.CrawlManager{
		Client: prowlredis.NewMockClient(), // Use the mock client
	}
	siteid := "emptyset"

	// No mock data added to the client

	links, err := RetrieveAndUnmarshalLinks(ctx, manager, siteid)
	assert.NoError(t, err)
	assert.Nil(t, links)
}

func TestCreateOutput(t *testing.T) {
	siteid := "testsite"
	var links []Link // Initialize with mock data

	output := CreateOutput(siteid, links)
	assert.Equal(t, siteid, output.Siteid)
	assert.WithinDuration(t, time.Now(), output.Timestamp, time.Second)
	assert.Equal(t, "success", output.Status)
	assert.Equal(t, "Links retrieved successfully", output.Message)
	assert.Equal(t, links, output.Links)
}

func TestMarshalOutput(t *testing.T) {
	output := Output{} // Initialize with mock data

	data, err := MarshalOutput(output)
	assert.NoError(t, err)
	assert.NotNil(t, data)
}
