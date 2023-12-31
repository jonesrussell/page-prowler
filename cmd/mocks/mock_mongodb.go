package mocks

import "context"

type MockMongoDBWrapper struct{}

func NewMockMongoDBWrapper() *MockMongoDBWrapper {
	// Initialize your MockMongoDBWrapper here, if necessary
	return &MockMongoDBWrapper{}
}

func (m *MockMongoDBWrapper) Connect(_ context.Context) error {
	// Implement your mock logic here
	return nil
}
