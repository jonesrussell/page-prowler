package mocks

import (
	"fmt"
)

type MockLogger struct{}

func (m *MockLogger) Info(msg string, keysAndValues ...interface{}) {
	fmt.Println("Info: ", msg, keysAndValues)
}

func (m *MockLogger) Debug(msg string, keysAndValues ...interface{}) {
	fmt.Println("Debug: ", msg, keysAndValues)
}

func (m *MockLogger) Error(msg string, keysAndValues ...interface{}) {
	fmt.Println("Error: ", msg, keysAndValues)
}

func (m *MockLogger) Fatal(msg string, keysAndValues ...interface{}) {
	fmt.Println("Fatal: ", msg, keysAndValues)
}

func (m *MockLogger) IsDebugEnabled() bool {
	return true
}

func (m *MockLogger) Warn(msg string, keysAndValues ...interface{}) {
	fmt.Println("Warn: ", msg, keysAndValues)
}
