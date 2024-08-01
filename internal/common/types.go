package common

import "github.com/jonesrussell/page-prowler/crawler"

type CrawlManagerKeyType = crawler.CrawlManagerInterface

var crawlManagerKey crawler.CrawlManagerInterface

// CrawlManagerKey is the key for storing and retrieving the CrawlManagerInterface from the context.
var CrawlManagerKey = &crawlManagerKey

const CrawlManagerKeyStr = "crawlManagerKey"
