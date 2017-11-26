package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "time"
    "github.com/gocolly/colly"
)

func generateFormData() map[string][]byte {
    f, _ := os.Open("image.jpg")
    defer f.Close()
    imgData, _ := ioutil.ReadAll(f)
    return map[string][]byte{
        "firstname": []byte("one"),
        "lastname":  []byte("two"),
        "email":     []byte("this@qq.com"),
        "file":      imgData,
    }
}

func setupServer() {
    var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("received request")
        err := r.ParseMultipartForm(1000000)
        if err != nil {
            fmt.Println("server: Error")
            w.WriteHeader(500)
            w.Write([]byte("<html><body>Internal Server Error</body></html>"))
            return
        }
        w.WriteHeader(200)
        fmt.Println("server: OK")
        w.Write([]byte("<html><body>Success</body></html>"))
    }
    go http.ListenAndServe(":9090",handler)
}

func main() {
    setupServer()
    c := colly.NewCollector()
    c.AllowURLRevisit = true
    c.MaxDepth = 5
    c.OnHTML("html",func(e *colly.HTMLElement){
        fmt.Println(e.Text)
        time.Sleep(1 * time.Second)
        e.Request.PostMultipart("http://localhost:9090/",generateFormData())
    })
    c.OnRequest(func(r *colly.Request){
        fmt.Println("Posting image.jpg to ", r.URL.String())
    })

    c.PostMultipart("http://localhost:9090/",generateFormData())
    c.Wait()
}





