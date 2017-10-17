package main

import (
	"log"
	"net/http"
    "io/ioutil"
)



func main(){
	http.HandleFunc("/500",func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(500)
		w.Write([]byte("NOT-OK"))
    })
	go http.ListenAndServe(":8080",nil)

	resp, err := http.Get("http://localhost:8080/500")

    if err != nil {
		log.Fatal(err)

	}

    if resp.StatusCode != http.StatusOK {
        b, _ := ioutil.ReadAll(resp.Body)
        log.Fatal(string(b))
	}
}
