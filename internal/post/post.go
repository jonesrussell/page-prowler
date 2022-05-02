package post

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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

func Create(href string) *http.Response {
	jsonData, err := prepare(href)
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

	return response
}

func prepare(href string) ([]byte, error) {
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

	return json.Marshal(postData)
}
