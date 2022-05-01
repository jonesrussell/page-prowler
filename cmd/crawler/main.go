package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/jonesrussell/crawler/internal/drug"
	"github.com/jonesrussell/crawler/internal/myredis"
)

func main() {
	// Retrieve URL to crawl
	if len(os.Args) < 2 {
		log.Fatalln("usage: ./crawler https://www.sudbury.com")
	}
	crawlUrl := os.Args[1]

	if godotenv.Load(".env") != nil {
		log.Println("error loading .env file")
	}

	redisClient := myredis.Connect()

	log.Println("crawler started")

	collector := colly.NewCollector(
		colly.Async(true),
		colly.URLFilters(
			regexp.MustCompile("https://www.sudbury.com/police"),
			regexp.MustCompile("https://www.midnorthmonitor.com/category/news"),
			regexp.MustCompile("https://www.midnorthmonitor.com/news"),
			regexp.MustCompile("https://www.midnorthmonitor.com/category/news/local-news"),
		),
		// colly.Debugger(&debug.LogDebugger{}),
	)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*sudbury.com" glob
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       3000 * time.Millisecond,
	})

	// Act on every link; <a href="foo">
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		foundHref := e.Request.AbsoluteURL(e.Attr("href"))

		// Determine if we will submit link to Redis
		matchedNewsMnm, _ := regexp.MatchString(`^https://www.midnorthmonitor.com/news/`, foundHref)
		matchedPoliceSc, _ := regexp.MatchString(`^https://www.sudbury.com/police/`, foundHref)

		if matchedPoliceSc || matchedNewsMnm {
			if drug.Related(foundHref) {
				fmt.Println(foundHref)
				myredis.SAdd(redisClient, foundHref)
			}

			/*if err = myredis.PublishHref(redisClient, foundHref); err != nil {
				log.Fatal(err)
			}

			if os.Getenv("CSV_WRITE") == "true" {
				mycsv.WriteHrefCsv(foundHref)
			}*/
		}

		if os.Getenv("CRAWL_MODE") != "single" {
			collector.Visit(foundHref)
		}
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// Set error handler
	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	collector.Visit(crawlUrl)
	collector.Wait()
}
