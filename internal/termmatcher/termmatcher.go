package termmatcher

import (
	"fmt"
	"strings"

	"github.com/adrg/strutil/metrics"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/utils"
	"github.com/xrash/smetrics"
)

const minTitleLength = 3

type TermMatcher struct {
	logger              loggo.LoggerInterface
	swg                 *metrics.SmithWatermanGotoh
	similarityThreshold float64
	contentProcessor    ContentProcessor
}

func NewTermMatcher(logger loggo.LoggerInterface, threshold float64, processor ContentProcessor) *TermMatcher {
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
		contentProcessor:    processor,
	}
}

func (tm *TermMatcher) GetMatchingTerms(href, anchorText string, searchTerms []string) ([]string, error) {
	content := utils.ExtractLastSegmentFromURL(href)
	processedContent := tm.contentProcessor.Process(content)
	processedAnchor := tm.contentProcessor.Process(anchorText)
	combinedContent := tm.combineContents(processedContent, processedAnchor)

	tm.logger.Debug(fmt.Sprintf("Combined content: %v", combinedContent))
	tm.logger.Debug(fmt.Sprintf("Search terms: %v", searchTerms))

	// Check if the combined content is too short
	if len(strings.Fields(combinedContent)) < minTitleLength {
		tm.logger.Debug("Combined content too short")
		return []string{}, nil
	}

	allSearchTerms := tm.flattenSearchTerms(searchTerms)
	tm.logger.Debug(fmt.Sprintf("Flattened search terms: %v", allSearchTerms))

	matchingTerms := tm.findMatchingTerms(combinedContent, allSearchTerms)
	tm.logger.Debug(fmt.Sprintf("Matching terms: %v", matchingTerms))

	return matchingTerms, nil
}

func (tm *TermMatcher) flattenSearchTerms(searchTerms []string) []string {
	var allSearchTerms []string
	for _, terms := range searchTerms {
		allSearchTerms = append(allSearchTerms, strings.Split(terms, ",")...)
	}
	return allSearchTerms
}

func (tm *TermMatcher) combineContents(content1 string, content2 string) string {
	if content2 == "" {
		return content1
	}
	return content1 + " " + content2
}

func (tm *TermMatcher) findMatchingTerms(content string, searchTerms []string) []string {
	var matchingTerms []string
	processedContent := tm.contentProcessor.Stem(strings.ToLower(content))
	words := strings.Fields(processedContent)

	tm.logger.Debug(fmt.Sprintf("Processed content: %s", processedContent))
	tm.logger.Debug(fmt.Sprintf("Words: %v", words))

	for _, searchTerm := range searchTerms {
		processedTerm := tm.contentProcessor.Stem(strings.ToLower(searchTerm))
		tm.logger.Debug(fmt.Sprintf("Processing search term: %s (stemmed: %s)", searchTerm, processedTerm))

		if strings.Contains(processedContent, processedTerm) {
			tm.logger.Debug(fmt.Sprintf("Exact match found for: %s", searchTerm))
			matchingTerms = append(matchingTerms, searchTerm)
		} else {
			if tm.isMultiTerm(processedTerm) {
				matchingTerms = append(matchingTerms, tm.compareMultiTerm(processedTerm, words)...)
			} else {
				matchingTerms = append(matchingTerms, tm.compareSingleTerm(processedTerm, words)...)
			}
		}
	}

	return tm.removeDuplicates(matchingTerms)
}

func (tm *TermMatcher) isMultiTerm(term string) bool {
	return len(strings.Fields(term)) > 1
}

func (tm *TermMatcher) compareSingleTerm(term string, words []string) []string {
	var matchingTerms []string
	for _, word := range words {
		similarity := tm.CompareTerms(term, word)
		tm.logger.Debug(fmt.Sprintf("Comparing single term '%s' with '%s', similarity: %.2f", term, word, similarity))
		if similarity >= tm.similarityThreshold {
			tm.logger.Debug(fmt.Sprintf("Similarity match found for: %s", term))
			matchingTerms = append(matchingTerms, term)
			break
		}
	}
	return matchingTerms
}

func (tm *TermMatcher) compareMultiTerm(term string, words []string) []string {
	matchingTerms := []string{} // Initialize as an empty slice
	termWords := strings.Fields(term)
	termLength := len(termWords)

	for i := 0; i <= len(words)-termLength; i++ {
		phrase := strings.Join(words[i:i+termLength], " ")
		similarity := tm.CompareTerms(term, phrase)
		tm.logger.Debug(fmt.Sprintf("Comparing multi-term '%s' with '%s', similarity: %.2f", term, phrase, similarity))
		if similarity >= tm.similarityThreshold {
			tm.logger.Debug(fmt.Sprintf("Similarity match found for: %s", term))
			matchingTerms = append(matchingTerms, term)
			break
		}
	}
	return matchingTerms
}

func (tm *TermMatcher) removeDuplicates(terms []string) []string {
	seen := make(map[string]struct{})
	unique := []string{}
	for _, term := range terms {
		if _, ok := seen[term]; !ok {
			seen[term] = struct{}{}
			unique = append(unique, term)
		}
	}
	return unique
}

func (tm *TermMatcher) CompareTerms(term1, term2 string) float64 {
	similarity := smetrics.JaroWinkler(strings.ToLower(term1), strings.ToLower(term2), 0.7, 4)
	tm.logger.Debug(fmt.Sprintf("Compared terms: term1=%s, term2=%s, similarity=%.2f", term1, term2, similarity))
	return similarity
}
