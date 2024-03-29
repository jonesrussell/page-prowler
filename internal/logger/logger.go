package logger

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/debug"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the interface for logging functions.
type Logger interface {
	Info(msg string)
	Error(msg string)
	Fatal(msg string)
	Debug(msg string)
	Warn(msg string)
	Init() error
	Event(e *debug.Event)
}

// zapLogger is a concrete implementation of Logger using Zap.
type zapLogger struct {
	logger *zap.Logger
	start  time.Time
}

// New creates a new Logger instance with the given log level.
func New(level zapcore.Level) (*zapLogger, error) {
	config := zap.NewProductionConfig() // Adjust for development if needed
	config.Level.SetLevel(level)

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %v", err)
	}

	return &zapLogger{
		logger: logger,
		start:  time.Now(),
	}, nil
}

// Info logs a message at the Info level.
func (l *zapLogger) Info(msg string) {
	l.logger.Info(msg)
}

// Error logs a message at the Error level.
func (l *zapLogger) Error(msg string) {
	l.logger.Error(msg)
}

// Fatal logs a message at the Fatal level, then exits the process.
func (l *zapLogger) Fatal(msg string) {
	l.logger.Fatal(msg)
}

// Debug logs a message at the Debug level.
func (l *zapLogger) Debug(msg string) {
	l.logger.Debug(msg)
}

// Warn logs a message at the Warn level.
func (l *zapLogger) Warn(msg string) {
	l.logger.Warn(msg)
}

// Init initializes the logger.
func (l *zapLogger) Init() error {
	// Add any custom initialization for Zap here (optional)
	return nil
}

// Event logs a debug event.
func (l *zapLogger) Event(e *debug.Event) {
	l.logger.Debug("Colly Debug Event",
		zap.String("Type", e.Type),
		zap.Uint32("RequestID", e.RequestID),
		zap.Uint32("CollectorID", e.CollectorID),
		zap.Any("Values", e.Values),
		zap.Duration("ElapsedTime", time.Since(l.start)),
	)
}
