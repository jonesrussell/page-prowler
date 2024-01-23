package crawler

//go:generate mockery --name=MatchedLinkProcessor
type MatchedLinkProcessor interface {
	IncrementMatchedLinks(options *CrawlOptions, statsManager *StatsManager)
	HandleMatchingLinks(href string) error
	UpdatePageData(pageData *PageData, href string, matchingTerms []string)
	AppendResult(options *CrawlOptions, pageData PageData)
}

type ConcreteMatchedLinkProcessor struct {
	CrawlManager *CrawlManager
}

func (p *ConcreteMatchedLinkProcessor) IncrementMatchedLinks(options *CrawlOptions, statsManager *StatsManager) {
	p.CrawlManager.incrementMatchedLinks(options, statsManager.LinkStats)
}

func (p *ConcreteMatchedLinkProcessor) HandleMatchingLinks(href string) error {
	return p.CrawlManager.handleMatchingLinks(href)
}

func (p *ConcreteMatchedLinkProcessor) UpdatePageData(pageData *PageData, href string, matchingTerms []string) {
	p.CrawlManager.updatePageData(pageData, href, matchingTerms)
}

func (p *ConcreteMatchedLinkProcessor) AppendResult(options *CrawlOptions, pageData PageData) {
	p.CrawlManager.appendResult(options, pageData)
}
