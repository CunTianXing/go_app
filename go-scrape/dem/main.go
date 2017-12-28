package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	urlRoot = "http://www.galerie-art-paris.com/"
)

var wg sync.WaitGroup

func gatherNodes(n *html.Node) bool {
	if n.DataAtom == atom.A && n.Parent != nil {
		return scrape.Attr(n.Parent, "class") == "menu"
	}
	return false
}

func scrapGalleries(url string) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}
	matcher := func(n *html.Node) bool {
		return n.DataAtom == atom.Span && scrape.Attr(n, "class") == "galerie-art-titre"
	}
	for _, g := range scrape.FindAll(root, matcher) {
		fmt.Println(scrape.Text(g))
	}
}

func main() {
	resp, err := http.Get(urlRoot)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	as := scrape.FindAll(root, gatherNodes)
	for _, link := range as {
		wg.Add(1)
		go scrapGalleries(urlRoot + scrape.Attr(link, "href"))
	}
	wg.Wait()
}
