package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	//  Go语言字符串的底层结构在reflect.StringHeader中定义：
	// type StringHeader struct {
	//     Data uintptr
	//     Len  int
	// }
	var data = [...]byte{'h', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd'}
	fmt.Printf("data: %#v\n", data) //data: [12]uint8{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64}
	var str string
	str = "hello, world"
	fmt.Printf("data: %#v\n", str)
	//字符串虽然不是切片，但是支持切片操作，不同位置的切片底层也访问的同一块内存数据（因为字符串是只读的，相同的字符串面值常量通常是对应同一个字符串常量）：
	s := "hello, world"
	hello := s[:5] // 左开右闭
	fmt.Printf("data: %#v\n", hello)
	world := s[7:]
	fmt.Printf("data: %#v\n", world)
	s1 := "hello, world"[:5]
	fmt.Printf("data: %#v\n", s1) //data: "hello"
	s2 := "hello, world"[7:]
	fmt.Printf("data: %#v\n", s2)                                             //data: "world"
	fmt.Println("len(s):", (*reflect.StringHeader)(unsafe.Pointer(&s)).Len)   // 12
	fmt.Println("len(s1):", (*reflect.StringHeader)(unsafe.Pointer(&s1)).Len) // 5
	fmt.Println("len(s2):", (*reflect.StringHeader)(unsafe.Pointer(&s2)).Len) // 5

	fmt.Printf("params: %#v\n", []byte("Hello, 世界"))
	//params: []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c}
	//分析可以发现0xe4, 0xb8, 0x96对应中文“世”，0xe7, 0x95, 0x8c对应中文“界”。我们也可以在字符串面值中直指定UTF8编码后的值（源文件中全部是ASCII码，可以避免出现多字节的字符）。
	fmt.Println("\xe4\xb8\x96") //世
	fmt.Println("\xe7\x95\x8c") //界
	for i, r := range "Hello, 世界" {
		fmt.Printf("%d\t%q\t%d\n", i, r, r)
	}
	// 0	'H'	72
	// 1	'e'	101
	// 2	'l'	108
	// 3	'l'	108
	// 4	'o'	111
	// 5	','	44
	// 6	' '	32
	// 7	'世'	19990
	// 10	'界'	30028
	fmt.Println("\xe4\x00\x00\xe7\x95\x8cabc") //?界abc
	for i, c := range "\xe4\x00\x00\xe7\x95\x8cabc" {
		fmt.Println(i, c)
	}
	// 0 65533  // \uFFFD, 对应 �
	// 1 0      // 空字符
	// 2 0      // 空字符
	// 3 30028  // 界
	// 6 97     // a
	// 7 98     // b
	// 8 99     // c
}
