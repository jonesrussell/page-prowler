package mocks

import (
	"time"

	"github.com/gocolly/colly/debug"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type MockLogger struct {
	observer *observer.ObservedLogs
	Logger   *zap.Logger
	start    time.Time
	events   []*debug.Event
}

func (m *MockLogger) Info(msg string) {
	m.Logger.Info(msg)
}

func (m *MockLogger) Debug(msg string) {
	m.Logger.Debug(msg)
}

func (m *MockLogger) Error(msg string) {
	m.Logger.Error(msg)
}

func (m *MockLogger) Warn(msg string) {
	m.Logger.Warn(msg)
}

func (m *MockLogger) Fatal(msg string) {
	m.Logger.Fatal(msg)
}

// Event implements logger.Logger.
func (m *MockLogger) Event(e *debug.Event) {
	// Store the event for later assertions
	m.events = append(m.events, e)
}

// Init implements logger.Logger.
func (m *MockLogger) Init() error {
	// Implement the method or leave it empty if it's not needed for your tests
	return nil
}

func NewMockLogger() *MockLogger {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	core, observed := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)
	return &MockLogger{
		observer: observed,
		Logger:   logger,
		start:    time.Now(),
	}
}

func (m *MockLogger) SetLevel(level zapcore.Level) {
	m.Logger.Core().Enabled(level)
}

// AllEntries returns all logged entries captured by the observer.
func (m *MockLogger) AllEntries() []observer.LoggedEntry {
	return m.observer.AllUntimed()
}
