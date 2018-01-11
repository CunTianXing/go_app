package main

import (
	"fmt"
	"runtime/debug"
)

// Recovered runtime error: index out of range
// goroutine 1 [running]:
// runtime/debug.Stack(0xc42000c018, 0xc42003bdc0, 0x2)
// 	/data/golang/go/src/runtime/debug/stack.go:24 +0xa7
// runtime/debug.PrintStack()
// 	/data/golang/go/src/runtime/debug/stack.go:16 +0x22
// main.r()
// 	/data/go/src/github.com/CunTianXing/go_app/go-test/Panic&&Recover/main23.go:11 +0xb3
// panic(0x10abb60, 0x1139220)
// 	/data/golang/go/src/runtime/panic.go:491 +0x283
// main.a()
// 	/data/go/src/github.com/CunTianXing/go_app/go-test/Panic&&Recover/main23.go:18 +0x68
// main.main()
// 	/data/go/src/github.com/CunTianXing/go_app/go-test/Panic&&Recover/main23.go:23 +0x22
// normally returned from main

func r() {
	if r := recover(); r != nil {
		fmt.Println("Recovered", r)
		debug.PrintStack()
	}
}

func a() {
	defer r()
	n := []int{5, 7, 4}
	fmt.Println(n[3])
	fmt.Println("normally returned from a")
}

func main() {
	a()
	fmt.Println("normally returned from main")
}
