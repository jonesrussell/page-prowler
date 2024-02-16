package crawler

import (
	"errors"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/termmatcher"
)

func (cs *CrawlManager) GetAnchorElementHandler(options *CrawlOptions) func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		href := cs.getHref(e)
		if href == "" {
			return
		}

		cs.processLink(e, href, options)
		err := cs.visitWithColly(href)
		if err != nil {
			cs.LoggerField.Debug("[GetAnchorElementHandler] Error visiting URL", "url", href, "error", err)
		}
	}
}

func (cs *CrawlManager) getHref(e *colly.HTMLElement) string {
	href := e.Request.AbsoluteURL(e.Attr("href"))
	if href == "" {
		cs.LoggerField.Debug("Found anchor element with no href attribute")
	} else {
		cs.LoggerField.Debug("Processing link", "href", href)
	}
	return href
}

func (cs *CrawlManager) incrementTotalLinks(options *CrawlOptions) {
	cs.StatsManager.LinkStats.IncrementTotalLinks()
	cs.LoggerField.Debug("Incremented total links count")
}

func (cs *CrawlManager) logCurrentURL(e *colly.HTMLElement) {
	cs.LoggerField.Debug("Current URL being crawled", "url", e.Request.URL.String())
}

func (cs *CrawlManager) createPageData(href string) PageData {
	return PageData{
		URL: href,
	}
}

func (cs *CrawlManager) logSearchTerms(options *CrawlOptions) {
	cs.LoggerField.Debug("Search terms", "terms", options.SearchTerms)
}

func (cs *CrawlManager) getMatchingTerms(href string, anchorText string, options *CrawlOptions) []string {
	return termmatcher.GetMatchingTerms(href, anchorText, options.SearchTerms, cs.Logger())
}

func (cs *CrawlManager) handleMatchingTerms(options *CrawlOptions, currentURL string, pageData PageData, matchingTerms []string) {
	if len(matchingTerms) > 0 {
		cs.ProcessMatchingLinkAndUpdateStats(options, currentURL, pageData, matchingTerms)
	} else {
		cs.incrementNonMatchedLinkCount(options)
		cs.LoggerField.Debug("Link does not match search terms", "link", pageData.URL)
	}
}

func (cs *CrawlManager) processLink(e *colly.HTMLElement, href string, options *CrawlOptions) {
	cs.incrementTotalLinks(options)
	cs.logCurrentURL(e)
	pageData := cs.createPageData(href)
	cs.logSearchTerms(options)
	matchingTerms := cs.getMatchingTerms(href, e.Text, options)
	cs.handleMatchingTerms(options, e.Request.URL.String(), pageData, matchingTerms)
}

func (cs *CrawlManager) handleSetupError(err error) error {
	cs.LoggerField.Error("Error setting up crawling logic", "error", err)
	return err
}

func (cs *CrawlManager) ProcessMatchingLinkAndUpdateStats(options *CrawlOptions, href string, pageData PageData, matchingTerms []string) {
	if cs == nil {
		cs.LoggerField.Error("CrawlManager instance is nil")
		return
	}

	if href == "" {
		cs.LoggerField.Error("Missing URL for matching link")
		return
	}

	cs.incrementMatchedLinks()
	cs.LoggerField.Debug("Incremented matched links count")

	pageData.UpdatePageData(href, matchingTerms)
	cs.AppendResult(options, pageData)
}

func (cs *CrawlManager) incrementMatchedLinks() {
	cs.StatsManager.LinkStats.IncrementMatchedLinks()
}

func (cs *CrawlManager) incrementNonMatchedLinkCount(options *CrawlOptions) {
	cs.StatsManager.LinkStats.IncrementNotMatchedLinks()
	cs.LoggerField.Debug("Incremented not matched links count")
}

func (cs *CrawlManager) createLimitRule() *colly.LimitRule {
	return &colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: DefaultParallelism,
		Delay:       DefaultDelay,
	}
}

func (cs *CrawlManager) splitSearchTerms(searchTerms string) []string {
	terms := strings.Split(searchTerms, ",")
	var validTerms []string
	for _, term := range terms {
		if term != "" {
			validTerms = append(validTerms, term)
		}
	}
	return validTerms
}

func (cs *CrawlManager) createStartCrawlingOptions(crawlSiteID string, searchTerms []string, debug bool) *CrawlOptions {
	var results []PageData
	return NewCrawlOptions(crawlSiteID, searchTerms, debug, &results)
}

// GetHostFromURL extracts the host from the given URL.
func GetHostFromURL(inputURL string, appLogger logger.Logger) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		appLogger.Error("Failed to parse URL", "url", inputURL, "error", err)
		return "", err
	}

	host := parsedURL.Hostname()
	if host == "" {
		appLogger.Errorf("failed to extract host from URL")
		return "", errors.New("failed to extract host from URL")
	}

	return host, nil
}
