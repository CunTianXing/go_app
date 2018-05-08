package main

import "time"
import "fmt"

func timer(d time.Duration) <-chan int {
	c := make(chan int)
	go func() {
		time.Sleep(d)
		fmt.Println("write")
		c <- 1
	}()
	return c
}

func main() {
	for i := 0; i < 24; i++ {
		c := timer(1 * time.Second)
		fmt.Println("received: ", <-c)
	}
}
