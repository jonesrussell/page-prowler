package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/jonesrussell/crawler/internal/myredis"
	"github.com/jonesrussell/crawler/internal/post"
)

type Response struct {
	Data []Entity `json:"data"`
}

type Entity struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

func main() {
	if godotenv.Load(".env") != nil {
		log.Fatal("error loading .env file")
	}

	// Connect to Redis
	redisClient := myredis.Connect()
	defer redisClient.Close()

	// Connect to Redis Stream
	err := myredis.Stream()
	if err != nil {
		log.Println(err)
	}

	log.Println("consumer started")

	for {
		entries, err := myredis.GetEntries()
		if err != nil {
			log.Fatal(err)
		}

		processEntries(entries)
	}
}

func processEntries(entries []redis.XStream) {
	messages := entries[0].Messages

	for i := 0; i < len(messages); i++ {
		processEntry(messages[i].Values, messages[i].ID)
	}
}

func processEntry(values map[string]interface{}, id string) {
	eventName := fmt.Sprintf("%v", values["eventName"])
	href := fmt.Sprintf("%v", values["href"])

	if eventName == "receivedUrl" {
		err := processHref(href)
		if err != nil {
			log.Fatal(err)
		}

		myredis.AckEntry(id)
	}
}

func processHref(href string) error {
	log.Println("checking href", href)

	// Assemble Streetcode API url that will search for link
	urlTest := fmt.Sprintf("%s%s", os.Getenv("API_FILTER_URL"), href)

	resData, err := checkHref(urlTest)
	if err != nil {
		log.Fatal(err)
	}

	// Process the JSON response data
	data := Response{}
	json.Unmarshal([]byte(resData), &data)

	// Finally, no data means we can publish to Streetcode
	if len(data.Data) == 0 {
		response := post.Create(href)
		defer response.Body.Close()

		log.Printf("INFO: [response] %s\n", response.Status)
	} else {
		log.Printf("INFO: [exists] %s", href)
	}

	return nil
}

func checkHref(href string) ([]byte, error) {
	// Call Streetcode
	res, err := http.Get(href)
	if err != nil {
		log.Fatal(err)
	}

	// Http call succeeded, check response code
	if res.StatusCode != 200 {
		b, _ := ioutil.ReadAll(res.Body)
		log.Fatal(string(b))
	}

	// Process good response from Streetcode
	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(" response.Body ", err)
	}

	return resData, err
}
