package main

//#include "hello.h"
import "C"

func main(){
    C.SayHello_in_c(C.CString("hello Golang"))
    C.SayHello_in_go(C.CString("hello clang"))
}
