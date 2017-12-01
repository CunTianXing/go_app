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
		fmt.Println("store start")
		r.Ctx.Put("url", r.URL.String())
		fmt.Println("store end")
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.Ctx.Get("url"))
	})

	c.Visit("https://en.wikipedia.org/")
}
