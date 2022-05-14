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
	// Retrieve URL to crawl from arguments
	if len(os.Args) < 3 {
		log.Println("usage: ./crawler https://www.sudbury.com c45fe232-0fbd-4fj8-b097-ff7bb863ae6b")
		os.Exit(0)
	}
	crawlUrl := os.Args[1]
	group := os.Args[2]

	// Load the environment variables
	if godotenv.Load(".env") != nil {
		log.Println("error loading .env file")
	}

	// Setup the Redis connection
	redisClient := myredis.Connect()
	defer redisClient.Close()

	// Create a new crawler
	collector := colly.NewCollector(
		colly.Async(true),
		//colly.Debugger(&debug.LogDebugger{}),
	)

	// Set reasonable limits for responsible crawling
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       3000 * time.Millisecond,
	})

	// When a url is found
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// Extract the full url
		href := e.Request.AbsoluteURL(e.Attr("href"))

		// Determine if we will submit link to Redis
		/*matchedNewsMnm, _ := regexp.MatchString(`^https://www.midnorthmonitor.com/news/`, href)
		matchedPoliceSc, _ := regexp.MatchString(`^https://www.sudbury.com/police/`, href)
		matchedNewsEls, _ := regexp.MatchString(`^https://www.elliotlakestandard.ca/category/news/`, href)

		if matchedPoliceSc || matchedNewsMnm || matchedNewsEls {*/
		// Determine if url is drug related
		if drug.Related(href) {
			// Announce the drug related url
			fmt.Println(href)

			// Add url to publishing queue
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

	// When a url has finished being crawled
	collector.OnScraped(func(r *colly.Response) {
		// Retrieve the urls to be published
		hrefs, err := myredis.SMembers()
		if err != nil {
			log.Fatal(err)
		}

		// Loop over urls
		for i := range hrefs {
			href := hrefs[i]

			// Send url to Redis stream
			err = myredis.PublishHref(href, group)
			if err != nil {
				log.Fatal(err)
			}

			_, err = myredis.Del()
			if err != nil {
				log.Fatal(err)
			}

			if os.Getenv("CSV_WRITE") == "true" {
				mycsv.WriteHrefCsv(href)
			}
		}
	})

	// When crawling a url
	collector.OnRequest(func(r *colly.Request) {
		// Announce the url
		fmt.Println("Visiting", r.URL)
	})

	// Set error handler
	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Everything is setup, time to crawl
	log.Println("Crawler started...")
	collector.Visit(crawlUrl)
	collector.Wait()
}
