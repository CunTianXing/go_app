package main

/*
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
*/
import "C"
import "fmt"

func main() {
	f := new(C.FILE)
	fmt.Println(f)
	//p := C.CString("string content")
	//C.fputs(p, (*C.FILE)(f))
}
