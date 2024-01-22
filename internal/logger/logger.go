package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the interface for logging functions.
type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type ZapLoggerWrapper struct {
	Logger *zap.SugaredLogger
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

const DefaultLogLevel = InfoLevel

func (z *ZapLoggerWrapper) Info(msg string, keysAndValues ...interface{}) {
	z.Logger.Infow(msg, keysAndValues...)
}

func (z *ZapLoggerWrapper) Error(msg string, keysAndValues ...interface{}) {
	z.Logger.Errorw(msg, keysAndValues...)
}

func (z *ZapLoggerWrapper) Fatal(msg string, keysAndValues ...interface{}) {
	z.Logger.Fatalw(msg, keysAndValues...)
}

func (z *ZapLoggerWrapper) Fatalf(format string, args ...interface{}) {
	z.Logger.Fatalf(format, args...)
}

func (z *ZapLoggerWrapper) Debug(msg string, keysAndValues ...interface{}) {
	z.Logger.Debugw(msg, keysAndValues...)
}

func (z *ZapLoggerWrapper) Warn(msg string, keysAndValues ...interface{}) {
	z.Logger.Warnw(msg, keysAndValues...)
}

func (z *ZapLoggerWrapper) Infof(format string, args ...interface{}) {
	z.Logger.Infof(format, args...)
}

func (z *ZapLoggerWrapper) Errorf(format string, args ...interface{}) {
	z.Logger.Errorf(format, args...)
}

// New returns a new Logger instance.
func New(level LogLevel) (*ZapLoggerWrapper, error) {
	fmt.Printf("Initializing logger with level: %v\n", level)

	var zapLevel zapcore.Level
	switch level {
	case DebugLevel:
		zapLevel = zapcore.DebugLevel
	case InfoLevel:
		zapLevel = zapcore.InfoLevel
	// ... handle other levels if necessary
	default:
		zapLevel = zapcore.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapLevel) // Set the atomic level using the explicit zapcore.Level

	fmt.Printf("Atomic level set to: %v\n", atomicLevel)

	config := zap.Config{
		Level:            atomicLevel,
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %v", err)
	}

	return &ZapLoggerWrapper{Logger: logger.Sugar()}, nil
}
