package main

import (
    "fmt"
    "github.com/gocolly/colly"
    "github.com/juju/persistent-cookiejar"
)

func main() {
    // create a new collector
    c := colly.NewCollector()
    // Reduce maximum response body size to 1M
    c.MaxBodySize = 1024 * 1024
    // Don't track visited urls automatically
    c.AllowURLRevisit = true
    // Turn off cookie handling
    //c.DisableCookies()
    j, err := cookiejar.New(&cookiejar.Options{Filename: "cookie.db"})
    if err == nil {
        c.SetCookieJar(j)
    }

    c.MaxDepth = 2

    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        fmt.Println(link)
        e.Request.Visit(link)
    })

    //start scraping
    c.Visit("https://www.amazon.cn/b/ref=sa_menu_top_applia_l2_b423695071?ie=UTF8&node=423695071")
}
