package termmatcher

import (
	"net/url"
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/bbalet/stopwords"
	"github.com/caneroj1/stemmer"
	"github.com/jonesrussell/page-prowler/internal/logger"
)

const minTitleLength = 5 // Set the minimum character limit as needed

// ExtractLastSegmentFromURL extracts the title part from a URL.
func ExtractLastSegmentFromURL(urlString string) string {
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

func RemoveHyphens(title string) string {
	return strings.ReplaceAll(title, "-", " ")
}

func RemoveStopwords(title string) string {
	return strings.TrimSpace(stopwords.CleanString(title, "en", true))
}

func ProcessAndStem(content string) string {
	content = strings.ToLower(content)
	words := strings.Fields(content)
	words = stemmer.StemMultiple(words)
	return strings.ToLower(strings.Join(words, " "))
}

func ProcessContent(content string) string {
	content = RemoveHyphens(content)
	content = RemoveStopwords(content)
	content = ProcessAndStem(content)
	return content
}

// GetMatchingTerms checks if the URL title matches any of the provided search terms and returns the matching terms.
func GetMatchingTerms(href string, anchorText string, searchTerms []string, logger logger.Logger) []string {
	content := ExtractLastSegmentFromURL(href)
	processedContent := ProcessContent(content)
	logger.Debug("Processed content from URL", map[string]interface{}{"processedContent": processedContent})

	anchorContent := ProcessContent(anchorText)
	logger.Debug("Processed anchor text", map[string]interface{}{"anchorContent": anchorContent})

	combinedContent := CombineContents(processedContent, anchorContent)
	logger.Debug("Combined content", map[string]interface{}{"combinedContent": combinedContent})

	if len(combinedContent) < minTitleLength {
		logger.Debug("Combined content is less than minimum title length", map[string]interface{}{"minTitleLength": minTitleLength})
		return []string{}
	}

	matchingTerms := FindMatchingTerms(combinedContent, searchTerms, logger)
	logger.Debug("Found matching terms", map[string]interface{}{"matchingTerms": matchingTerms})

	seen := make(map[string]bool)
	var result []string
	for _, v := range matchingTerms {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	// Instead of returning nil, return an empty slice if no matching terms are found
	if len(result) == 0 {
		logger.Debug("No matching terms found", map[string]interface{}{})
		return []string{}
	}

	logger.Debug("Matching terms result", map[string]interface{}{"result": result})
	return result
}

func CombineContents(content1 string, content2 string) string {
	if content2 == "" {
		return content1
	}
	return content1 + " " + content2
}

func ConvertToLowercase(content string) string {
	return strings.ToLower(content)
}

func StemContent(content string) string {
	words := strings.Fields(content)
	stemmedWords := stemmer.StemMultiple(words)
	lowercaseStemmedWords := make([]string, len(stemmedWords))
	for i, word := range stemmedWords {
		lowercaseStemmedWords[i] = strings.ToLower(word)
	}
	return strings.Join(lowercaseStemmedWords, " ")
}

func CompareTerms(searchTerm string, content string, swg *metrics.SmithWatermanGotoh, mylogger logger.Logger) float64 {
	searchTerm = strings.ToLower(searchTerm)
	similarity := strutil.Similarity(searchTerm, content, swg)

	// Log the searchTerm, content, and their similarity
	mylogger.Debug("Compared terms", map[string]interface{}{"searchTerm": searchTerm, "content": content, "similarity": similarity})

	return similarity
}

func CreateSWG() *metrics.SmithWatermanGotoh {
	swg := metrics.NewSmithWatermanGotoh()
	swg.CaseSensitive = false
	swg.GapPenalty = -0.1
	swg.Substitution = metrics.MatchMismatch{
		Match:    1,
		Mismatch: -0.5,
	}
	return swg
}

func CompareAndAppendTerm(searchTerm string, content string, swg *metrics.SmithWatermanGotoh, matchingTerms *[]string, mylogger logger.Logger) {
	similarity := CompareTerms(searchTerm, content, swg, mylogger)
	mylogger.Debug("Compared terms", map[string]interface{}{"searchTerm": searchTerm, "similarity": similarity})
	if similarity >= 0.9 { // Increase the threshold to 0.9
		*matchingTerms = append(*matchingTerms, searchTerm)
		mylogger.Debug("Matching term found", map[string]interface{}{"searchTerm": searchTerm})
	}
}

func FindMatchingTerms(content string, searchTerms []string, mylogger logger.Logger) []string {
	var matchingTerms []string
	swg := CreateSWG()

	content = ConvertToLowercase(content)
	contentStemmed := StemContent(content)

	// Debug statement
	mylogger.Debug("Stemmed content", map[string]interface{}{"contentStemmed": contentStemmed})

	for _, searchTerm := range searchTerms {
		// Convert the search term to lowercase and apply stemming
		searchTerm = ConvertToLowercase(searchTerm)
		searchTermStemmed := StemContent(searchTerm)

		words := strings.Fields(searchTermStemmed)
		for _, word := range words {
			CompareAndAppendTerm(word, contentStemmed, swg, &matchingTerms, mylogger)
		}
	}

	// Ensure an empty slice is returned instead of nil
	if len(matchingTerms) == 0 {
		mylogger.Debug("No matching terms found", map[string]interface{}{})
		return []string{}
	}

	mylogger.Debug("Matching terms result", map[string]interface{}{"matchingTerms": matchingTerms})
	return matchingTerms
}
