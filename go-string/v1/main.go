package main

import (
    "fmt"
)

func main() {
    const World = "world"
    var hello = "hello"

    //concat strings.
    var helloWorld = hello + " " + World
    helloWorld += "!"
    fmt.Println(helloWorld)

    //Compare strings.
    fmt.Println(hello == "hello")
    fmt.Println(hello > helloWorld)
}
