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
	logger              loggo.LoggerInterface
	swg                 *metrics.SmithWatermanGotoh
	similarityThreshold float64
	processedCache      map[string]string
}

func NewTermMatcher(logger loggo.LoggerInterface, threshold float64) *TermMatcher {
	swg := metrics.NewSmithWatermanGotoh()
	swg.CaseSensitive = false
	swg.GapPenalty = -0.1
	swg.Substitution = metrics.MatchMismatch{
		Match:    1,
		Mismatch: -0.5,
	}

	return &TermMatcher{
		logger:              logger,
		swg:                 swg,
		similarityThreshold: threshold,
		processedCache:      make(map[string]string),
	}
}

func (tm *TermMatcher) GetMatchingTerms(href string, anchorText string, searchTerms []string) ([]string, error) {
	content := utils.ExtractLastSegmentFromURL(href)
	processedContent := tm.processContent(content)
	tm.logger.Debug(fmt.Sprintf("Processed content from URL: %v", processedContent))

	anchorContent := tm.processContent(anchorText)
	tm.logger.Debug(fmt.Sprintf("Processed anchor text: %v", anchorContent))

	combinedContent := tm.combineContents(processedContent, anchorContent)
	tm.logger.Debug(fmt.Sprintf("Combined content: %v", combinedContent))

	if len(combinedContent) < minTitleLength {
		tm.logger.Warn(fmt.Sprintf("Combined content is less than minimum title length: %d", minTitleLength))
		return []string{}, fmt.Errorf("combined content too short: %d characters", len(combinedContent))
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
		return []string{}, nil
	}

	tm.logger.Debug(fmt.Sprintf("Found matching terms: %v", matchingTerms))
	tm.logger.Debug(fmt.Sprintf("Matching terms result: %v", result))
	return result, nil
}

func (tm *TermMatcher) processContent(content string) string {
	if processed, ok := tm.processedCache[content]; ok {
		return processed
	}
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

func (tm *TermMatcher) CompareTerms(searchTerm string, content string) float64 {
	searchTerm = strings.ToLower(searchTerm)
	similarity := strutil.Similarity(searchTerm, content, tm.swg)

	tm.logger.Debug(fmt.Sprintf("Compared terms: searchTerm=%s, content=%s, similarity=%.2f", searchTerm, content, similarity))

	return similarity
}

func (tm *TermMatcher) findMatchingTerms(content string, searchTerms []string) []string {
	var matchingTerms []string

	contentWords := strings.Fields(tm.stemContent(tm.convertToLowercase(content)))
	contentSet := make(map[string]struct{}, len(contentWords))
	for _, word := range contentWords {
		contentSet[word] = struct{}{}
	}

	for _, searchTerm := range searchTerms {
		processedTerm := tm.stemContent(tm.convertToLowercase(searchTerm))
		words := strings.Fields(processedTerm)
		for _, word := range words {
			if _, exists := contentSet[word]; exists {
				matchingTerms = append(matchingTerms, searchTerm)
				break
			}
		}
	}

	// Use similarity comparison only for terms not found by exact match
	// ... implement this part ...

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
