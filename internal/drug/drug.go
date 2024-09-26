package drug

import (
	"strings"

	"github.com/jonesrussell/page-prowler/internal/matcher"
)

type Matcher struct {
	*matcher.BaseMatcher // Embed BaseMatcher
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

	// Calculate similarities
	similarityDrug := m.Similarity("DRUG", title)
	similaritySmokeJoint := m.Similarity("SMOKE JOINT", title)
	similarityGrowop := m.Similarity("GROW OP", title)
	similarityCannabi := m.Similarity("CANNABI", title)
	similarityImpair := m.Similarity("IMPAIR", title)
	similarityShoot := m.Similarity("SHOOT", title)
	similarityFirearm := m.Similarity("FIREARM", title)
	similarityMurder := m.Similarity("MURDER", title)
	similarityCocain := m.Similarity("COCAIN", title)
	similarityPossess := m.Similarity("POSSESS", title)
	similarityBreakEnter := m.Similarity("BREAK ENTER", title)

	return similarityDrug == 1 ||
		similaritySmokeJoint == 1 ||
		similarityGrowop == 1 ||
		similarityCannabi == 1 ||
		similarityImpair == 1 ||
		similarityShoot == 1 ||
		similarityFirearm == 1 ||
		similarityMurder == 1 ||
		similarityCocain == 1 ||
		similarityPossess == 1 ||
		similarityBreakEnter == 1
}
