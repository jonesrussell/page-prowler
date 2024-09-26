package matcher

import "strings"

// Matcher defines the interface for matching content.
type Matcher interface {
	// Match checks if the given content matches certain criteria.
	Match(content string, pattern string) (bool, error) // Accept pattern as an argument
}

// BaseMatcherInterface defines methods for processing content.
type BaseMatcherInterface interface {
	// ProcessContent processes the input content and returns a modified version.
	ProcessContent(content string) string
}

type DefaultMatcher struct{}

func (dm *DefaultMatcher) Match(content string, pattern string) (bool, error) { // Updated signature
	if content == "" {
		return false, nil // Return false if content is empty, no error
	}
	return strings.Contains(content, pattern), nil // Check for the provided pattern in the content
}
