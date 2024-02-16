package mocks

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type MockLogger struct {
	observer *observer.ObservedLogs
	Logger   *zap.SugaredLogger
}

func NewMockLogger() *MockLogger {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	core, observed := observer.New(zapcore.DebugLevel)
	sugar := zap.New(core).Sugar()
	return &MockLogger{
		observer: observed,
		Logger:   sugar,
	}
}

func (m *MockLogger) SetLevel(level zapcore.Level) {
	m.Logger.Desugar().Core().Enabled(level)
}

func (m *MockLogger) Info(msg string, keysAndValues ...interface{}) {
	m.Logger.Infow(msg, keysAndValues...)
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Logger.Infof(format, args...)
}

func (m *MockLogger) Debug(msg string, keysAndValues ...interface{}) {
	m.Logger.Debugw(msg, keysAndValues...)
}

func (m *MockLogger) Error(msg string, keysAndValues ...interface{}) {
	m.Logger.Errorw(msg, keysAndValues...)
}

func (m *MockLogger) Fatal(msg string, keysAndValues ...interface{}) {
	m.Logger.Fatalw(msg, keysAndValues...)
}

func (m *MockLogger) Fatalf(msg string, keysAndValues ...interface{}) {
	m.Logger.Fatalf(msg, keysAndValues...)
}

func (m *MockLogger) Warn(msg string, keysAndValues ...interface{}) {
	m.Logger.Warnw(msg, keysAndValues...)
}

func (m *MockLogger) AllEntries() []observer.LoggedEntry {
	return m.observer.AllUntimed()
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Logger.Errorf(format, args...)
}
