package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/jonesrussell/crawler/internal/myredis"
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
		log.Fatalln(err)
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
		err := handleNewHref(href)
		if err != nil {
			log.Fatal(err)
		}

		myredis.AckEntry(id)
	}
}

func handleNewHref(href string) error {
	// Check if link has already been submitted to Streetcode
	// Assemble Streetcode API url that will search for link
	urlTest := fmt.Sprintf("%s%s", os.Getenv("API_FILTER_URL"), href)

	log.Println("checking url", urlTest)

	// Call Streetcode
	res, err := http.Get(urlTest)
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

	// Process the JSON response data
	data := Response{}
	json.Unmarshal([]byte(resData), &data)

	// Finally, no data means we can publish to Streetcode
	if len(data.Data) == 0 {
		// POST to Streetcode
		jsonData := fmt.Sprintf(`{"data":{"type":"post--photo","attributes":{"field_post":{"value":"%s","format":"basic_html"},"field_visibility": "1"},"relationships":{"field_recipient_group":{"data":{"type":"group--public_group","id":"b55fe232-0fbf-4fa8-b697-ff7bb863ae6a"}}}}}`, href)
		request, _ := http.NewRequest("POST", os.Getenv("API_URL"), bytes.NewBuffer([]byte(jsonData)))
		request.Header.Set("Content-Type", "application/vnd.api+json")
		request.Header.Set("Accept", "application/vnd.api+json")
		request.SetBasicAuth(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))

		client := &http.Client{}
		response, error := client.Do(request)
		if error != nil {
			panic(error)
		}
		defer response.Body.Close()

		fmt.Printf("INFO: [response] %s\n", response.Status)
	} else {
		log.Printf("INFO: [exists] %s", href)
	}
	return nil
}
