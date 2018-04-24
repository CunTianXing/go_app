// +build go1.10
package main

// extern void SayHello(_GoString_ s);
import "C"
import "fmt"

func main() {
	C.SayHello("hello world\n")
}

//export SayHello
func SayHello(s string) {
	fmt.Print(s)
}

// Go 1.10 增加了 _GoString_ 类型
// _GoString_ 是预定义的类型, 和 GoString 等价
// 避免手写函数声明时出现循环依赖
