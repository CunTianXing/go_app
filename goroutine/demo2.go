package main

import (
	"fmt"
	"runtime"
	"sync"
)

var count int = 0

func counter(lock *sync.Mutex) {
	lock.Lock()
	count++
	fmt.Println(count)
	lock.Unlock()
}

func main() {
	lock := &sync.Mutex{}
	for i := 0; i < 10; i++ {
		//传递指针是为了防止 函数内的锁和 调用锁不一致
		go counter(lock)
	}
	for {
		lock.Lock()
		c := count
		lock.Unlock()
		///把时间片给别的goroutine  未来某个时刻运行该routine
		runtime.Gosched()
		if c >= 10 {
			fmt.Println("goroutine end")
			break
		}
	}
}
