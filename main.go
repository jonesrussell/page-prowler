package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
//	"net"
    "net/url"
	"os"
//	"strconv"

	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Fact struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type Webpage struct {
	Head	string	`json:"head"`
	Body	string	`json:"body"`
}

var collection *mongo.Collection
var ctx = context.TODO()

func init() {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
	str := os.Args[1]
	fmt.Printf("URL : %s\n", str)

	u, err := url.Parse(str)
    if err != nil {
        panic(err)
    }

	fmt.Println(u.Host)
    // host, port, _ := net.SplitHostPort(u.Host)
    // fmt.Println(host)
	
	preview(str)

	allFacts := make([]Fact, 0)

	collector := colly.NewCollector(
		// colly.AllowedDomains(host),
		colly.AllowedDomains(u.Host),
	)

	/*collector.OnHTML("article", func(element *colly.HTMLElement) {
		factId, err := strconv.Atoi(element.Attr("id"))
		if err != nil {
			log.Println("Could not get id")
		}
		factDesc := element.Text

		fact := Fact{
			ID:          factId,
			Description: factDesc,
		}

		allFacts = append(allFacts, fact)
	})*/

	collector.OnHTML("html", func(element *colly.HTMLElement) {
		fmt.Println(element)
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	collector.Visit(str)

	writeJSON(allFacts)
}

func writeJSON(data []Fact) {
	file, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	_ = ioutil.WriteFile("rhinofacts.json", file, 0644)

}
