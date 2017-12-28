package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

//atfResults ul li
//pagn
//pagnDisabled
func main() {
	fName := "cryptocoinmarketcap.csv"
	file, err := os.OpenFile(fName, os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	c := colly.NewCollector()
	//c.SetDebugger(&debug.LogDebugger{})
	//s-result-list s-col-3 s-result-list-hgrid s-height-equalized s-grid-view s-text-condensed
	c.OnHTML("div#atfResults ul", func(e *colly.HTMLElement) {
		fmt.Println(e.DOM.Children())
		//fmt.Println(e.DOM.Children().Eq(5).Children().Eq(0).Attr("href"))
	})
	//下一页
	// c.OnHTML("div[id=pagn] > span[class=pagnRA] > a[href]", func(e *colly.HTMLElement) {
	// 	fmt.Println(e.Attr("href"))
	// 	e.Request.Visit(e.Attr("href"))
	//
	// })
	c.Visit("https://www.amazon.cn/s?field-keywords=%E6%89%8B%E6%9C%BA%E5%A3%B3")
}
