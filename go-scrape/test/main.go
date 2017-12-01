package main

import (
	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()

	c.Visit("https://www.indonesiax.co.id/courses")
}
