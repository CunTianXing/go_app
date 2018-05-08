package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	stopChan := make(chan struct{})

	tickStoppedChan := make(chan struct{})
	go tick(stopChan, tickStoppedChan)
	tockStoppedChan := make(chan struct{})
	go tock(stopChan, tockStoppedChan)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("main: received C-c - shutting down")
	fmt.Println("main: telling goroutines to stop")
	close(stopChan)
	<-tickStoppedChan
	<-tockStoppedChan
	fmt.Println("main: all goroutines have told us they've finished")
}

func tick(stop, stopped chan struct{}) {
	defer close(stopped)

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

func tock(stop, stopped chan struct{}) {
	defer close(stopped)

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
