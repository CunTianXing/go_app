package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	stopChan := make(chan struct{})
	wg := sync.WaitGroup{}
	//tickStoppedChan := make(chan struct{})
	wg.Add(1)
	go tick(stopChan, &wg)
	//tockStoppedChan := make(chan struct{})
	wg.Add(1)
	go tock(stopChan, &wg)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("main: received C-c - shutting down")
	fmt.Println("main: telling goroutines to stop")
	close(stopChan)
	//<-tickStoppedChan
	//<-tockStoppedChan
	wg.Wait()
	fmt.Println("main: all goroutines have told us they've finished")
}

func tick(stop chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			fmt.Printf("tick: tick %s\n", now.UTC().Format("20060102-150405.000000000"))
		case <-stop:
			fmt.Println("tick: caller has told us to stop")
			return
		}
	}
}

func tock(stop chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			fmt.Printf("tock: tick %s\n", now.UTC().Format("20060102-150405.000000000"))
		case <-stop:
			fmt.Println("tock: caller has told us to stop")
			return
		}
	}
}
