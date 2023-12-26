package cmd

import (
	"testing"
)

func TestArticlesCmd(t *testing.T) {
	if articlesCmd == nil {
		t.Errorf("articlesCmd is not initialized")
	}
}

func TestSaveResultsToRedis(t *testing.T) {
	// TODO: Add your test code here
}

func TestPrintResults(t *testing.T) {
	// TODO: Add your test code here
}
