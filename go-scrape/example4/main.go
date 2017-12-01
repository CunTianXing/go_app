package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

const DATE_FORMAT = "Jan 02, 2006"

type Course struct {
	CourseID  string
	Run       string
	Name      string
	Number    string
	StartDate string
	EndDate   string
	URL       string
}

func main() {
	c := colly.NewCollector()
	//c.SetDebugger(&debug.LogDebugger{})
	c.AllowedDomains = []string{"indonesiax.co.id", "www.indonesiax.co.id"}
	c.CacheDir = "./cache"
	courses := make([]Course, 0, 200)
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		//fmt.Println(link)
		if !strings.HasPrefix(link, "/courses/") {
			return
		}
		fmt.Println(link)
		e.Request.Visit(link)
	})

	c.OnHTML("div[class=content-wrapper]", func(e *colly.HTMLElement) {
		if e.DOM.Find("section.course-info").Length() == 0 {
			return
		}
		title := strings.Split(e.ChildText(".course-info__title"), "\n")[0]
		fmt.Println(title)
		course_id := e.ChildAttr("input[name=course_id]", "value")
		fmt.Println(course_id)
		start_date := e.DOM.Find("div.course-info__meta__item").Eq(1).Children().Eq(1).Text()
		fmt.Println("start====>:", start_date)
		end_date := strings.TrimSpace(e.DOM.Find("div.course-info__meta__item").Eq(2).Children().Eq(1).Text())
		fmt.Println("end====>:", end_date)
		number := e.DOM.Find("div.course-info__meta__item").Eq(0).Children().Eq(1).Text()
		fmt.Println("number====>:", number)
		//fmt.Printf("item: %+v\n", item)
		//start_date, _ = time.Parse(DATE_FORMAT, start_date)
		var run string
		if len(strings.Split(course_id, "_")) <= 1 {
			return
		}
		run = strings.Split(course_id, "_")[1]
		fmt.Println(run)
		course := Course{
			CourseID:  course_id,
			Run:       run,
			Name:      title,
			Number:    number,
			StartDate: start_date,
			EndDate:   end_date,
			URL:       e.Request.AbsoluteURL(fmt.Sprintf("/courses/%s/about", course_id)),
		}
		courses = append(courses, course)
		fmt.Println(course)

	})

	c.Visit("https://www.indonesiax.co.id/courses")
	jsonData, err := json.MarshalIndent(courses, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonData))
}
