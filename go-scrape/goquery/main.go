package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

//ExampleScrape ...
func ExampleScrape() {
	doc, err := goquery.NewDocument("http://metalsucks.net")
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
		band := s.Find("a").Text()
		title := s.Find("i").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})
}

func main() {
	ExampleScrape()
}
