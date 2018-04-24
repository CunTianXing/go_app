package main

//#include "./hello.h"
import "C"

func main() {
	C.SayHello(C.CString("hello, world\n"))
}

//C代码模块化 - 改用Go重写C模块
//函数参数去掉 const 修饰符
