package main

import (
	"fmt"
	"time"
)

func main() {
	a()
	fmt.Println("normally returned from main")
}

// Inside A
// Inside B
// panic: oh! B panicked
//
// goroutine 5 [running]:
// main.b()
// 	/data/go/src/github.com/CunTianXing/go_app/go-test/Panic&&Recover/main2.go:28 +0x83
// created by main.a
// 	/data/go/src/github.com/CunTianXing/go_app/go-test/Panic&&Recover/main2.go:22 +0x9f
// exit status 2
func recovery() {
	if r := recover(); r != nil {
		fmt.Println("recovered:", r)
	}
}

func a() {
	defer recovery()
	fmt.Println("Inside A")
	go b()
	time.Sleep(1 * time.Second)
}

func b() {
	fmt.Println("Inside B")
	panic("oh! B panicked")
}
