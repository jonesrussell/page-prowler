package post

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/google/jsonapi"
	"github.com/joho/godotenv"
)

type Response struct {
	Data []Entity `jsonapi:"attr,data"`
}

type Entity struct {
	Id   string `jsonapi:"id"`
	Type string `jsonapi:"attr,type"`
}

type PostPhoto struct {
	FieldPost           FieldPost           `jsonapi:"attr,field_post"`
	FieldVisibility     int                 `jsonapi:"attr,field_visibility"`
	FieldRecipientGroup FieldRecipientGroup `jsonapi:"relation,field_recipient_group"`
}

type FieldPost struct {
	Value  string `jsonapi:"attr,value"`
	Format string `jsonapi:"attr,format"`
}

type FieldRecipientGroup struct {
	ID   string `jsonapi:"primary,id"`
	Type string `jsonapi:"attr,type"`
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
		log.Println(err)
		return err
	}

	if !resData {
		log.Println("Yes")
		response := create(href)
		defer response.Body.Close()

		log.Printf("INFO: [response] %s\n", response.Status)
		foo, _ := ioutil.ReadAll(response.Body)
		log.Printf("INFO: [response] %s\n", foo)
	} else {
		log.Println("No")
		log.Printf("INFO: [exists] %s", href)
	}

	return nil
}

func prepare(href string) []byte {
	// Open our JSON template
	jsonFile, err := os.Open(Template)
	if err != nil {
		log.Fatalln(err)
	}

	defer jsonFile.Close()

	w := httptest.NewRecorder()
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)

	// we initialize our PostPhoto array
	postData := new(PostPhoto)

	if err := jsonapi.UnmarshalPayload(jsonFile, postData); err != nil {
		log.Println("prepare() UnmarshalPayload")
		http.Error(w, err.Error(), 500)
	}

	postData.FieldPost.Value = href
	postData.FieldRecipientGroup.ID = GroupSudbury

	if err := jsonapi.MarshalPayload(w, postData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return w.Body.Bytes()
}

func checkHref(href string) (bool, error) {
	c := http.Client{Timeout: time.Duration(3) * time.Second}

	req, err := http.NewRequest("GET", href, nil)
	if err != nil {
		log.Printf("error %s", err)
		panic(err)
	}

	req.Header.Add("Accept", `application/json`)
	resp, err := c.Do(req)
	if err != nil {
		log.Printf("error %s", err)
		panic(err)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	data := Response{}
	json.Unmarshal([]byte(respBody), &data)

	if len(data.Data) == 0 {
		return false, nil
	}

	return true, nil
}

func create(href string) *http.Response {
	jsonData := prepare(href)

	// POST to Streetcode
	request, err := http.NewRequest(
		"POST",
		os.Getenv("API_URL"),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Printf("%s", err)
		// return
	}

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
