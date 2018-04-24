package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var m sync.Mutex
	m.Lock()
	go func() {
		time.Sleep(time.Second)
		fmt.Println("Hi.")
		m.Unlock()
	}()
	m.Lock()
	fmt.Println("Bye.")
	//m.Unlock()
}
