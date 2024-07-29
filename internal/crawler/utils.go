package crawler

import (
	"errors"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/page-prowler/models"
)

func (cm *CrawlManager) getHref(e *colly.HTMLElement) (string, error) {
	href := e.Request.AbsoluteURL(e.Attr("href"))
	if href == "" {
		return "", errors.New("Found anchor element with no href attribute")
	}
	return href, nil
}

func (cm *CrawlManager) createPageData(href string) models.PageData {
	return models.PageData{
		URL: href,
	}
}

func (cm *CrawlManager) handleMatchingTerms(options *CrawlOptions, currentURL string, pageData models.PageData, matchingTerms []string) error {
	err := cm.ProcessMatchingLink(currentURL, pageData, matchingTerms)
	if err != nil {
		return err
	}
	cm.UpdateStats(options, matchingTerms)

	return nil
}

func (cm *CrawlManager) processLink(e *colly.HTMLElement, href string) error {
	cm.StatsManager.LinkStats.IncrementTotalLinks()
	pageData := cm.createPageData(href)
	matchingTerms := cm.TermMatcher.GetMatchingTerms(href, e.Text, cm.Options.SearchTerms)
	err := cm.handleMatchingTerms(cm.Options, e.Request.URL.String(), pageData, matchingTerms)
	if err != nil {
		return err
	}

	return nil
}

func (cm *CrawlManager) ProcessMatchingLink(href string, pageData models.PageData, matchingTerms []string) error {
	if href == "" {
		return errors.New("Missing URL for matching link")
	}

	pageData.UpdatePageData(href, matchingTerms)
	cm.AppendResult(pageData)
	return nil
}

func (cm *CrawlManager) UpdateStats(_ *CrawlOptions, matchingTerms []string) {
	if len(matchingTerms) > 0 {
		cm.StatsManager.LinkStats.IncrementMatchedLinks()
	} else {
		cm.StatsManager.LinkStats.IncrementNotMatchedLinks()
	}
}
