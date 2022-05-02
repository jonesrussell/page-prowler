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

type PostPhoto struct {
	Data Data `json:"data"`
}

type Data struct {
	Type          string        `json:"type"`
	Attributes    Attributes    `json:"attributes"`
	Relationships Relationships `json:"relationships"`
}

type Attributes struct {
	FieldPost       FieldPost `json:"field_post"`
	FieldVisibility string    `json:"field_visibility"`
}

type Relationships struct {
	FieldRecipientGroup FieldRecipientGroup `json:"field_recipient_group"`
}

type FieldPost struct {
	Value  string `json:"value"`
	Format string `json:"format"`
}

type FieldRecipientGroup struct {
	Data RData `json:"data"`
}

type RData struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

const (
	GroupSudbury       = "b55fe232-0fbf-4fa8-b697-ff7bb863ae6a"
	GroupEspanola      = "85c9a7c9-bb3a-42d2-b2a0-a4dead6e9d77"
	GroupElliotLake    = "01123a12-3837-4883-9d4a-6642ff690fae"
	GroupNorthBay      = "e1bb6e47-76be-4781-85b7-1c541a108da1"
	GroupSturgeonFalls = "b4462f3b-d305-43c6-bf9f-98ed121fcd74"
)

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
		// Open our jsonFile
		jsonFile, err := os.Open("api/post.json")
		if err != nil {
			log.Fatalln(err)
		}
		defer jsonFile.Close()

		// read our opened jsonFile as a byte array.
		byteValue, _ := ioutil.ReadAll(jsonFile)

		// we initialize our PostPhoto array
		postData := PostPhoto{}

		// we unmarshal our byteArray which contains our
		// jsonFile's content into 'postData' which we defined above
		json.Unmarshal(byteValue, &postData)

		postData.Data.Attributes.FieldPost.Value = href
		postData.Data.Relationships.FieldRecipientGroup.Data.Id = GroupSudbury

		jsonData, err := json.Marshal(postData)
		if err != nil {
			log.Fatalln(err)
		}

		// POST to Streetcode
		request, _ := http.NewRequest(
			"POST",
			os.Getenv("API_URL"),
			bytes.NewBuffer([]byte(jsonData)),
		)
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
