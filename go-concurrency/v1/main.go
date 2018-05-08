package main

import "fmt"

func main() {
	ch := make(chan int)

	go func() {
		ch <- 45
	}()
	fmt.Println("received :", <-ch)
}
