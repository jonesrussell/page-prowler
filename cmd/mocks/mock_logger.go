package mocks

type MockLogger struct{}

func (m *MockLogger) Info(msg string, keysAndValues ...interface{}) {
	// Implement your mock logic here
}

func (m *MockLogger) Debug(msg string, keysAndValues ...interface{}) {
	// Implement your mock logic here
}

func (m *MockLogger) Error(msg string, keysAndValues ...interface{}) {
	// Implement your mock logic here
}

func (m *MockLogger) Fatal(msg string, keysAndValues ...interface{}) {
	// Implement your mock logic here
}

func (m *MockLogger) IsDebugEnabled() bool {
	// Implement your mock logic here
	return false
}

func (m *MockLogger) Warn(msg string, keysAndValues ...interface{}) {
	// Implement your mock logic here
}
