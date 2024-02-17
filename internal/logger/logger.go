package logger

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/debug"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the interface for logging functions and implements the debug.Debugger interface.
type Logger interface {
	Info(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})
	Debug(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Init() error
	Event(e *debug.Event)
}

// LoggerWrapper is a wrapper around zap.Logger that implements the debug.Debugger interface.
type LoggerWrapper struct {
	Logger *zap.Logger
	start  time.Time
}

// ConvertFields converts a map of fields to a slice of zapcore.Field.
func (lw *LoggerWrapper) ConvertFields(fields map[string]interface{}) []zapcore.Field {
	zapFields := make([]zapcore.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

// Debug implements Logger.
func (lw *LoggerWrapper) Debug(msg string, fields map[string]interface{}) {
	lw.Logger.Debug(msg, lw.ConvertFields(fields)...)
}

// Error implements Logger.
func (lw *LoggerWrapper) Error(msg string, fields map[string]interface{}) {
	lw.Logger.Error(msg, lw.ConvertFields(fields)...)
}

// Fatal implements Logger.
func (lw *LoggerWrapper) Fatal(msg string, fields map[string]interface{}) {
	lw.Logger.Fatal(msg, lw.ConvertFields(fields)...)
}

// Info implements Logger.
func (lw *LoggerWrapper) Info(msg string, fields map[string]interface{}) {
	lw.Logger.Info(msg, lw.ConvertFields(fields)...)
}

// Warn implements Logger.
func (lw *LoggerWrapper) Warn(msg string, fields map[string]interface{}) {
	lw.Logger.Warn(msg, lw.ConvertFields(fields)...)
}

// Init initializes the LoggerWrapper.
func (lw *LoggerWrapper) Init() error {
	// No initialization needed for LoggerWrapper
	return nil
}

// Event logs a debug event.
func (lw *LoggerWrapper) Event(e *debug.Event) {
	lw.Logger.Debug("Colly Debug Event",
		zap.String("Type", e.Type),
		zap.Uint32("RequestID", e.RequestID),
		zap.Uint32("CollectorID", e.CollectorID),
		zap.Any("Values", e.Values),
		zap.Duration("ElapsedTime", time.Since(lw.start)),
	)
}

type LogLevel int

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel LogLevel = iota
	// InfoLevel is the default logging priority.
	InfoLevel
	// ... other log levels if needed
)

// New returns a new LoggerWrapper instance that implements the debug.Debugger interface.
func New(level LogLevel) (*LoggerWrapper, error) {
	var config zap.Config
	var zapLevel zapcore.Level

	switch level {
	case DebugLevel:
		config = zap.NewDevelopmentConfig()
		zapLevel = zapcore.DebugLevel
	default:
		config = zap.NewProductionConfig()
		zapLevel = zapcore.InfoLevel
	}

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapLevel) // Set the atomic level using the explicit zapcore.Level

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %v", err)
	}

	// Create a LoggerWrapper that implements the debug.Debugger interface
	loggerWrapper := &LoggerWrapper{
		Logger: logger,
		start:  time.Now(),
	}

	return loggerWrapper, nil
}
