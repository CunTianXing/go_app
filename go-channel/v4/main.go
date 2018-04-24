package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan struct{})

	//done := make(chan struct{},1)
	//done <- struct{}{}

	go func() {
		fmt.Print("Hello")
		time.Sleep(2 * time.Second)
		<-done
		fmt.Println(222)
	}()
	done <- struct{}{}
	fmt.Println(" world!")
}
