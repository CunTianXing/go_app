package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
    "github.com/urfave/negroni"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/",func(w http.ResponseWriter, r *http.Request){
        fmt.Fprintf(w,"Welcome to the home page!")
    })
    n := negroni.Classic()
    n.UseHandler(mux)
    
    s := &http.Server{
        Addr:       ":8080",
        Handler:    n,
        ReadTimeout: 10 * time.Second,
        WriteTimeout: 10 * time.Second,
        MaxHeaderBytes: 1 << 20,

    }
    log.Println("http server start:8080")
    if err := s.ListenAndServe(); err != nil {
        log.Fatal(err)

    }

}

