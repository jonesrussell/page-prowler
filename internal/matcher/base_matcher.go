package matcher

import (
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/bbalet/stopwords"
	"github.com/caneroj1/stemmer"
)

type BaseMatcher struct {
	swg *metrics.SmithWatermanGotoh
}

// NewBaseMatcher creates a new BaseMatcher with a provided SmithWatermanGotoh instance.
func NewBaseMatcher(swg *metrics.SmithWatermanGotoh) *BaseMatcher {
	if swg == nil {
		swg = metrics.NewSmithWatermanGotoh() // Default instance if none provided
		swg.CaseSensitive = false
		swg.GapPenalty = -0.1
		swg.Substitution = metrics.MatchMismatch{
			Match:    1,
			Mismatch: -0.5,
		}
	}
	return &BaseMatcher{swg: swg}
}

// ProcessContent processes the content by removing hyphens, stopwords, and stemming.
func (bm *BaseMatcher) ProcessContent(content string) string {
	content = strings.ReplaceAll(content, "-", " ")                         // Remove hyphens
	content = strings.TrimSpace(stopwords.CleanString(content, "en", true)) // Remove stopwords
	return bm.StemAndLowerContent(content)                                  // Combine stemming and lowercasing
}

// StemAndLowerContent stems the content and returns the processed string.
func (bm *BaseMatcher) StemAndLowerContent(content string) string {
	words := strings.Fields(content)
	stemmedWords := stemmer.StemMultiple(words)
	lowercaseStemmedWords := make([]string, len(stemmedWords))
	for i, word := range stemmedWords {
		lowercaseStemmedWords[i] = strings.ToLower(word)
	}
	return strings.Join(lowercaseStemmedWords, " ")
}

// Similarity checks the similarity between two terms.
func (bm *BaseMatcher) Similarity(term1, term2 string) float64 {
	return strutil.Similarity(term1, term2, bm.swg)
}
