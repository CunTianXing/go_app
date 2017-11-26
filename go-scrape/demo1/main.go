package main

import (
    "fmt"
    "github.com/gocolly/colly"
)

func main() {
    //Instantiate default collector
    c := colly.NewCollector()

    //Visit only domains: hackerspaces.org,
    c.AllowedDomains = []string{"hackerspaces.org", "wiki.hackerspaces.org"}

    // On every a element which has href attribute call callback
    c.OnHTML("a[href]",func(e *colly.HTMLElement) {
        link := e.Attr("href")
        fmt.Printf("Link found: %q -> %s\n", e.Text, link)
        c.Visit(e.Request.AbsoluteURL(link))
    })

    c.OnRequest(func(r *colly.Request){
        fmt.Println("Visiting", r.URL.String())
    })

    c.Visit("https://hackerspaces.org/")
}
