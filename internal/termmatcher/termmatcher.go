package termmatcher

import (
	"fmt"
	"math"
	"strings"

	"github.com/adrg/strutil/metrics"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/utils"
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

	// Check if the combined content is too short
	if len(strings.Fields(combinedContent)) < minTitleLength {
		return []string{}, nil
	}

	allSearchTerms := tm.flattenSearchTerms(searchTerms)
	return tm.findMatchingTerms(combinedContent, allSearchTerms), nil
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

func (tm *TermMatcher) CompareTerms(term1, term2 string) float64 {
	term1 = strings.ToLower(term1)
	term2 = strings.ToLower(term2)

	similarity := tm.swg.Compare(term1, term2)

	// Normalize the similarity score
	maxLen := float64(max(len(term1), len(term2)))
	if maxLen > 0 {
		similarity /= maxLen
	}

	// Adjust similarity to be more lenient
	similarity = math.Pow(similarity, 0.5)

	if similarity < tm.similarityThreshold {
		similarity = 0
	}

	tm.logger.Debug(fmt.Sprintf("Compared terms: term1=%s, term2=%s, similarity=%.2f", term1, term2, similarity))

	return similarity
}

func (tm *TermMatcher) findMatchingTerms(content string, searchTerms []string) []string {
	var matchingTerms []string
	processedContent := tm.contentProcessor.Stem(strings.ToLower(content))

	for _, searchTerm := range searchTerms {
		processedTerm := tm.contentProcessor.Stem(strings.ToLower(searchTerm))
		if strings.Contains(processedContent, processedTerm) {
			matchingTerms = append(matchingTerms, searchTerm)
		} else {
			words := strings.Fields(processedContent)
			for _, word := range words {
				if similarity := tm.CompareTerms(processedTerm, word); similarity >= tm.similarityThreshold {
					matchingTerms = append(matchingTerms, searchTerm)
					break
				}
			}
		}
	}

	return tm.removeDuplicates(matchingTerms)
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
