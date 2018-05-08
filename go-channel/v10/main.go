package main

import (
	"fmt"
	"runtime"
	"time"
)

type Xingcuntian struct {
	xch chan string
	bch chan bool
}

func main() {
	x := New()

	for {
		x.sendingGoRoutine()
		time.Sleep(time.Second * 2)
		x.sendingGoRoutineBool()
	}
}

func New() (x *Xingcuntian) {
	x = &Xingcuntian{
		xch: make(chan string, 10),
		bch: make(chan bool, 10),
	}
	fmt.Println(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go x.receivingGoRoutine()
	}
	return x
}

func (x *Xingcuntian) sendingGoRoutine() {
	x.xch <- "xingcuntian.com"
}

func (x *Xingcuntian) sendingGoRoutineBool() {
	x.bch <- true
}

func (x *Xingcuntian) receivingGoRoutine() {
	for {
		select {
		case v := <-x.xch:
			fmt.Println("Received value ", v)
		case u := <-x.bch:
			fmt.Println("Received bool ", u)
		default:
		}
	}
}
