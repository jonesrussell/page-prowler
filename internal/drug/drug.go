package drug

import (
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/bbalet/stopwords"
	"github.com/caneroj1/stemmer"
)

type Matcher struct { // Renamed from DrugMatcher to Matcher
	swg *metrics.SmithWatermanGotoh
}

func NewMatcher() *Matcher { // Updated constructor name
	swg := metrics.NewSmithWatermanGotoh()
	swg.CaseSensitive = false
	swg.GapPenalty = -0.1
	swg.Substitution = metrics.MatchMismatch{
		Match:    1,
		Mismatch: -0.5,
	}

	return &Matcher{swg: swg}
}

func (m *Matcher) Match(href string) bool { // Updated receiver name
	// We only want the title part of the URL
	sliced := strings.Split(href, "/")
	title := sliced[len(sliced)-1]

	// Clean up the title if it's empty
	if title == "" {
		if len(sliced)-2 < 0 {
			return false
		}
		title = sliced[len(sliced)-2]
	}

	// Remove hyphens from title
	sliced = strings.Split(title, "-")
	title = strings.Join(sliced, " ")

	// Remove stopwords
	title = stopwords.CleanString(title, "en", false)

	// Trim whitespace
	title = strings.TrimSpace(title)

	// Stem the remaining words
	sliced = strings.Split(title, " ")
	sliced = stemmer.StemMultiple(sliced)

	// Convert slice back to string
	title = strings.Join(sliced, " ")

	if title == "" {
		return false
	}

	// Calculate similarities
	similarityDrug := strutil.Similarity("DRUG", title, m.swg)
	similaritySmokeJoint := strutil.Similarity("SMOKE JOINT", title, m.swg)
	similarityGrowop := strutil.Similarity("GROW OP", title, m.swg)
	similarityCannabi := strutil.Similarity("CANNABI", title, m.swg)
	similarityImpair := strutil.Similarity("IMPAIR", title, m.swg)
	similarityShoot := strutil.Similarity("SHOOT", title, m.swg)
	similarityFirearm := strutil.Similarity("FIREARM", title, m.swg)
	similarityMurder := strutil.Similarity("MURDER", title, m.swg)
	similarityCocain := strutil.Similarity("COCAIN", title, m.swg)
	similarityPossess := strutil.Similarity("POSSESS", title, m.swg)
	similarityBreakEnter := strutil.Similarity("BREAK ENTER", title, m.swg)

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
