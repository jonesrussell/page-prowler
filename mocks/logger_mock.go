package mocks

import (
	"github.com/gocolly/colly/debug"
	"github.com/jonesrussell/loggo"
)

type MockLogger struct {
	events []*debug.Event
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	// Implement the method here.
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	// Implement the method here.
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	// Implement the method here.
}

func (m *MockLogger) Error(msg string, err error, args ...interface{}) {
	// Implement the method here.
}

func (m *MockLogger) Fatal(msg string, err error, args ...interface{}) {
	// Implement the method here.
}

func (m *MockLogger) WithOperation(operationID string) loggo.LoggerInterface {
	// Implement the method here.
	return m
}

// Event implements logger.Logger.
func (m *MockLogger) Event(e *debug.Event) {
	// Store the event for later assertions
	m.events = append(m.events, e)
}

func NewMockLogger() *MockLogger {
	return &MockLogger{}
}
