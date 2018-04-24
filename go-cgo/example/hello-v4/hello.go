package main

import "C"
import "fmt"

//SayHello export
func SayHello(s *C.char) {
	fmt.Print(C.GoString(s))
}
