package cmd

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestArticlesCmd(t *testing.T) {
	if matchlinksCmd == nil {
		t.Errorf("matchlinksCmd is not initialized")
	}
}

func TestArticlesCmdFlags(t *testing.T) {
	// Set the flags
	if err := matchlinksCmd.Flags().Set("crawlsiteid", "test"); err != nil {
		t.Fatalf("Error setting crawlsiteid flag: %v", err)
	}
	if err := matchlinksCmd.Flags().Set("searchterms", "test"); err != nil {
		t.Fatalf("Error setting searchterms flag: %v", err)
	}
	if err := matchlinksCmd.Flags().Set("url", "test"); err != nil {
		t.Fatalf("Error setting url flag: %v", err)
	}

	// Check if the flags are correctly set
	assert.Equal(t, "test", viper.GetString("crawlsiteid"))
	assert.Equal(t, "test", viper.GetString("searchterms"))
	assert.Equal(t, "test", viper.GetString("url"))
}
