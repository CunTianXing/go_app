package main

/*
#include <string.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	//record := (**C.uchar)(unsafe.Pointer(uintptr(0)))
	// var record **C.uchar
	buf := []byte{1, 2, 3}
	fmt.Println(buf)
	record := (*C.uchar)(C.malloc(C.size_t(len(buf))))
	fmt.Println(record)
	C.memcpy(unsafe.Pointer(record), unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
	fmt.Println(*record)

}
