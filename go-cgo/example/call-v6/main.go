package main

// extern void SayHello(GoString s);
import "C"
import "fmt"

func main() {
	C.SayHello("ddddd")
}

//export SayHello
func SayHello(s string) {
	fmt.Print(s)
}

//导出Go函数
// GoString 在哪定义?
//导出函数的参数是Go字符串
//C类型为 GoString, 在 _cgo_export.h 文件定义
//要使用 GoString 类型就要引用 _cgo_export.h 文件
//这时候该如何手写 SayHello 函数的声明?
