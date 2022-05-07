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

	pass := string(resData[:])
	log.Println(pass)
	// os.Exit(0)

	// Process the JSON response data
	/*data := Response{}
	json.Unmarshal(resData, &data)*/

	// postPhoto := PostPhoto{}
	op := jsonapi.OnePayload{}
	log.Println(op)

	err = jsonapi.UnmarshalPayload(bytes.NewBuffer(resData), &op)
	if err != nil {
		log.Println("ProcessHref() UnmarshalPayload")
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	log.Println(op)
	log.Print("POST to Streetcode? ")
	// Finally, no data means we can publish to Streetcode
	if len(resData) == 0 {
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

func checkHref(href string) ([]byte, error) {
	// log.Println("checkHref()")
	c := http.Client{Timeout: time.Duration(3) * time.Second}

	req, err := http.NewRequest("GET", href, nil)
	if err != nil {
		log.Printf("error %s", err)
		return []byte{}, err
	}

	req.Header.Add("Accept", `application/json`)
	resp, err := c.Do(req)
	if err != nil {
		log.Printf("error %s", err)
		return []byte{}, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	pass := string(respBody[:])
	log.Println(pass)

	var decoded []interface{}
	err = json.Unmarshal(respBody, &decoded)
	if err != nil {
		log.Println(err)
	}

	log.Println(decoded)

	os.Exit(1)

	/*op := jsonapi.OnePayload{}
	err = jsonapi.UnmarshalPayload(bytes.NewBuffer(respBody), &op.Data)
	if err != nil {
		log.Print("checkHref() UnmarshalPayload ")
		log.Print(err)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		os.Exit(0)
	}*/

	return respBody, err

	// return ioutil.ReadAll(resp.Body)

	/*	if err != nil {
		fmt.Printf("error %s", err)
		return []byte{}, err
	}*/

	/*
		postPhoto := new(PostPhoto)

		if err := jsonapi.UnmarshalPayload(r.Body, postPhoto); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// ...save your blog...

		w.Header().Set("Content-Type", jsonapi.MediaType)
		w.WriteHeader(http.StatusCreated)

		if err := jsonapi.MarshalPayload(w, postPhoto); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// Call Streetcode
		res, err := http.Get(href)
		if err != nil {
			log.Fatal(err)
		}
	*/
	// Http call succeeded, check response code
	/*if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		log.Fatal(string(b))
	}

	return body, err*/
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
