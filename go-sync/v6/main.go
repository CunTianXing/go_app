package main

import (
	"fmt"
	"sync"
)

func main() {
	const N = 10
	var values [N]string

	var m sync.RWMutex

	for i := 0; i < N; i++ {
		i := i
		go func() {
			m.RLock()
			values[i] = string('a' + i)
			m.RUnlock()
		}()
	}

	done := func() bool {
		m.Lock()
		defer m.Unlock()
		for i := 0; i < N; i++ {
			if values[i] == "" {
				return false
			}
		}
		fmt.Println(values)
		return true
	}

	for !done() {
	}
}
