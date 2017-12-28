package main

import (
	"fmt"
	"net/http"
)

func main() {
	go DeferEnd()
	http.ListenAndServe(":8080", nil)
}

func DeferEnd() bool {
	defer fmt.Println("defer1")
	Defer3()
	return true
}

func Defer3() bool {
	fmt.Println("defer2")
	return true
}
