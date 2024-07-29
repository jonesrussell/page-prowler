package termmatcher

import (
	"fmt"
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/bbalet/stopwords"
	"github.com/caneroj1/stemmer"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/utils"
)

const minTitleLength = 5 // Set the minimum character limit as needed

type TermMatcher struct {
	logger loggo.LoggerInterface
	swg    *metrics.SmithWatermanGotoh
}

func NewTermMatcher(logger loggo.LoggerInterface) *TermMatcher {
	swg := metrics.NewSmithWatermanGotoh()
	swg.CaseSensitive = false
	swg.GapPenalty = -0.1
	swg.Substitution = metrics.MatchMismatch{
		Match:    1,
		Mismatch: -0.5,
	}

	return &TermMatcher{
		logger: logger,
		swg:    swg,
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

	matchingTerms := tm.findMatchingTerms(combinedContent, allSearchTerms)

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

func (tm *TermMatcher) compareTerms(searchTerm string, content string) float64 {
	searchTerm = strings.ToLower(searchTerm)
	similarity := strutil.Similarity(searchTerm, content, tm.swg)

	tm.logger.Debug(fmt.Sprintf("Compared terms: searchTerm=%s, content=%s, similarity=%.2f", searchTerm, content, similarity))

	return similarity
}

func (tm *TermMatcher) compareAndAppendTerm(searchTerm string, content string) bool {
	similarity := tm.compareTerms(searchTerm, content)
	if similarity >= 0.9 { // Increase the threshold to 0.9
		tm.logger.Debug(fmt.Sprintf("Matching term found: %v", searchTerm))
		return true
	}
	return false
}

func (tm *TermMatcher) findMatchingTerms(content string, searchTerms []string) []string {
	var matchingTerms []string

	content = tm.convertToLowercase(content)
	contentStemmed := tm.stemContent(content)

	tm.logger.Debug(fmt.Sprintf("Stemmed content: %v", contentStemmed))

	for _, searchTerm := range searchTerms {
		searchTerm = tm.convertToLowercase(searchTerm)
		searchTermStemmed := tm.stemContent(searchTerm)

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
