package main

import (
	"fmt"
	"unsafe"
)

//The Sizes Of nil Values With Types Of Different Kinds May Be Different
func main() {
	var p *struct{} = nil
	fmt.Println(unsafe.Sizeof(p)) // 8

	var s []int = nil
	fmt.Println(unsafe.Sizeof(s)) //24

	var s2 []byte = nil
	fmt.Println(unsafe.Sizeof(s2)) //24

	var m map[int]bool = nil
	fmt.Println(unsafe.Sizeof(m)) // 8

	var m2 map[float32]bool = nil
	fmt.Println("sss", unsafe.Sizeof(m2)) //8

	var c chan string = nil
	fmt.Println(unsafe.Sizeof(c)) // 8

	var f func() = nil
	fmt.Println(unsafe.Sizeof(f)) // 8

	var i interface{} = nil
	fmt.Println(unsafe.Sizeof(i)) //16

	//Two nil Values Of Two Different Types May Be Not Comparable
	//mismatched types *int and *bool
	//var _ = (*int)(nil) == (*bool)(nil)
	//mismatched types chan int and chan bool
	//var _ = (chan int)(nil) == (chan bool)(nil)

}
