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
}

func fieldsToZapFields(fields map[string]interface{}) []zapcore.Field {
	zapFields := make([]zapcore.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

func (m *MockLogger) Info(msg string, fields map[string]interface{}) {
	m.Logger.Info(msg, fieldsToZapFields(fields)...)
}

func (m *MockLogger) Debug(msg string, fields map[string]interface{}) {
	m.Logger.Debug(msg, fieldsToZapFields(fields)...)
}

func (m *MockLogger) Error(msg string, fields map[string]interface{}) {
	m.Logger.Error(msg, fieldsToZapFields(fields)...)
}

func (m *MockLogger) Warn(msg string, fields map[string]interface{}) {
	m.Logger.Warn(msg, fieldsToZapFields(fields)...)
}

func (m *MockLogger) Fatal(msg string, fields map[string]interface{}) {
	m.Logger.Fatal(msg, fieldsToZapFields(fields)...)
}

// Event implements logger.Logger.
func (*MockLogger) Event(e *debug.Event) {
	panic("unimplemented")
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

func (m *MockLogger) AllEntries() []observer.LoggedEntry {
	return m.observer.AllUntimed()
}
