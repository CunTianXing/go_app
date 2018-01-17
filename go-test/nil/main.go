package main

import (
	"fmt"
)

func main() {
	// fmt.Println("nil")
	// fmt.Printf("nil type %T\n", nil) //nil type <nil>
	// var m map[int]string
	// var ptr *int
	// fmt.Printf("%p\n", m)
	// fmt.Printf("%p\n", ptr)

	//fmt.Printf(m == ptr)

	// var m2 map[string]string
	// fmt.Printf("m value %v\n", m2)
	// var s2 []string
	// fmt.Printf("s2 value %v\n", s2)
	// var ch chan int
	// fmt.Printf("ch value %v\n", ch)
	// var ptr2 *int
	// fmt.Printf("ptr2 value %v\n", ptr2)
	// var f func()
	// fmt.Printf("%T\n", f)
	//
	// //fmt.Printf("%#v\n", f())
	// var i interface{}
	// fmt.Printf("i value %v\n", i)

	var m2 map[int]string
	var ptr2 *int
	var c chan int
	var sl []int
	var f func()
	var i interface{}
	fmt.Printf("%#v\n", m2)
	if m2 == nil {
		fmt.Println("ddddd")
	}
	fmt.Printf("%#v\n", ptr2)
	if ptr2 == nil {
		fmt.Println("ppppp")
	}
	fmt.Printf("%#v\n", c)
	if c == nil {
		fmt.Println("ccccc")
	}
	fmt.Printf("%#v\n", sl)
	if sl == nil {
		fmt.Println("ssssssl")
	}
	fmt.Printf("%#v\n", f)
	if f == nil {
		fmt.Println("fffff")
	}
	fmt.Printf("%#v\n", i)
	if i == nil {
		fmt.Println("iiiiii")
	}

	m3 := make(map[string]int)
	fmt.Printf("%#v\n", m3)
	fmt.Println(m3)
	if m3 == nil {
		fmt.Println("m3")
	}
	s3 := []string{}
	fmt.Printf("%#v\n", s3)
	fmt.Println(s3)
	if s3 == nil {
		fmt.Println("s3")
	}

	var m4 map[int]string
	m4 = map[int]string{}
	fmt.Printf("%#v\n", m4)
	fmt.Println(m4)
	if m4 == nil {
		fmt.Println("m4")
	}
}

//1. nil 是不能比较的
//fmt.Println(nil == nil)

//2. 默认 nil 是 typed
//fmt.Printf("nil type %T", nil) //nil type <nil>

//3.不同类型 nil 的地址是一样的
//m 和 ptr 的地址都是 0x0

//4.不同类型的 nil 是不能比较的
//fmt.Printf(m == ptr)

//5.nil 是 map，slice，pointer，channel，func，interface 的零值

//*** zero value 是 go 中变量在声明之后,但是未初始化被赋予的该类型的一个默认值。
