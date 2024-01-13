package termmatcher

import (
	"net/url"
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/bbalet/stopwords"
	"github.com/caneroj1/stemmer"
)

const minTitleLength = 5 // Set the minimum character limit as needed

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
func removeHyphens(title string) string {
	return strings.ReplaceAll(title, "-", " ")
}

func removeStopwords(title string) string {
	return strings.TrimSpace(stopwords.CleanString(title, "en", true))
}

func stemTitle(title string) string {
	title = strings.ToLower(title)
	words := strings.Fields(title)
	words = stemmer.StemMultiple(words)
	return strings.ToLower(strings.Join(words, " "))
}

func processTitle(title string) string {
	title = removeHyphens(title)
	title = removeStopwords(title)
	title = stemTitle(title)
	title = strings.ToLower(title)
	return title
}

// GetMatchingTerms checks if the URL title matches any of the provided search terms and returns the matching terms.
func GetMatchingTerms(href string, anchorText string, searchTerms []string) []string {
	combined := href + " " + anchorText
	title := extractTitleFromURL(combined)
	if title == "" {
		return nil
	}

	processedTitle := processTitle(title)
	if processedTitle == "" {
		return nil
	}

	// Check if the title meets the minimum character limit
	if len(processedTitle) < minTitleLength {
		return nil
	}

	return findMatchingTerms(processedTitle, searchTerms)
}

// findMatchingTerms finds the search terms that match the given title.
func findMatchingTerms(title string, searchTerms []string) []string {
	var matchingTerms []string
	swg := metrics.NewSmithWatermanGotoh()
	swg.CaseSensitive = false
	swg.GapPenalty = -0.1
	swg.Substitution = metrics.MatchMismatch{
		Match:    1,
		Mismatch: -0.5,
	}

	title = strings.ToLower(title)
	titleStemmed := stemmer.Stem(title)

	for _, term := range searchTerms {
		originalTerm := term
		term = strings.ToLower(term)
		termStemmed := stemmer.Stem(term)
		similarity := strutil.Similarity(termStemmed, titleStemmed, swg)
		if similarity >= 0.8 { // Adjust the threshold as needed
			matchingTerms = append(matchingTerms, originalTerm)
		}
	}

	return matchingTerms
}
