package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type Mail struct {
	Title   string
	Link    string
	Author  string
	Date    string
	Message string
}

func main() {
	var groupName string
	flag.StringVar(&groupName, "group", "go-kit", "Google Groups group name")
	flag.Parse()

	//threads := make(map[string][]Mail)

	threadCollector := colly.NewCollector()
	//threadCollector.SetDebugger(&debug.LogDebugger{})
	//mailCollector := colly.NewCollector()

	threadCollector.OnHTML("tr", func(e *colly.HTMLElement) {
		// <tbody>
		// 	<tr>
		// 	    <td class="subject"><a href="https://groups.google.com/d/topic/go-kit/hTQThEhvdTo" title="HTTP Redirects Handling strategy with go-kit endpoints">HTTP Redirects Handling strategy with go-kit endpoints</a></td>
		//     	<td class="author"><span>Joel Unzain</span></td>
		// 	    <td class="lastPostDate">17-10-2</td>
		// 	</tr>
		//  </tbody>
		ch := e.DOM.Children()
		//fmt.Printf("data: %+v\n", *ch)
		author := ch.Eq(1).Text()
		if author == "" {
			return
		}
		fmt.Println(author)
		title := ch.Eq(0).Text()
		fmt.Println("subject:", title)
		link, _ := ch.Eq(0).Children().Eq(0).Attr("href")
		link = strings.Replace(link, ".com/d/topic", ".com/forum/?_escaped_fragment_=topic", 1)
		fmt.Println("link:", link)
	})

	threadCollector.Visit("https://groups.google.com/forum/?_escaped_fragment_=forum/" + groupName)
}
