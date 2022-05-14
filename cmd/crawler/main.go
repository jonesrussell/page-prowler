package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/jonesrussell/crawler/internal/drug"
	"github.com/jonesrussell/crawler/internal/mycsv"
	"github.com/jonesrussell/crawler/internal/myredis"
)

func main() {
	// Retrieve URL to crawl
	if len(os.Args) < 2 {
		log.Println("usage: ./crawler https://www.sudbury.com")
		os.Exit(0)
	}
	crawlUrl := os.Args[1]

	if godotenv.Load(".env") != nil {
		log.Println("error loading .env file")
	}

	redisClient := myredis.Connect()
	defer redisClient.Close()

	log.Println("crawler started")

	collector := colly.NewCollector(
		colly.Async(true),
		// colly.Debugger(&debug.LogDebugger{}),
	)

	// Limit the number of threads started by colly to two
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       3000 * time.Millisecond,
	})

	// Act on every link; <a href="foo">
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Request.AbsoluteURL(e.Attr("href"))

		// Determine if we will submit link to Redis
		/*matchedNewsMnm, _ := regexp.MatchString(`^https://www.midnorthmonitor.com/news/`, href)
		matchedPoliceSc, _ := regexp.MatchString(`^https://www.sudbury.com/police/`, href)
		matchedNewsEls, _ := regexp.MatchString(`^https://www.elliotlakestandard.ca/category/news/`, href)

		if matchedPoliceSc || matchedNewsMnm || matchedNewsEls {*/
		if drug.Related(href) {
			fmt.Println(href)

			_, err := myredis.SAdd(href)
			if err != nil {
				log.Fatal(err)
			}
		}
		//}

		if os.Getenv("CRAWL_MODE") != "single" {
			collector.Visit(href)
		}
	})

	collector.OnScraped(func(r *colly.Response) {
		href, err := myredis.SPop()
		if err != nil {
			log.Fatal(err)
		}

		err = myredis.PublishHref(href)
		if err != nil {
			log.Fatal(err)
		}

		if os.Getenv("CSV_WRITE") == "true" {
			mycsv.WriteHrefCsv(href)
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
