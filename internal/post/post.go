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

const Template = "api/post.json"

var (
	username = ""
	password = ""
)

func SetUsername(user string) {
	username = user
}

func SetPassword(pass string) {
	password = pass
}

func Process(msg myredis.MsgPost, url string) error {
	log.Println("checking href", msg.Href)

	// Assemble Streetcode API url that will search for link
	urlTest := fmt.Sprintf("%s%s", url, msg.Href)

	resData, err := checkHref(urlTest)
	if err != nil {
		return err
	}

	// Process the JSON response data
	data := Response{}
	json.Unmarshal([]byte(resData), &data)

	// Finally, no data means we can publish to Streetcode
	if len(data.Data) == 0 {
		response := create(msg, url)
		defer response.Body.Close()

		log.Printf("INFO: [response] %s\n", response.Status)
	} else {
		log.Printf("INFO: [exists] %s", msg.Href)
	}

	return nil
}

func prepare(msg myredis.MsgPost) ([]byte, error) {
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

	postData.Data.Attributes.FieldPost.Value = msg.Href
	postData.Data.Relationships.FieldRecipientGroup.Data.Id = msg.Group

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

func create(msg myredis.MsgPost, url string) *http.Response {
	jsonData, err := prepare(msg)
	if err != nil {
		log.Fatalln(err)
	}

	// POST to Streetcode
	request, _ := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)
	request.Header.Set("Content-Type", "application/vnd.api+json")
	request.Header.Set("Accept", "application/vnd.api+json")
	request.Header.Add("Authorization", "Basic "+basicAuth(username, password))

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
