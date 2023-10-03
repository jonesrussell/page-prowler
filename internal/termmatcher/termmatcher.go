package termmatcher

import (
	"net/url"
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/bbalet/stopwords"
	"github.com/caneroj1/stemmer"
)

// Related checks if the URL title matches any of the provided search terms.
func Related(href string, searchTerms []string) bool {
	title := extractTitleFromURL(href)
	if title == "" {
		return false
	}

	processedTitle := processTitle(title)
	if processedTitle == "" {
		return false
	}

	return matchSearchTerms(processedTitle, searchTerms)
}

// extractTitleFromURL extracts the title part from a URL.
func extractTitleFromURL(urlString string) string {
	u, err := url.Parse(urlString)
	if err != nil {
		// Handle the error, e.g., log it or return an error string
		return ""
	}

	// Check if the URL has a path component
	if u.Path == "" || u.Path == "/" {
		// If there's no path component, return an empty string
		return ""
	}

	// Split the path and get the last segment as the title
	segments := strings.Split(u.Path, "/")
	title := segments[len(segments)-1]

	return title
}

// processTitle processes the title by removing '-', stopwords, and stemming.
func processTitle(title string) string {
	title = strings.ReplaceAll(title, "-", " ")
	title = stopwords.CleanString(title, "en", false)
	title = strings.TrimSpace(title)

	words := strings.Split(title, " ")
	words = stemmer.StemMultiple(words)

	// Lemmatize (if needed)
	// ...

	return strings.Join(words, " ")
}

// matchSearchTerms checks if the processed title matches any of the search terms.
func matchSearchTerms(title string, searchTerms []string) bool {
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
