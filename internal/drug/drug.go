package drug

import (
	"strings"

	"github.com/adrg/strutil/metrics"
	"github.com/jonesrussell/page-prowler/internal/matcher"
)

const SimilarityThreshold = 1

type Matcher struct {
	*matcher.BaseMatcher
}

func NewMatcher(swg *metrics.SmithWatermanGotoh) *Matcher {
	return &Matcher{BaseMatcher: matcher.NewBaseMatcher(swg)}
}

func (m *Matcher) Match(href string) bool {
	// Extract the title part of the URL
	sliced := strings.Split(href, "/")
	title := sliced[len(sliced)-1]

	// Clean up the title if it's empty
	if title == "" {
		if len(sliced)-2 < 0 {
			return false
		}
		title = sliced[len(sliced)-2]
	}

	// Use BaseMatcher's ProcessContent method
	title = m.ProcessContent(title)

	if title == "" {
		return false
	}

	// Define drug-related terms
	drugTerms := []string{
		"drug", "smoke joint", "prescription", "medication", "pharmacy", "medicine",
		"treatment", "health", "wellness", "pharmaceutical", "dosage", "side effects",
		"prescription drug", "over the counter", "drug interaction", "drug-abuse",
		"drug addiction", "drug rehabilitation", "drug policy", "drug regulation",
	}

	// Check for matches
	for _, term := range drugTerms {
		if m.Similarity(term, title) >= SimilarityThreshold {
			return true
		}
	}

	return false
}
