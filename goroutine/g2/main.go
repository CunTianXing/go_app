package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(1)

	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("starting Go Routines")

	go func() {
		defer wg.Done()

		for char := 'a'; char < 'a'+26; char++ {
			fmt.Printf("%c ", char)
		}
	}()

	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Microsecond)
		for number := 1; number < 27; number++ {
			fmt.Printf("%d ", number)
		}
	}()
	fmt.Println("Wating To Finish")

	wg.Wait()
	fmt.Println("\n Terminating Program")
}
