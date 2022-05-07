package post

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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
	Value string `json:"value"`
}

type FieldRecipientGroup struct {
	Data RData `json:"data"`
}

type RData struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

const Template = "api/post.json"

var (
	GroupSudbury       = ""
	GroupEspanola      = ""
	GroupElliotLake    = ""
	GroupNorthBay      = ""
	GroupSturgeonFalls = ""
)

func init() {
	if godotenv.Load(".env") != nil {
		log.Fatal("error loading .env file")
	}

	GroupSudbury = os.Getenv("GROUP_SUDBURY")
	GroupEspanola = os.Getenv("GROUP_ESPANOLA")
	GroupElliotLake = os.Getenv("GROUP_ELLIOTLAKE")
	GroupNorthBay = os.Getenv("GROUP_NORTHBAY")
	GroupSturgeonFalls = os.Getenv("GROUP_STURGEONFALLS")
}

func ProcessHref(href string) error {
	log.Println("checking href", href)

	// Assemble Streetcode API url that will search for link
	urlTest := fmt.Sprintf("%s%s", os.Getenv("API_FILTER_URL"), href)

	resData, err := checkHref(urlTest)
	if err != nil {
		return err
	}

	// Process the JSON response data
	data := Response{}
	json.Unmarshal([]byte(resData), &data)

	// Finally, no data means we can publish to Streetcode
	if len(data.Data) == 0 {
		response := create(href)
		defer response.Body.Close()

		log.Printf("INFO: [response] %s\n", response.Status)
	} else {
		log.Printf("INFO: [exists] %s", href)
	}

	return nil
}

func prepare(href string) ([]byte, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(Template)
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

	return json.Marshal(postData)
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

func create(href string) *http.Response {
	jsonData, err := prepare(href)
	if err != nil {
		log.Fatalln(err)
	}

	// POST to Streetcode
	request, _ := http.NewRequest(
		"POST",
		os.Getenv("API_URL"),
		bytes.NewBuffer(jsonData),
	)
	request.Header.Set("Content-Type", "application/vnd.api+json")
	request.Header.Set("Accept", "application/vnd.api+json")
	request.Header.Add("Authorization", "Basic "+basicAuth(os.Getenv("USERNAME"), os.Getenv("PASSWORD")))

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}

	return response
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
