package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockCrawlManager struct {
	mock.Mock
}

func (m *MockCrawlManager) StartCrawling(ctx context.Context, url string, searchterms string, crawlsiteid string, maxdepth int, debug bool) error {
	args := m.Called(ctx, url, searchterms, crawlsiteid, maxdepth, debug)
	return args.Error(0)
}
