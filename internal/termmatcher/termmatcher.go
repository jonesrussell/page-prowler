package termmatcher

import (
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/bbalet/stopwords"
	"github.com/caneroj1/stemmer"
)

// Related checks if the URL title matches any of the provided search terms.
func Related(href string, searchTerms []string) bool {
	// We only want the title part of the URL
	sliced := strings.Split(href, "/")
	title := sliced[len(sliced)-1]

	// If title is empty, use the previous part of the URL as the title
	if title == "" {
		if len(sliced)-2 < 0 {
			return false
		}
		title = sliced[len(sliced)-2]
	}

	// Remove '-' from title
	title = strings.ReplaceAll(title, "-", " ")

	// Remove stopwords
	title = stopwords.CleanString(title, "en", false)

	// Trim
	title = strings.TrimSpace(title)

	// Stem the remaining words
	words := strings.Split(title, " ")
	words = stemmer.StemMultiple(words)

	// Lemmatize (if needed)
	// ...

	// Convert slice back to string
	title = strings.Join(words, " ")

	if title == "" {
		return false
	}

	swg := metrics.NewSmithWatermanGotoh()
	swg.CaseSensitive = false
	swg.GapPenalty = -0.1
	swg.Substitution = metrics.MatchMismatch{
		Match:    1,
		Mismatch: -0.5,
	}

	for _, term := range searchTerms {
		similarity := strutil.Similarity(term, title, swg)
		if similarity == 1 {
			return true
		}
	}

	return false
}
