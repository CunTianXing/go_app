package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var x struct {
		a int64
		z byte

		c string
		b bool
	}
	const M, N = unsafe.Sizeof(x.c), unsafe.Sizeof(x)
	fmt.Printf("M=%d, N=%d\n", M, N) //M=16, N=32
	fmt.Println("======Alignof======")
	fmt.Println(unsafe.Alignof(x.a)) //8

	fmt.Println(unsafe.Alignof(x.z)) //1
	fmt.Println(unsafe.Alignof(x.c)) //8
	fmt.Println(unsafe.Alignof(x.b)) //1
	fmt.Println("=====Offsetof======")
	fmt.Println(unsafe.Offsetof(x.a)) //0
	fmt.Println(unsafe.Offsetof(x.z)) //8
	fmt.Println(unsafe.Offsetof(x.c)) //16
	fmt.Println(unsafe.Offsetof(x.b)) //32

	var y struct {
		a int8
		b byte
		d float32
		e float32
		f int64
		c bool
	}
	fmt.Println(unsafe.Sizeof(y.d))
	const H, L = unsafe.Sizeof(y.a), unsafe.Sizeof(y)
	fmt.Printf("M=%d, N=%d\n", H, L)  //M=1, N=32
	fmt.Println(unsafe.Offsetof(y.a)) //0
	fmt.Println(unsafe.Offsetof(y.b)) //1
	fmt.Println(unsafe.Offsetof(y.d)) //4
	fmt.Println(unsafe.Offsetof(y.e)) //8
	fmt.Println(unsafe.Offsetof(y.f)) //16
	fmt.Println(unsafe.Offsetof(y.c)) //24

}
