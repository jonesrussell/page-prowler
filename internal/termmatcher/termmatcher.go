package termmatcher

import (
	"fmt"
	"strings"

	"github.com/adrg/strutil/metrics"
	"github.com/bbalet/stopwords"
	"github.com/caneroj1/stemmer"
	"github.com/jonesrussell/loggo"
	"github.com/jonesrussell/page-prowler/utils"
)

const minTitleLength = 3

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

func (tm *TermMatcher) GetMatchingTerms(href, anchorText string, searchTerms []string) ([]string, error) {
	content := utils.ExtractLastSegmentFromURL(href)
	processedContent := tm.processContent(content)
	processedAnchor := tm.processContent(anchorText)
	combinedContent := tm.combineContents(processedContent, processedAnchor)

	tm.logger.Debug(fmt.Sprintf("Combined content: %v", combinedContent))

	if len(combinedContent) < minTitleLength {
		return []string{}, nil
	}

	allSearchTerms := tm.flattenSearchTerms(searchTerms)
	matchingTerms := tm.findMatchingTerms(combinedContent, allSearchTerms)

	return tm.removeDuplicates(matchingTerms), nil
}

func (tm *TermMatcher) flattenSearchTerms(searchTerms []string) []string {
	var allSearchTerms []string
	for _, terms := range searchTerms {
		allSearchTerms = append(allSearchTerms, strings.Split(terms, ",")...)
	}
	return allSearchTerms
}

func (tm *TermMatcher) processContent(content string) string {
	if processed, ok := tm.processedCache[content]; ok {
		return processed
	}

	content = strings.ReplaceAll(content, "-", " ")
	content = stopwords.CleanString(content, "en", true)
	words := strings.Fields(content)

	// Debug log before stemming
	tm.logger.Debug(fmt.Sprintf("Before stemming: %v", words))

	words = stemmer.StemMultiple(words)

	// Convert words to lowercase after stemming
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}

	// Debug log after stemming
	tm.logger.Debug(fmt.Sprintf("After stemming: %v", words))

	processed := strings.Join(words, " ")
	tm.processedCache[content] = processed
	return processed
}

func (tm *TermMatcher) combineContents(content1 string, content2 string) string {
	if content2 == "" {
		return content1
	}
	return content1 + " " + content2
}

func (tm *TermMatcher) CompareTerms(searchTerm string, content string) float64 {
	searchTerm = strings.ToLower(searchTerm)
	content = strings.ToLower(content)
	similarity := tm.swg.Compare(searchTerm, content)

	if similarity < tm.similarityThreshold {
		similarity = 0
	}

	tm.logger.Debug(fmt.Sprintf("Compared terms: searchTerm=%s, content=%s, similarity=%.2f", searchTerm, content, similarity))

	return similarity
}

func (tm *TermMatcher) findMatchingTerms(content string, searchTerms []string) []string {
	var matchingTerms []string
	processedContent := tm.stemContent(strings.ToLower(content))

	for _, searchTerm := range searchTerms {
		processedTerm := tm.stemContent(strings.ToLower(searchTerm))
		if strings.Contains(processedContent, processedTerm) {
			matchingTerms = append(matchingTerms, searchTerm)
		} else if tm.CompareTerms(searchTerm, content) >= tm.similarityThreshold {
			matchingTerms = append(matchingTerms, searchTerm)
			tm.logger.Debug(fmt.Sprintf("Matched term '%s' with similarity", searchTerm))
		}
	}

	// Ensure only unique and expected terms are returned
	return tm.removeDuplicates(matchingTerms)
}

func (tm *TermMatcher) stemContent(content string) string {
	words := strings.Fields(content)
	stemmedWords := stemmer.StemMultiple(words)
	return strings.Join(stemmedWords, " ")
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
