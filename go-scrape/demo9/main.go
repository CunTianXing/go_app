package main

import (
    "fmt"
    "regexp"
    "github.com/gocolly/colly"
    "github.com/gocolly/colly/debug"
)

func main() {
    c := colly.NewCollector()
    c.SetDebugger(&debug.LogDebugger{})
    // Visit only root url and urls which start with "e" or "h" on httpbin.org
    c.URLFilters = []*regexp.Regexp{
        //regexp.MustCompile("http://httpbin\\.org/(e.+)$"),
        regexp.MustCompile("https://www.amazon\\.cn/b.+"),
    }
    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        fmt.Printf("Link found: %q -> %s\n", e.Text, link)
        c.Visit(e.Request.AbsoluteURL(link))
    })

    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL.String())
    })

    c.Visit("https://www.amazon.cn/b/ref=sa_menu_top_pc_l1?ie=UTF8&node=42689071")
}
