package main

import (
	"fmt"
	"sync"
	"time"
)

func f(wg *sync.WaitGroup, val string) {
	time.Sleep(3 * time.Second)
	fmt.Printf("Finished: %v - %v\n", val, time.Now())
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	go f(&wg, "groutine A")

	go func(wg *sync.WaitGroup, val string) {
		time.Sleep(3 * time.Second)
		fmt.Printf("Finished: %v - %v\n", val, time.Now())
		wg.Done()
	}(&wg, "goroutine B")

	go f(&wg, "goroutine C")

	wg.Wait()
	fmt.Printf("Finished all goroutines: %v\n", time.Now())
}
