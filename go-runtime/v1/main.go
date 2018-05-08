package main

import "fmt"

func showNumber(i int) {
    fmt.Println(i)
}

func main() {
    for i := 0; i < 10; i++ {
        go showNumber(i)
    }

    fmt.Println("gary")
}
