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
	"github.com/go-redis/redis/v8"
	"github.com/gocolly/colly"
	"github.com/golang-module/carbon/v2"
	"github.com/segmentio/kafka-go"
)

type Article struct {
	Href  string `json:"href"`
	Title string `json:"title"`
	Image string `json:"image"`
	Body  string `json:"body"`
}

const (
	redis_uri  = "localhost"
	redis_port = "6379"
)

func main() {
	log.Println("Publisher started")
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redis_uri, redis_port),
	})

	_, err := redisClient.Ping(redisClient.Context()).Result()
	if err != nil {
		log.Fatal("Unable to connect to Redis", err)
	}

	log.Println("Connected to Redis server")

	os.Exit(0)

	// Retrieve URL parameter
	webpage := os.Args[1]
	u, err := url.Parse(webpage)
	if err != nil {
		panic(err)
	}
	// articles := make([]Article, 0)

	collector := colly.NewCollector(
		colly.AllowedDomains(u.Host),
		colly.Async(true),
		colly.URLFilters(
			// regexp.MustCompile(`https://www.sudbury\.com(|/local-news.+)$`),
			// regexp.MustCompile(`https://www.sudbury\.com/membership.+`),
			regexp.MustCompile(`https://www.sudbury\.com/police.+`),
			// regexp.MustCompile(`https://www.sudbury\.com/weather.+`),
		),
		// colly.Debugger(&debug.LogDebugger{}),
	)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*sudbury.com" glob
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*sudbury.com",
		Parallelism: 2,
		Delay:       1000 * time.Millisecond,
	})

	/*collector.OnHTML("a.section-item", func(element *colly.HTMLElement) {
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
	})*/

	collector.OnHTML(`a[href]`, func(e *colly.HTMLElement) {
		foundURL := e.Request.AbsoluteURL(e.Attr("href"))
		collector.Visit(foundURL)
	})

	// Article page
	collector.OnHTML(`section .details`, func(e *colly.HTMLElement) {
		datetime, found := e.DOM.FindMatcher(goquery.Single("time")).Attr("datetime")
		if !found {
			panic(err)
		}
		today := carbon.Now()
		articleDate := carbon.Parse(datetime)
		diff := articleDate.DiffInMonths(today)

		title := e.DOM.FindMatcher(goquery.Single("h1")).Text()

		href, found := e.DOM.FindMatcher(goquery.Single(".details-share .nav")).Attr("data-url")
		if !found {
			panic(err)
		}

		src, found := e.DOM.FindMatcher(goquery.Single("img")).Attr("src")
		if !found {
			panic(err)
		}

		body, err := e.DOM.FindMatcher(goquery.Single("#details-body")).Html()
		if err != nil {
			panic(err)
		}

		article := Article{
			Title: title,
			Href:  href,
			Image: src,
			Body:  body,
		}

		if diff <= 1 {
			// fmt.Println("Date:", articleDate)
			writeTopic(article)
		}

	})

	collector.OnRequest(func(request *colly.Request) {
		// fmt.Println("Visiting", request.URL.String())
	})

	collector.OnScraped(func(r *colly.Response) {
		// writeTopic(articles)
	})

	fmt.Println("Webpage:", webpage)
	collector.Visit(webpage)
	collector.Wait()
}

func writeTopic(article Article) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafka_uri, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	messages := make([]kafka.Message, 0)

	// Convert article struct to json
	articleJson, err := json.Marshal(article)
	if err != nil {
		panic(err)
	}

	messages = append(messages, kafka.Message{
		Key:   []byte(article.Href),
		Value: articleJson,
	})

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
