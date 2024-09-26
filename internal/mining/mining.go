package mining

import (
	"fmt"
	"strings"

	"github.com/adrg/strutil/metrics"
	"github.com/jonesrussell/page-prowler/internal/matcher"
)

const SimilarityThreshold = 0.6

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

	// Define mining-related terms
	miningTerms := []string{
		"mining", "gold", "silver", "copper", "coal", "ore", "excavation", "drilling",
		"exploration", "mineral", "quarry", "sustainability", "environmental impact",
		"resource", "extraction", "geology", "mineral rights", "mine safety",
		"junior mining stocks", "gold mining", "silver mining", "copper mining",
		"lead mining", "zinc mining", "exploration projects", "mining companies",
		"market data", "stock quotes", "real-time news", "mining sectors",
		"mining regions", "high-grade deposits", "mining districts", "mining news",
	}

	// Define terms that should automatically score 0
	excludedTerms := []string{
		"technology", "sports", "entertainment", "fashion", "music", "movies",
	}

	// Check for excluded terms
	for _, term := range excludedTerms {
		if strings.Contains(title, term) {
			fmt.Printf("Excluding term '%s' found in title '%s'\n", term, title)
			return false
		}
	}

	// Check for matches
	for _, term := range miningTerms {
		score := m.Similarity(term, title)
		fmt.Printf("Matching '%s' with '%s': score = %f\n", term, title, score)
		if score >= SimilarityThreshold {
			return true
		}
	}

	return false
}
