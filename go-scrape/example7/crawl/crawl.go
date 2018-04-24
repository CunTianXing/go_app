package crawl

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ScrapeResult struct {
	URL   string
	Title string
	H1    string
}

//Defining Our Parser Interface
type Parser interface {
	ParsePage(*goquery.Document) ScrapeResult
}

//Making HTTP Requests
func getRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func extractLinks(doc *goquery.Document) []string {
	foundUrls := []string{}
	if doc != nil {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			url, exist := s.Attr("href")
			if exist {
				foundUrls = append(foundUrls, url)
			}
		})
		return foundUrls
	}
	return foundUrls
}

func resolveRelative(baseURL string, hrefs []string) []string {
	internalUrls := []string{}
	// func HasPrefix(s, prefix string) bool
	// 判断s是否有前缀字符串prefix。
	// func HasSuffix(s, suffix string) bool
	// 判断s是否有后缀字符串suffix。
	for _, href := range hrefs {
		if strings.HasPrefix(href, baseURL) {
			internalUrls = append(internalUrls, href)
		}
		if strings.HasPrefix(href, "/") {
			resolvedURL := fmt.Sprintf("%s%s", baseURL, href)
			internalUrls = append(internalUrls, resolvedURL)
		}
	}
	return internalUrls
}

//Crawling A Page
func crawlPage(baseURL, targetURL string, parser Parser, token chan struct{}) ([]string, ScrapeResult) {
	token <- struct{}{}
	fmt.Println("Requesting: ", targetURL)
	res, err := getRequest(targetURL)
	if err != nil {
		log.Println(err)
	}
	<-token
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Println(err)
	}
	pageResults := parser.ParsePage(doc)
	links := extractLinks(doc)
	foundUrls := resolveRelative(baseURL, links)
	return foundUrls, pageResults
}

//Getting Base URL
func parseStartURL(u string) string {
	parsed, err := url.Parse(u)
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
}

//Crawl Function
func Crawl(startURL string, parser Parser, concurrency int) []ScrapeResult {
	results := []ScrapeResult{}
	worklist := make(chan []string)
	var n int
	n++
	var tokens = make(chan struct{}, concurrency)
	go func() { worklist <- []string{startURL} }()
	seen := make(map[string]bool)
	baseDomain := parseStartURL(startURL)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(baseDomain, link string, parser Parser, token chan struct{}) {
					foundLinks, pageResults := crawlPage(baseDomain, link, parser, token)
					results = append(results, pageResults)
					if foundLinks != nil {
						worklist <- foundLinks
					}
				}(baseDomain, link, parser, tokens)
			}
		}
	}
	return results
}
