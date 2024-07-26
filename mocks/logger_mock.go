package mocks

import (
	"fmt"
	"os"

	"github.com/jonesrussell/loggo"
)

type MockLogger struct {
	// Removed unused 'events' field
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+msg+"\n", args...)
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	fmt.Printf("[INFO] "+msg+"\n", args...)
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	fmt.Printf("[WARN] "+msg+"\n", args...)
}

func (m *MockLogger) Error(msg string, err error, args ...interface{}) {
	fmt.Printf("[ERROR] "+msg+": %v\n", append(args, err)...)
}

func (m *MockLogger) Fatal(msg string, err error, args ...interface{}) {
	fmt.Printf("[FATAL] "+msg+": %v\n", append(args, err)...)
	os.Exit(1)
}

func (m *MockLogger) WithOperation(_ string) loggo.LoggerInterface {
	// Renamed 'operationID' to '_'
	return m
}

func NewMockLogger() loggo.LoggerInterface {
	return &MockLogger{}
}
