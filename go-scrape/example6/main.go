package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Course struct {
	Title       string
	Description string
	Creator     string
	Level       string
	URL         string
	Language    string
	Commitment  string
	HowToPass   string
	Rating      string
}

func main() {
	c := colly.NewCollector()

	c.AllowedDomains = []string{"coursera.org", "www.coursera.org"}
	c.CacheDir = "./coursera_cache"
	// Create another collector to scrape course details
	detailCollector := c.Clone()
	courses := make([]Course, 0, 200)
	c.OnHTML("a.rc-DomainNavItem", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		//fmt.Println(link)
		e.Request.Visit(link)
	})
	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})
	c.OnHTML(`a[name]`, func(e *colly.HTMLElement) {
		courseURL := e.Request.AbsoluteURL(e.Attr("href"))
		//fmt.Println(courseURL)
		if strings.Index(courseURL, "coursera.org/learn") != -1 {
			detailCollector.Visit(courseURL)
		}
	})
	// Extract details of the course
	detailCollector.OnHTML(`div[id=rendered-content]`, func(e *colly.HTMLElement) {
		//log.Println("Course found", e.Request.URL)
		title := e.ChildText(".course-title")
		if title == "" {
			log.Println("No title found", e.Request.URL)
		}
		course := Course{
			Title:       title,
			URL:         e.Request.URL.String(),
			Description: e.ChildText("div.content"),
			Creator:     e.ChildText("div.creator-names > span"),
		}
		e.DOM.Find("table.basic-info-table tr").Each(func(_ int, s *goquery.Selection) {
			//fmt.Println(s.Find("td:first-child").Text())
			switch s.Find("td:first-child").Text() {
			case "Language":
				course.Language = s.Find("td:nth-child(2)").Text()
			case "Level":
				course.Level = s.Find("td:nth-child(2)").Text()
			case "Commitment":
				course.Commitment = s.Find("td:nth-child(2)").Text()
			case "How To Pass":
				course.HowToPass = s.Find("td:nth-child(2)").Text()
			case "User Ratings":
				course.Rating = s.Find("td:nth-child(2) div:nth-of-type(2)").Text()
			}
		})
		courses = append(courses, course)
		//fmt.Println(course)
	})
	c.Visit("https://coursera.org/browse")
	jsonData, err := json.MarshalIndent(courses, "", "  ")
	if err != nil {
		panic(err)
	}

	// Dump json to the standard output (can be redirected to a file)
	fmt.Println(string(jsonData))
}
