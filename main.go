package main

import (
	"fmt"
	"log"
	"time"

	"net/url"
	"os"

	"context"

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
	)

	collector.OnHTML("a.section-item", func(element *colly.HTMLElement) {
		href := element.Attr("href")
		title := element.Text
		element.DOM.Find("img").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the src
			src, foo := s.Attr("src")
			fmt.Printf("Review %d: %s -- %t\n", i, src, foo)
		})

		article := Article{
			Title: title,
			Href:  href,
			// Image: image,
		}

		articles = append(articles, article)
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	collector.Visit(webpage)

	writeTopic(articles)
}

func writeTopic(data []Article) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafka_uri, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	messages := make([]kafka.Message, 0)
	for i := 0; i < len(data); i++ {
		messages = append(messages, kafka.Message{
			Key:   []byte(data[i].Href),
			Value: []byte(data[i].Title),
		})
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(messages...)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
