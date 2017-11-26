package main

import (
    "log"
    "github.com/gocolly/colly"
)

func main() {
    // create a  new collector
    c := colly.NewCollector()
    // authenticate
    err := c.Post("https://gocn.io/account/login/",map[string]string{"user_name":"","password":""})
    if err != nil {
        log.Fatal(err)
    }

    // attach callbacks after login
    c.OnResponse(func(r *colly.Response) {
        log.Println("response received", r.StatusCode)
    })

    // start scraping
    c.Visit("https://gocn.io/")
}
