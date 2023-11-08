// Package post provides the functionality to process and create posts
// based on web crawl results using a hypothetical Streetcode API.
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

	"github.com/jonesrussell/crawler/internal/rediswrapper"
)

// Response represents the structure of an API response.
type Response struct {
	Data []Entity `json:"data"` // List of entities returned from the API
}

// Entity represents a data entity within the API response.
type Entity struct {
	Type string `json:"type"` // Type of the entity
	ID   string `json:"id"`   // ID of the entity
}

// Photo represents the structure of a photo post request.
type Photo struct {
	Data Data `json:"data"` // Data of the post
}

// Data holds the structure of the post data.
type Data struct {
	Type          string        `json:"type"`          // Type of the post
	Attributes    Attributes    `json:"attributes"`    // Attributes of the post
	Relationships Relationships `json:"relationships"` // Relationships associated with the post
}

// Attributes contains the attributes of a post.
type Attributes struct {
	FieldPost       FieldPost `json:"field_post"`       // The main post content
	FieldVisibility string    `json:"field_visibility"` // Visibility setting of the post
}

// Relationships contains the relationships of a post.
type Relationships struct {
	FieldRecipientGroup FieldRecipientGroup `json:"field_recipient_group"` // Group recipient of the post
}

// FieldPost defines the content and format of a post.
type FieldPost struct {
	Value  string `json:"value"`  // Content of the post
	Format string `json:"format"` // Format of the post content
}

// FieldRecipientGroup defines the recipient group of a post.
type FieldRecipientGroup struct {
	Data RData `json:"data"` // Data about the recipient group
}

// RData holds the type and ID for a relationship data.
type RData struct {
	Type string `json:"type"` // Type of the relationship
	ID   string `json:"id"`   // ID of the relationship
}

// Template is the path to the JSON template file for post creation.
const Template = "api/post.json"

var (
	username = "" // Username for API authentication
	password = "" // Password for API authentication
)

// SetUsername sets the username for API authentication.
func SetUsername(user string) {
	username = user
}

// SetPassword sets the password for API authentication.
func SetPassword(pass string) {
	password = pass
}

// Process handles the processing of a MsgPost message and creates a post if needed.
func Process(msg rediswrapper.MsgPost, url string) error {
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

	// If no data is found, create a new post on Streetcode
	if len(data.Data) == 0 {
		response := create(msg, url)
		defer response.Body.Close()

		log.Printf("INFO: [response] %s\n", response.Status)
	} else {
		log.Printf("INFO: [exists] %s", msg.Href)
	}

	return nil
}

// prepare constructs the JSON payload for a new post based on MsgPost.
func prepare(msg rediswrapper.MsgPost) ([]byte, error) {
	// Open the JSON template file
	jsonFile, err := os.Open(Template)
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()

	// Read the file's content as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Initialize the Photo structure
	postData := Photo{}

	// Unmarshal the byte array into postData
	json.Unmarshal(byteValue, &postData)

	// Populate the postData fields with the message details
	postData.Data.Attributes.FieldPost.Value = msg.Href
	postData.Data.Relationships.FieldRecipientGroup.Data.ID = msg.Group

	// Marshal the postData back into JSON
	return json.Marshal(postData)
}

// checkHref checks if the href exists on the Streetcode API.
func checkHref(href string) ([]byte, error) {
	// Make an HTTP GET request to Streetcode
	res, err := http.Get(href)
	if err != nil {
		log.Fatal(err)
	}

	// Check the response code
	if res.StatusCode != 200 {
		b, _ := ioutil.ReadAll(res.Body)
		log.Fatal(string(b))
	}

	// Read the response body
	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(" response.Body ", err)
	}

	return resData, err
}

// create sends a POST request to create a new post on the Streetcode API.
func create(msg rediswrapper.MsgPost, url string) *http.Response {
	jsonData, err := prepare(msg)
	if err != nil {
		log.Fatalln(err)
	}

	// Prepare the POST request
	request, _ := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)
	request.Header.Set("Content-Type", "application/vnd.api+json")
	request.Header.Set("Accept", "application/vnd.api+json")
	request.Header.Add("Authorization", "Basic "+basicAuth(username, password))

	// Make the POST request
	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}

	return response
}

// basicAuth encodes the username and password for HTTP basic authentication.
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
