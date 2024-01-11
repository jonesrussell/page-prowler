package common

import "testing"

func TestManagerKey(t *testing.T) {
	if ManagerKey != 0 {
		t.Errorf("Expected ManagerKey to be 0, got %d", ManagerKey)
	}
}
