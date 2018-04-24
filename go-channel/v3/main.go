package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"sort"
)

func main() {
	values := make([]byte, 32*1024*1024)
	//本函数是一个使用io.ReadFull调用Reader.Read的辅助性函数。当且仅当err == nil时，返回值n == len(b)。
	if _, err := rand.Read(values); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	done := make(chan struct{})
	go func() {
		sort.Slice(values, func(i, j int) bool {
			return values[i] < values[j]
		})
		done <- struct{}{}
	}()
	<-done
	//0 255
	fmt.Println(len(values) - 1)

	fmt.Println(values[0], values[len(values)-1])
}
