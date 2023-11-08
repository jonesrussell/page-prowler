package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the interface for logging functions.
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
}

// zapLogger is an internal struct that will implement the Logger interface using zap.
type zapLogger struct {
	sugar *zap.SugaredLogger
}

// New returns a new Logger instance.
func New(debug bool) Logger {
	var logger *zap.Logger
	var err error

	if debug {
		// Development logger is more verbose and writes to standard output.
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // optional, colorizes the output.
		logger, err = config.Build()
	} else {
		// Production logger is less verbose and could be set to log to a file.
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}

	return &zapLogger{sugar: logger.Sugar()}
}

// Debug logs a message at the Debug level.
func (l *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.sugar.Debugw(msg, keysAndValues...)
}

// Info logs a message at the Info level.
func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.sugar.Infow(msg, keysAndValues...)
}

// Warn logs a message at the Warn level.
func (l *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.sugar.Warnw(msg, keysAndValues...)
}

// Error logs a message at the Error level.
func (l *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.sugar.Errorw(msg, keysAndValues...)
}

// Fatal logs a message at the Fatal level and exits the application.
func (l *zapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.sugar.Fatalw(msg, keysAndValues...)
}
