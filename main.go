package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/segmentio/kafka-go"
)

type Article struct {
	Href  string
	Title string
	Image string
}

const (
	kafka_uri = "kafka:9092"
	topic     = "streetcode"
	partition = 0
)

func main() {
	webpage := os.Args[1]

	u, err := url.Parse(webpage)
	if err != nil {
		panic(err)
	}

	articles := make([]Article, 0)

	collector := colly.NewCollector(
		colly.AllowedDomains(u.Host),
		colly.Async(true),
		colly.URLFilters(
			regexp.MustCompile(`https://www.sudbury\.com(|/local-news.+)$`),
			regexp.MustCompile(`https://www.sudbury\.com/membership.+`),
			regexp.MustCompile(`https://www.sudbury\.com/weather.+`),
		),
		// colly.Debugger(&debug.LogDebugger{}),
	)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*sudbury.com" glob
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*sudbury.com",
		Parallelism: 2,
		Delay:       500 * time.Millisecond,
	})

	collector.OnHTML("a.section-item", func(element *colly.HTMLElement) {
		href := element.Attr("href")
		title := element.Text
		src := ""

		element.DOM.Find("img").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the src
			dataSrc, found := s.Attr("data-src")
			src = ""
			if found {
				src = dataSrc
			}
		})

		article := Article{
			Title: title,
			Href:  href,
			Image: src,
		}

		articles = append(articles, article)
	})

	collector.OnHTML(`a[href]`, func(e *colly.HTMLElement) {
		foundURL := e.Request.AbsoluteURL(e.Attr("href"))
		// foundURL := e.Attr("href")
		collector.Visit(foundURL)
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	collector.OnScraped(func(r *colly.Response) {
		// writeTopic(articles)
	})

	fmt.Println("Webpage:", webpage)
	collector.Visit(webpage)
	collector.Wait()
}

func writeTopic(articles []Article) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafka_uri, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	messages := make([]kafka.Message, 0)
	for i := 0; i < len(articles); i++ {
		// Convert article struct to json
		article, err := json.Marshal(articles[i])
		if err != nil {
			panic(err)
		}

		messages = append(messages, kafka.Message{
			Key:   []byte(articles[i].Href),
			Value: article,
		})
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	fmt.Println("Write")
	_, err = conn.WriteMessages(messages...)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
