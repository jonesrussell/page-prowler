package common

import (
	"github.com/jonesrussell/page-prowler/internal/crawler"
	"github.com/jonesrussell/page-prowler/internal/worker"
)

type CrawlManagerKeyType = crawler.CrawlManagerInterface
type CrawlerWorkerKeyType = worker.CrawlerWorkerInterface

var crawlManagerKey crawler.CrawlManagerInterface
var crawlerWorkerKey worker.CrawlerWorkerInterface

var CrawlManagerKey = &crawlManagerKey
var CrawlerWorkerKey = &crawlerWorkerKey

const CrawlManagerKeyStr = "crawlManagerKey"
const CrawlerWorkerKeyStr = "crawlerWorkerKey"
