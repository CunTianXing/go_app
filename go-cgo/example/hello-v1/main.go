package main

//#include <stdio.h>
import "C"

func main() {
	C.puts(C.CString("hello world"))
}

// import "C" 表示启用 CGO
// import "C" 前的注释表示包含C头文件: <stdio.h>
// C.CString 表示将 Go 字符串转为 C 字符串
// C.puts 调用C语言的puts函数输出 C 字符串
