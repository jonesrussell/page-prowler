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
	crawlsiteid := "testsite"

	// Add some mock data to the client
	link1 := `{"url": "http://example.com/1", "matching_terms": ["term1", "term2"]}`
	link2 := `{"url": "http://example.com/2", "matching_terms": ["term3", "term4"]}`
	manager.Client.SAdd(ctx, crawlsiteid, link1, link2)

	links, err := RetrieveAndUnmarshalLinks(ctx, manager, crawlsiteid)
	assert.NoError(t, err)
	assert.NotNil(t, links)
}

func TestRetrieveAndUnmarshalLinksEmptySet(t *testing.T) {
	ctx := context.Background()
	manager := &crawler.CrawlManager{
		Client: prowlredis.NewMockClient(), // Use the mock client
	}
	crawlsiteid := "emptyset"

	// No mock data added to the client

	links, err := RetrieveAndUnmarshalLinks(ctx, manager, crawlsiteid)
	assert.NoError(t, err)
	assert.Nil(t, links)
}

func TestCreateOutput(t *testing.T) {
	crawlsiteid := "testsite"
	links := []Link{} // Initialize with mock data

	output := CreateOutput(crawlsiteid, links)
	assert.Equal(t, crawlsiteid, output.Crawlsiteid)
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
