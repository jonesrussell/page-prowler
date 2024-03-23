package crawler

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/internal/logger"
	"github.com/jonesrussell/page-prowler/internal/termmatcher"
)

func (cm *CrawlManager) GetAnchorElementHandler(options *CrawlOptions) func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		href := cm.getHref(e)
		if href == "" {
			return
		}

		cm.processLink(e, href, options)
		err := cm.visitWithColly(href)
		if err != nil {
			cm.LoggerField.Debug(fmt.Sprintf("[GetAnchorElementHandler] Error visiting URL: %s, Error: %v", href, err))
		}
	}
}

func (cm *CrawlManager) getHref(e *colly.HTMLElement) string {
	href := e.Request.AbsoluteURL(e.Attr("href"))
	if href == "" {
		cm.LoggerField.Debug("Found anchor element with no href attribute")
	} else {
		cm.LoggerField.Debug(fmt.Sprintf("Processing link: %s", href))
	}
	return href
}

func (cm *CrawlManager) incrementTotalLinks() {
	cm.StatsManager.LinkStats.IncrementTotalLinks()
	cm.LoggerField.Debug("Incremented total links count")
}

func (cm *CrawlManager) logCurrentURL(e *colly.HTMLElement) {
	cm.LoggerField.Debug(fmt.Sprintf("Current URL being crawled: %s", e.Request.URL.String()))
}

func (cm *CrawlManager) createPageData(href string) PageData {
	return PageData{
		URL: href,
	}
}

func (cm *CrawlManager) logSearchTerms(options *CrawlOptions) {
	cm.LoggerField.Debug(fmt.Sprintf("Search terms: %v", options.SearchTerms))
}

func (cm *CrawlManager) getMatchingTerms(href string, anchorText string, options *CrawlOptions) []string {
	return termmatcher.GetMatchingTerms(href, anchorText, options.SearchTerms, cm.Logger())
}

func (cm *CrawlManager) handleMatchingTerms(options *CrawlOptions, currentURL string, pageData PageData, matchingTerms []string) {
	if len(matchingTerms) > 0 {
		cm.ProcessMatchingLinkAndUpdateStats(options, currentURL, pageData, matchingTerms)
	} else {
		cm.incrementNonMatchedLinkCount()
		cm.LoggerField.Debug(fmt.Sprintf("Link does not match search terms: %s", pageData.URL))
	}
}

func (cm *CrawlManager) processLink(e *colly.HTMLElement, href string, options *CrawlOptions) {
	cm.incrementTotalLinks()
	cm.logCurrentURL(e)
	pageData := cm.createPageData(href)
	cm.logSearchTerms(options)
	matchingTerms := cm.getMatchingTerms(href, e.Text, options)
	cm.handleMatchingTerms(options, e.Request.URL.String(), pageData, matchingTerms)
}

func (cm *CrawlManager) handleSetupError(err error) error {
	cm.LoggerField.Error(fmt.Sprintf("Error setting up crawling logic: %v", err))
	return err
}

func (cm *CrawlManager) ProcessMatchingLinkAndUpdateStats(options *CrawlOptions, href string, pageData PageData, matchingTerms []string) {
	if href == "" {
		cm.LoggerField.Error("Missing URL for matching link")
		return
	}

	cm.incrementMatchedLinks()
	cm.LoggerField.Debug("Incremented matched links count")

	pageData.UpdatePageData(href, matchingTerms)
	cm.AppendResult(options, pageData)
}

func (cm *CrawlManager) incrementMatchedLinks() {
	cm.StatsManager.LinkStats.IncrementMatchedLinks()
}

func (cm *CrawlManager) incrementNonMatchedLinkCount() {
	cm.StatsManager.LinkStats.IncrementNotMatchedLinks()
	cm.LoggerField.Debug("Incremented not matched links count")
}

func (cm *CrawlManager) createLimitRule() *colly.LimitRule {
	return &colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: DefaultParallelism,
		Delay:       DefaultDelay,
	}
}

func (cm *CrawlManager) splitSearchTerms(searchTerms string) []string {
	terms := strings.Split(searchTerms, ",")
	var validTerms []string
	for _, term := range terms {
		if term != "" {
			validTerms = append(validTerms, term)
		}
	}
	return validTerms
}

func (cm *CrawlManager) createStartCrawlingOptions(crawlSiteID string, searchTerms []string, debug bool) *CrawlOptions {
	var results []PageData
	return NewCrawlOptions(crawlSiteID, searchTerms, debug, &results)
}

// GetHostFromURL extracts the host from the given URL.
func GetHostFromURL(inputURL string, appLogger logger.Logger) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		appLogger.Error(fmt.Sprintf("Failed to parse URL: %s, Error: %v", inputURL, err))
		return "", err
	}

	host := parsedURL.Hostname()
	if host == "" {
		appLogger.Error("failed to extract host from URL")
		return "", errors.New("failed to extract host from URL")
	}

	return host, nil
}
