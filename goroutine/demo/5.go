package main

import (
	"fmt"
	"sync"
)

func main() {
	messages := make(chan int)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		//time.Sleep(time.Second * 3)
		messages <- 1
	}()
	go func() {
		defer wg.Done()
		//time.Sleep(time.Second * 2)
		messages <- 2
	}()
	go func() {
		defer wg.Done()
		//time.Sleep(time.Second * 1)
		messages <- 3
	}()
	go func() {
		for i := range messages {
			fmt.Println(i)
		}
	}()

	wg.Wait()
}
