package main

import (
	"fmt"
)

func main() {
	messages := make(chan int)
	done := make(chan bool)

	go func() {
		//time.Sleep(time.Second * 3)
		messages <- 1
		done <- true
	}()
	go func() {
		//time.Sleep(time.Second * 2)
		messages <- 2
		done <- true
	}()
	go func() {
		//time.Sleep(time.Second * 1)
		messages <- 3
		done <- true
	}()
	go func() {
		for i := range messages {
			fmt.Println(i)
		}
	}()
	for i := 0; i < 3; i++ {
		<-done
	}
}
