package main

import (
    "fmt"
    "github.com/gocolly/colly"
)

func main() {
    c := colly.NewCollector()
    c.Limit(&colly.LimitRule{DomainGlob:"*", Parallelism: 5})
    c.MaxDepth = 2
    c.OnHTML("a[href]",func(e *colly.HTMLElement) {
        link := e.Attr("href")
        fmt.Println(link)
        go e.Request.Visit(link)
    })
    c.Visit("https://en.wikipedia.org/")
    c.Wait()
}
