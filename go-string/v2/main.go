package main

import (
    "fmt"
    "strings"
)

func main(){
    var helloworld = "hello world!"
    var hello = helloworld[:5]
    fmt.Println(hello[0])
    fmt.Printf("%T\n",hello[0])//uint8

    fmt.Println(len(hello),len(helloworld))
    fmt.Println(strings.HasPrefix(helloworld,hello))
}
