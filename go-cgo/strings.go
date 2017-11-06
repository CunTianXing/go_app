package main

/*
#include <stdlib.h>
char* cstring = "C string example";
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	var gstring string = "Go string example"

	// Go to C String, Outpit: *C.char
	cs := C.CString(gstring)
	defer C.free(unsafe.Pointer(cs))
	fmt.Println("value:", cs, "type:", reflect.TypeOf(cs))

	// C to Go String, Output: string
	gs := C.GoString(C.cstring)
	fmt.Println(gs)

	// C string, length to Go string
	gs2 := C.GoStringN(C.cstring, (C.int)(len(gs)))
	fmt.Println(gs2)

	gbyte := C.GoBytes(unsafe.Pointer(C.cstring), (C.int)(len(gs)))
	fmt.Println(gbyte)
	fmt.Println([]byte(gstring[:len(gs)]))

}

// func C.CString(goString string) *C.char
// func C.GoString(cString *C.char) string
// func C.GoStringN(cString *C.char, length C.int) string

// var cmsg *C.char = C.CString("hi")
// 	defer C.free(unsafe.Pointer(cmsg))

// var i uint64 = 0xdeedbeef01020304
// slice := (*[1 << 30]byte)(unsafe.Pointer(&i))[:8:8]
// fmt.Println(slice)

//func C.GoBytes(cArray unsafe.Pointer, length C.int) []byte
