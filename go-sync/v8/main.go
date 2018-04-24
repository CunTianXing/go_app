package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	fooIsDone = false
	barIsDone = false
	cond      = sync.NewCond(&sync.Mutex{})
)

func doFoo() {
	time.Sleep(time.Second)
	cond.L.Lock()
	fooIsDone = true
	cond.Signal()
	cond.L.Unlock()
}

func doBar() {
	time.Sleep(time.Second * 2)
	cond.L.Lock()
	barIsDone = true
	cond.L.Unlock()
	cond.Signal()
}

func main() {
	cond.L.Lock()
	go doFoo()
	go doBar()

	checkCondition := func() bool {
		fmt.Println(fooIsDone, barIsDone)
		return fooIsDone && barIsDone
	}
	for !checkCondition() {
		cond.Wait()
	}
	cond.L.Unlock()
}
