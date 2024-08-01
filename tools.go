//go:build tools
// +build tools

package main

//go:generate mockery --name=github.com/jonesrussell/page-prowler/crawler.CrawlManagerInterface
//go:generate mockery --name=github.com/jonesrussell/page-prowler/crawler.CollectorInterface
//go:generate oapi-codegen api.yaml > internal/api/api.gen.go
