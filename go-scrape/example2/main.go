package main

import (
    "encoding/csv"
    "log"
    "os"
    "github.com/gocolly/colly"
    "github.com/gocolly/colly/debug"
)

func main() {
    fName := "xkcd_store_items.csv"
    file, err := os.Create(fName)
    if err != nil {
        log.Fatalf("Cannot create file %q: %s\n", fName, err)
        return
    }
    defer file.Close()
    writer := csv.NewWriter(file)
    defer writer.Flush()
    writer.Write([]string{"Name","Price","URL","Image URL"})
    c := colly.NewCollector()
    c.SetDebugger(&debug.LogDebugger{})
    c.AllowedDomains = []string{"store.xkcd.com"}
    c.OnHTML(".product-grid-item",func(e *colly.HTMLElement) {
        writer.Write([]string{
            e.ChildAttr("a","title"),
            e.ChildText("span"),
            e.Request.AbsoluteURL(e.ChildAttr("a","href")),
            "https"+ e.ChildAttr("img","src"),
        })
    })
    c.OnHTML(`.next a[href]`, func(e *colly.HTMLElement) {
        e.Request.Visit(e.Attr("href"))
    })
    c.Visit("https://store.xkcd.com/collections/everything")
    log.Printf("Scraping finished, check file %q for results\n", fName)

}
