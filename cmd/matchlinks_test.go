package cmd

import (
	"testing"
)

func TestMatchlinksCmd(t *testing.T) {
	if matchlinksCmd == nil {
		t.Errorf("matchlinksCmd is not initialized")
	}
}
