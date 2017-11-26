package main

import (
    "fmt"
    "github.com/gocolly/colly"
)

func main() {
    // Create a collector
    c := colly.NewCollector()
    // Set HTML callback
    // Won't be called if error occurs
    c.OnHTML("*",func(e *colly.HTMLElement) {
        fmt.Println(e)
    })

    // Set error handler
    c.OnError(func(r *colly.Response, err error) {
        fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
    })

    // Start scraping
    c.Visit("https://www.amazon.cn/dp/B076M4XQ5K?_encoding=UTF8&ref_=pc_cxrd_658390051_newTab_658390051_a_new_2&pf_rd_p=efc7b5e2-17f0-4684-8f11-b51e9cf3855e&pf_rd_s=merchandised-search-9&pf_rd_t=101&pf_rd_i=658390051&pf_rd_m=A1AJ19PSB66TGU&pf_rd_r=X92DSYYSPJTAQVCKAGY4&pf_rd_r=X92DSYYSPJTAQVCKAGY4&pf_rd_p=efc7b5e2-17f0-4684-8f11-b51e9cf3855e")
}
