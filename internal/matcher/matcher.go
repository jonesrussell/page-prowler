package matcher

type Matcher interface {
	Match(content string) bool
}
