package termmatcher

import (
	"fmt"
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/bbalet/stopwords"
	"github.com/caneroj1/stemmer"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/internal/matcher" // Import the matcher interface
	"github.com/jonesrussell/page-prowler/utils"
)

const (
	minTitleLength      = 5
	similarityThreshold = 0.9
)

type TermMatcher struct {
	logger   loggo.LoggerInterface
	swg      *metrics.SmithWatermanGotoh
	matchers []matcher.Matcher // List of matchers
}

func NewTermMatcher(logger loggo.LoggerInterface, matchers []matcher.Matcher) *TermMatcher {
	swg := metrics.NewSmithWatermanGotoh()
	swg.CaseSensitive = false
	swg.GapPenalty = -0.1
	swg.Substitution = metrics.MatchMismatch{
		Match:    1,
		Mismatch: -0.5,
	}

	return &TermMatcher{
		logger:   logger,
		swg:      swg,
		matchers: matchers, // Initialize with provided matchers
	}
}

func (tm *TermMatcher) GetMatchingTerms(href string, anchorText string, searchTerms []string) []string {
	content := utils.ExtractLastSegmentFromURL(href)
	processedContent := tm.processContent(content)
	tm.logger.Debug(fmt.Sprintf("Processed content from URL: %v", processedContent))

	anchorContent := tm.processContent(anchorText)
	tm.logger.Debug(fmt.Sprintf("Processed anchor text: %v", anchorContent))

	combinedContent := tm.combineContents(processedContent, anchorContent)
	tm.logger.Debug(fmt.Sprintf("Combined content: %v", combinedContent))

	if len(combinedContent) < minTitleLength {
		tm.logger.Debug(fmt.Sprintf("Combined content is less than minimum title length: %d", minTitleLength))
		return []string{}
	}

	var allSearchTerms []string
	for _, terms := range searchTerms {
		allSearchTerms = append(allSearchTerms, strings.Split(terms, ",")...)
	}

	// Check each matcher for matches
	var matchingTerms []string
	for _, m := range tm.matchers {
		matched, err := m.Match(combinedContent, "")
		if err != nil {
			tm.logger.Error("Error matching term", err)
			continue // Skip to the next matcher if there's an error
		}
		if matched {
			matchingTerms = append(matchingTerms, allSearchTerms...) // Add search terms if matched
		}
	}

	// Use findMatchingTerms to check for additional matches
	matchingTerms = append(matchingTerms, tm.findMatchingTerms(combinedContent, allSearchTerms)...)

	// Remove duplicates
	seen := make(map[string]bool)
	var result []string
	for _, v := range matchingTerms {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	if len(result) == 0 {
		tm.logger.Debug("No matching terms found")
		return []string{}
	}

	tm.logger.Debug(fmt.Sprintf("Found matching terms: %v", matchingTerms))
	tm.logger.Debug(fmt.Sprintf("Matching terms result: %v", result))
	return result
}

func (tm *TermMatcher) findMatchingTerms(content string, searchTerms []string) []string {
	var matchingTerms []string

	content = tm.convertToLowercase(content)
	contentStemmed := tm.stemContent(content)

	tm.logger.Debug(fmt.Sprintf("Stemmed content: %v", contentStemmed))

	for _, searchTerm := range searchTerms {
		searchTerm = tm.convertToLowercase(searchTerm)
		searchTermStemmed := searchTerm // Avoid stemming the search term

		words := strings.Fields(searchTermStemmed)
		for _, word := range words {
			if tm.compareAndAppendTerm(word, contentStemmed) {
				matchingTerms = append(matchingTerms, word)
			}
		}
	}

	if len(matchingTerms) == 0 {
		return []string{}
	}

	tm.logger.Debug(fmt.Sprintf("Matching terms result: %v", matchingTerms))
	return matchingTerms
}

// New CompareTerms method
func (tm *TermMatcher) CompareTerms(term1, term2 string) float64 {
	return strutil.Similarity(term1, term2, tm.swg)
}

func (tm *TermMatcher) processContent(content string) string {
	content = strings.ReplaceAll(content, "-", " ")                         // Remove hyphens
	content = strings.TrimSpace(stopwords.CleanString(content, "en", true)) // Remove stopwords

	// Process and stem
	content = strings.ToLower(content)
	words := strings.Fields(content)
	words = stemmer.StemMultiple(words)
	return strings.ToLower(strings.Join(words, " "))
}

func (tm *TermMatcher) combineContents(content1 string, content2 string) string {
	if content2 == "" {
		return content1
	}
	return content1 + " " + content2
}

func (tm *TermMatcher) compareAndAppendTerm(searchTerm string, content string) bool {
	// Check for exact match
	words := strings.Fields(content)
	for _, word := range words {
		if word == searchTerm {
			tm.logger.Debug(fmt.Sprintf("Exact matching term found: %v", searchTerm))
			return true
		}
	}

	// If no exact match, use SWG for comparison
	similarity := tm.CompareTerms(searchTerm, content)
	if similarity >= similarityThreshold { // Use constant for threshold
		tm.logger.Debug(fmt.Sprintf("Matching term found: %v", searchTerm))
		return true
	}
	return false
}

func (tm *TermMatcher) convertToLowercase(content string) string {
	return strings.ToLower(content)
}

func (tm *TermMatcher) stemContent(content string) string {
	words := strings.Fields(content)
	stemmedWords := stemmer.StemMultiple(words)
	lowercaseStemmedWords := make([]string, len(stemmedWords))
	for i, word := range stemmedWords {
		lowercaseStemmedWords[i] = strings.ToLower(word)
	}
	return strings.Join(lowercaseStemmedWords, " ")
}
