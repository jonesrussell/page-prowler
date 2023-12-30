package logger

import (
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
	IsDebugEnabled() bool
}

type ZapLoggerWrapper struct {
	Logger *zap.SugaredLogger
}

func (z *ZapLoggerWrapper) Info(msg string, keysAndValues ...interface{}) {
	z.Logger.Infow(msg, keysAndValues...)
}

func (z *ZapLoggerWrapper) Error(msg string, keysAndValues ...interface{}) {
	z.Logger.Errorw(msg, keysAndValues...)
}

func (z *ZapLoggerWrapper) Fatal(msg string, keysAndValues ...interface{}) {
	z.Logger.Fatalw(msg, keysAndValues...)
}

func (z *ZapLoggerWrapper) Debug(msg string, keysAndValues ...interface{}) {
	z.Logger.Debugw(msg, keysAndValues...)
}

func (z *ZapLoggerWrapper) Warn(msg string, keysAndValues ...interface{}) {
	z.Logger.Warnw(msg, keysAndValues...)
}

func (z *ZapLoggerWrapper) IsDebugEnabled() bool {
	return z.Logger.Desugar().Core().Enabled(zap.DebugLevel)
}

// New returns a new Logger instance.
func New(debug bool, level zapcore.Level) *ZapLoggerWrapper {
	var logger *zap.Logger
	var err error

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	if debug {
		config := zap.Config{
			Level:            atomicLevel,
			Development:      true,
			Encoding:         "console",
			EncoderConfig:    encoderConfig,
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}
		logger, err = config.Build()
	} else {
		config := zap.Config{
			Level:            atomicLevel,
			Development:      false,
			Encoding:         "json",
			EncoderConfig:    encoderConfig,
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}
		logger, err = config.Build()
	}

	if err != nil {
		panic(err)
	}

	return &ZapLoggerWrapper{logger.Sugar()}
}
