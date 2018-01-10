package main

import (
	"fmt"
	"sync"
)

func ingest() <-chan []string {
	out := make(chan []string)
	go func() {
		out <- []string{"aaaa", "bbbb"}
		out <- []string{"cccc", "dddd"}
		out <- []string{"eeee", "ffff"}
		close(out)
	}()
	return out
}

func process(concurrency int, in <-chan []string) <-chan int {
	var wg sync.WaitGroup
	wg.Add(concurrency)
	out := make(chan int)
	work := func() {
		for data := range in {
			for _, word := range data {
				out <- len(word)
			}
		}
		wg.Done()
	}

	go func() {
		for i := 0; i < concurrency; i++ {
			go work()
		}
	}()

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func store(in <-chan int) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for data := range in {
			fmt.Println(data)
		}
	}()
	return done
}

func main() {
	in := ingest()
	reduced := process(4, in)
	<-store(reduced)
}
