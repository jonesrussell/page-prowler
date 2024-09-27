package crawler

type MockStatsManager struct {
	// Add fields to store the state of the MockStatsManager if needed
}

// Implement the methods of StatsManager with the behavior you want in your tests.
// For example, if you have a method called IncrementTotalLinks, you could implement it like this:

func (m *MockStatsManager) IncrementTotalLinks() {
	// In this mock method, we do nothing.
	// In your tests, you can check if this method was called by adding a field to MockStatsManager
	// and incrementing it here.
}

// Implement other methods of StatsManager...
