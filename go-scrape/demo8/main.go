package main

import (
    "fmt"
    "github.com/gocolly/colly"
    "github.com/gocolly/colly/debug"
)

func main() {
    c := colly.NewCollector()
    c.SetDebugger(&debug.LogDebugger{})
    c.OnRequest(func(r *colly.Request) {
        r.Ctx.Put("url", r.URL.String())
    })

    c.OnResponse(func(r *colly.Response) {
        fmt.Println(r.Ctx.Get("url"))
    })

    c.Visit("https://en.wikipedia.org/")
}

