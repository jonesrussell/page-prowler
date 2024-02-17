package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogger(t *testing.T) {
	// Create an observer to capture log entries
	core, recorded := observer.New(zapcore.DebugLevel)

	// Create a new logger that writes to the buffer
	appLogger, err := New(LogLevel(zapcore.DebugLevel))
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	appLogger.Logger = zap.New(core)

	// Log a message at each level
	appLogger.Debug("debug", map[string]interface{}{"key": "value"})
	appLogger.Info("info", map[string]interface{}{"key": "value"})
	appLogger.Warn("warn", map[string]interface{}{"key": "value"})
	appLogger.Error("error", map[string]interface{}{"key": "value"})

	// Check the recorded log entries
	entries := recorded.All()
	assert.Equal(t, 4, len(entries))
	assert.Equal(t, "debug", entries[0].Message)
	assert.Equal(t, "info", entries[1].Message)
	assert.Equal(t, "warn", entries[2].Message)
	assert.Equal(t, "error", entries[3].Message)
}
