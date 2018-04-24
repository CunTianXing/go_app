package main

import (
	"fmt"
	"math/rand"
	"time"
)

func longTimeRequest() <-chan int32 {
	r := make(chan int32)
	go func() {
		time.Sleep(3 * time.Second)
		r <- rand.Int31n(100)
	}()
	return r
}

func sumSquares(a, b int32) int32 {
	return a*a + b*b
}

func main() {
	//使用给定的seed来初始化生成器到一个确定的状态。
	rand.Seed(time.Now().UnixNano())
	a, b := longTimeRequest(), longTimeRequest()
	fmt.Println(sumSquares(<-a, <-b))
}
