package main

// extern void SayHello(char* s);
import "C"
import "fmt"

func main() {
	C.SayHello(C.CString("Hello, World\n"))
}

//export SayHello
func SayHello(s *C.char) {
	fmt.Print(C.GoString(s))
}

//
// C 语言版本 SayHello 函数实现只存在于心中
// 面向纯 C 接口的 Go 语言编程
