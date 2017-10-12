package main

import (
    "fmt"
    "net/http"
    "github.com/urfave/negroni"
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/",HomeHandler)
    r.HandleFunc("/products/{key}",ProductsHandler)
    n := negroni.Classic()
    n.UseHandler(r)
    n.Run(":3000")
}

func HomeHandler(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w,"Welcome home page")
}

func ProductsHandler(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    key := vars["key"]
    fmt.Fprintf(w,"get key: %s\n",key)
}
