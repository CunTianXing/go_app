package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	const N = 5
	var values [N]int32

	var wg = &sync.WaitGroup{}
	wg.Add(1)
	for i := 0; i < N; i++ {
		i := i
		go func() {
			wg.Wait()
			fmt.Printf("values[%v]=%v \n", i, values[i])
		}()
	}
	for i := 0; i < N; i++ {
		values[i] = 50 + rand.Int31n(50)
	}
	wg.Done()
}
