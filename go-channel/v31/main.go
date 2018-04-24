package main

import (
	"fmt"
	"math/rand"
	"time"
)

func longTimeRequest(r chan<- int32) {
	time.Sleep(3 * time.Second)
	r <- rand.Int31n(100)
}

func sumSquares(a, b int32) int32 {
	return a*a + b*b
}

func main() {
	//使用给定的seed来初始化生成器到一个确定的状态。
	rand.Seed(time.Now().UnixNano())

	ra, rb := make(chan int32), make(chan int32)
	go longTimeRequest(ra)
	go longTimeRequest(rb)
	fmt.Println(sumSquares(<-ra, <-rb))

	r := make(chan int32, 2)
	go longTimeRequest(r)
	go longTimeRequest(r)
	fmt.Println(sumSquares(<-r, <-r))
}
