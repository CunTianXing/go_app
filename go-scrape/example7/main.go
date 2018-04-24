package main

import (
	"github.com/CunTianXing/go_app/go-scrape/example7/crawl"
	"github.com/PuerkitoBio/goquery"
)

type DummyParser struct{}

func (d DummyParser) ParsePage(doc *goquery.Document) crawl.ScrapeResult {
	data := crawl.ScrapeResult{}
	data.Title = doc.Find("title").First().Text()
	data.H1 = doc.Find("h1").First().Text()
	return data
}

func main() {
	d := DummyParser{}
	startURL := "https://www.theguardian.com/uk"
	crawl.Crawl(startURL, d, 10)
}
