package main

import (
    "fmt"
    "github.com/gocolly/colly"
    "github.com/gocolly/colly/debug"
)

func main() {
    url := "https://httpbin.org/delay/2"

    c := colly.NewCollector()

    c.SetDebugger(&debug.LogDebugger{})

    c.Limit(&colly.LimitRule{
        DomainGlob: "*httpbin.*",
        Parallelism: 2,
        //Delay: 5 * time.Second,
    })

    for i:=0; i<4; i++ {
        go c.Visit(fmt.Sprintf("%s?n=%d",url,i))
    }
    c.Visit(url)
    c.Wait()
}
