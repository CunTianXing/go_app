package main

import (
	"fmt"
	"time"
)

func main() {
	message := make(chan int)
	go func() {
		time.Sleep(time.Second * 3)
		message <- 1
	}()
	go func() {
		time.Sleep(time.Second * 2)
		message <- 2
	}()
	go func() {
		time.Sleep(time.Second * 1)
		message <- 3
	}()
	go func() {
		for i := range message {
			fmt.Println(i)
		}
	}()
	time.Sleep(5 * time.Second)

}
