package dbmanager

import (
	"context"

	"github.com/jonesrussell/page-prowler/internal/prowlredis"
	"github.com/jonesrussell/page-prowler/models"
)

type MockDBManager struct {
	// You can add more fields if needed
	SavedResults []models.PageData
}

func NewMockDBManager() *MockDBManager {
	return &MockDBManager{}
}
func (m *MockDBManager) SaveResultsToRedis(_ context.Context, results []models.PageData, _ string) error {
	m.SavedResults = append(m.SavedResults, results...)
	return nil
}

func (m *MockDBManager) ClearRedisSet(_ context.Context, _ string) error {
	// Implement this if you use it in your tests
	return nil
}

func (m *MockDBManager) GetLinksFromRedis(_ context.Context, _ string) ([]string, error) {
	// Implement this if you use it in your tests
	return nil, nil
}

func (m *MockDBManager) GetResultsFromRedis(_ context.Context, _ string) ([]models.PageData, error) {
	// Return the saved results
	return m.SavedResults, nil
}

func (m *MockDBManager) RedisOptions() prowlredis.Options {
	// Implement this if you use it in your tests
	return prowlredis.Options{}
}
