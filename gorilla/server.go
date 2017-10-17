package main

import (
    "log"
    "time"
    "flag"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
    var dir string

    flag.StringVar(&dir,"dir",".","the directory to serve files from. Default to the current dir")
    flag.Parse()
    r := mux.NewRouter()
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",http.FileServer(http.Dir(dir))))

    srv := &http.Server{
        Handler:  r,
        Addr:     "127.0.0.1:8000",
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    log.Fatal(srv.ListenAndServe())
}
