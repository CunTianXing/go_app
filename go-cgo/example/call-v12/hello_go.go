package main

//#include "hello.h"
import "C"
import "fmt"

//export SayHello_in_go
func SayHello_in_go(s *C.char) {
	fmt.Printf("SayHello_in_go: %s\n", C.GoString(s))
}
