package matcher

import (
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/bbalet/stopwords"
	"github.com/caneroj1/stemmer"
)

type BaseMatcher struct {
	swg *metrics.SmithWatermanGotoh // Unexported field
}

func NewBaseMatcher() *BaseMatcher {
	swg := metrics.NewSmithWatermanGotoh()
	swg.CaseSensitive = false
	swg.GapPenalty = -0.1
	swg.Substitution = metrics.MatchMismatch{
		Match:    1,
		Mismatch: -0.5,
	}

	return &BaseMatcher{swg: swg}
}

func (bm *BaseMatcher) ProcessContent(content string) string {
	content = strings.ReplaceAll(content, "-", " ")                         // Remove hyphens
	content = strings.TrimSpace(stopwords.CleanString(content, "en", true)) // Remove stopwords

	// Process and stem
	content = strings.ToLower(content)
	words := strings.Fields(content)
	words = stemmer.StemMultiple(words)
	return strings.ToLower(strings.Join(words, " "))
}

// New method to perform similarity check
func (bm *BaseMatcher) Similarity(term1, term2 string) float64 {
	return strutil.Similarity(term1, term2, bm.swg)
}
