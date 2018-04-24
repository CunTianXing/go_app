package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	m sync.RWMutex
	n uint64
}

func (c *Counter) Increase(delta uint64) {
	c.m.Lock()
	c.n += delta
	c.m.Unlock()
}

func (c *Counter) Value() uint64 {
	c.m.Lock()
	defer c.m.Unlock()
	return c.n
}

func main() {
	var c Counter
	for i := 0; i < 10; i++ {
		go func() {
			for k := 0; k < 10; k++ {
				c.Increase(1)
			}
		}()
		for c.Value() < 10 {
		}
		fmt.Println(c.Value())
	}
}
