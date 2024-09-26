package mining

import (
	"strings"

	"github.com/jonesrussell/page-prowler/internal/matcher"
)

const SimilarityThreshold = 1

type Matcher struct {
	*matcher.BaseMatcher
}

func NewMatcher() *Matcher {
	return &Matcher{BaseMatcher: matcher.NewBaseMatcher()}
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

	// Define mining-related terms
	miningTerms := []string{
		"mining", "gold", "silver", "copper", "coal", "ore", "excavation", "drilling",
		"exploration", "mineral", "quarry", "sustainability", "environmental impact",
		"resource", "extraction", "geology", "mineral rights", "mine safety",
		"junior mining stocks", "gold mining", "silver mining", "copper mining",
		"lead mining", "zinc mining", "exploration projects", "mining companies",
		"market data", "stock quotes", "real-time news", "mining sectors",
		"mining regions", "high-grade deposits", "mining districts",
	}

	// Check for matches
	for _, term := range miningTerms {
		if m.Similarity(term, title) >= SimilarityThreshold {
			return true
		}
	}

	return false
}
