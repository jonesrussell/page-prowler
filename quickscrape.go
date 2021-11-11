package main

import (
	"fmt"

	"github.com/badoux/goscraper"
)

func preview(url string) {
	s, err := goscraper.Scrape(url, 5)
	
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Icon : %s\n", s.Preview.Icon)
	fmt.Printf("Name : %s\n", s.Preview.Name)
	fmt.Printf("Title : %s\n", s.Preview.Title)
	fmt.Printf("Description : %s\n", s.Preview.Description)
	fmt.Printf("Image: %s\n", s.Preview.Images[0])
	fmt.Printf("Url : %s\n", s.Preview.Link)
}
