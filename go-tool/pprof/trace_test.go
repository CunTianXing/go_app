package trace_110

import (
	"testing"
	"sync"
	"fmt"
	"time"
	"net/http"
)

func takeCPU(start time.Time, wg *sync.WaitGroup) {
	defer wg.Done()
	j := 3
	for time.Since(start) < time.Second {
		for i :=1; i < 1000000; i++ {
			j *= i
		}
	}
	fmt.Println(j)
}

func takeTimeOnly(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Second * 3)
}

func takeIO(start time.Time, wg *sync.WaitGroup) {
	defer wg.Done()
	errCount := 0
	for time.Since(start) < time.Second * 4{
		_, err := http.Get("https://www.google.com")
		if err != nil {
			errCount++
		}
	}
	fmt.Println(errCount)
}

func startServer() {
	wg := sync.WaitGroup{}
	wg.Add(3)
	start := time.Now()
	go takeCPU(start, &wg)
	go takeTimeOnly(&wg)
	go takeIO(start, &wg)
	wg.Wait()
}


func TestServer(t *testing.T) {
	startServer()
}
