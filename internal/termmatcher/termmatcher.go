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

func (tm *TermMatcher) RemoveHyphens(title string) string {
	return strings.ReplaceAll(title, "-", " ")
}

func (tm *TermMatcher) RemoveStopwords(title string) string {
	return strings.TrimSpace(stopwords.CleanString(title, "en", true))
}

func (tm *TermMatcher) ProcessAndStem(content string) string {
	content = strings.ToLower(content)
	words := strings.Fields(content)
	words = stemmer.StemMultiple(words)
	return strings.ToLower(strings.Join(words, " "))
}

func (tm *TermMatcher) ProcessContent(content string) string {
	content = tm.RemoveHyphens(content)
	content = tm.RemoveStopwords(content)
	content = tm.ProcessAndStem(content)
	return content
}

func (tm *TermMatcher) GetMatchingTerms(href string, anchorText string, searchTerms []string) []string {
	content := utils.ExtractLastSegmentFromURL(href)
	processedContent := tm.ProcessContent(content)
	tm.logger.Debug(fmt.Sprintf("Processed content from URL: %v", processedContent))

	anchorContent := tm.ProcessContent(anchorText)
	tm.logger.Debug(fmt.Sprintf("Processed anchor text: %v", anchorContent))

	combinedContent := tm.CombineContents(processedContent, anchorContent)
	tm.logger.Debug(fmt.Sprintf("Combined content: %v", combinedContent))

	if len(combinedContent) < minTitleLength {
		tm.logger.Debug(fmt.Sprintf("Combined content is less than minimum title length: %d", minTitleLength))
		return []string{}
	}

	matchingTerms := tm.FindMatchingTerms(combinedContent, searchTerms)

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

func (tm *TermMatcher) CombineContents(content1 string, content2 string) string {
	if content2 == "" {
		return content1
	}
	return content1 + " " + content2
}

func (tm *TermMatcher) ConvertToLowercase(content string) string {
	return strings.ToLower(content)
}

func (tm *TermMatcher) StemContent(content string) string {
	words := strings.Fields(content)
	stemmedWords := stemmer.StemMultiple(words)
	lowercaseStemmedWords := make([]string, len(stemmedWords))
	for i, word := range stemmedWords {
		lowercaseStemmedWords[i] = strings.ToLower(word)
	}
	return strings.Join(lowercaseStemmedWords, " ")
}

func (tm *TermMatcher) CompareTerms(searchTerm string, content string) float64 {
	searchTerm = strings.ToLower(searchTerm)
	similarity := strutil.Similarity(searchTerm, content, tm.swg)

	tm.logger.Debug(fmt.Sprintf("Compared terms: searchTerm=%s, content=%s, similarity=%.2f", searchTerm, content, similarity))

	return similarity
}

func (tm *TermMatcher) CompareAndAppendTerm(searchTerm string, content string, matchingTerms *[]string) {
	similarity := tm.CompareTerms(searchTerm, content)
	tm.logger.Debug(fmt.Sprintf("Compared terms: searchTerm=%s, similarity=%.2f", searchTerm, similarity))
	if similarity >= 0.9 { // Increase the threshold to 0.9
		*matchingTerms = append(*matchingTerms, searchTerm)
		tm.logger.Debug(fmt.Sprintf("Matching term found: %v", searchTerm))
	}
}

func (tm *TermMatcher) FindMatchingTerms(content string, searchTerms []string) []string {
	var matchingTerms []string

	content = tm.ConvertToLowercase(content)
	contentStemmed := tm.StemContent(content)

	tm.logger.Debug(fmt.Sprintf("Stemmed content: %v", contentStemmed))

	for _, searchTerm := range searchTerms {
		searchTerm = tm.ConvertToLowercase(searchTerm)
		searchTermStemmed := tm.StemContent(searchTerm)

		words := strings.Fields(searchTermStemmed)
		for _, word := range words {
			tm.CompareAndAppendTerm(word, contentStemmed, &matchingTerms)
		}
	}

	if len(matchingTerms) == 0 {
		tm.logger.Debug("No matching terms found")
		return []string{}
	}

	tm.logger.Debug(fmt.Sprintf("Matching terms result: %v", matchingTerms))
	return matchingTerms
}
