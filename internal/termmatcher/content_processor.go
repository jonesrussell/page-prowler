package termmatcher

import (
	"strings"

	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
)

type ContentProcessor interface {
	Process(content string) string
	Stem(content string) string
}

type DefaultContentProcessor struct {
	processedCache map[string]string
}

func NewDefaultContentProcessor() *DefaultContentProcessor {
	return &DefaultContentProcessor{
		processedCache: make(map[string]string),
	}
}

func (cp *DefaultContentProcessor) Process(content string) string {
	if processed, ok := cp.processedCache[content]; ok {
		return processed
	}

	content = strings.ReplaceAll(content, "-", " ")
	content = stopwords.CleanString(content, "en", false) // Change to false to keep words like "over"
	words := strings.Fields(content)

	for i, word := range words {
		stemmed, err := snowball.Stem(word, "english", true)
		if err == nil {
			words[i] = stemmed
		}
		words[i] = strings.ToLower(words[i])
	}

	processed := strings.Join(words, " ")
	cp.processedCache[content] = processed
	return processed
}

func (cp *DefaultContentProcessor) Stem(content string) string {
	words := strings.Fields(content)
	for i, word := range words {
		stemmed, err := snowball.Stem(word, "english", true)
		if err == nil {
			words[i] = stemmed
		}
		words[i] = strings.ToLower(words[i])
	}
	return strings.Join(words, " ")
}
