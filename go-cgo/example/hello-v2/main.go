package main

/*
#include <stdio.h>

static void SayHello(const char* s){
    puts(s);
}
*/
import "C"

func main() {
	C.SayHello(C.CString("hello world\n"))
}

//调用自定义的C函数
