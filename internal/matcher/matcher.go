package matcher

type Matcher interface {
	Match(content string) bool
}

type BaseMatcherInterface interface {
	ProcessContent(content string) string
}
