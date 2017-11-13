package main

/*
#include <stdlib.h>
#include <stdint.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	size := 1024
	buf := cstrNew(size)
	fmt.Println(unsafe.Sizeof(buf))
	fmt.Println(C.size_t(1))
	charSize := unsafe.Sizeof(new(C.char))
	buf = (*C.char)(C.calloc(C.size_t(size), C.size_t(charSize)))
	defer C.free(unsafe.Pointer(buf))
	fmt.Println(unsafe.Sizeof(buf))
}

// Helper functions
// Calls C malloc
func malloc(size int) unsafe.Pointer {
	return (unsafe.Pointer(C.calloc(C.size_t(size), C.size_t(1))))
}

// Calls C free
func free(ptr unsafe.Pointer) {
	C.free(ptr)
}

// Allocates a string with the given byte length
// don't forget a call to defer s.free() !
func cstrNew(size int) *C.char {
	return (*C.char)(malloc(size))
}

// free is a method on C char * strings to method to free the associated memory
func (self *C.char) free() {
	free(unsafe.Pointer(self))
}

// free is a method on C int * pointers to method to free the associated memory
func (self *C.int) free() {
	C.free(unsafe.Pointer(self))
}

// cstring converts a string to a C string. This allocates memory,
// so don't forget to add a "defer s.free()"
func cstr(self string) *C.char {
	buf := cstrNew(len(self) + 1)
	// Allocate buffer
	if buf == nil {
		panic("Could not allocate memory for string")
	}
	// Some nice pointer math
	// for i:=0 ; i < len(self) ; i ++ {
	// ch  := self[i]
	// pto := (*byte)(ptr(uintptr(ptr(buf)) + uintptr(i)))
	// *pto = ch
	// }
	// // Don't forget to zero-terminate
	// ptoe := (*byte)(ptr(uintptr(ptr(buf)) + uintptr(len(self))))
	// *ptoe = byte(0)

	return buf

	// Strangely enough, C.String does NOT work for me. :p
	// return C.CString(self)
}

// Converts an int pointer to a C.int pointer
func cintptr(ptr *int) *C.int {
	return (*C.int)(unsafe.Pointer(ptr))
}

// Converts a byte pointer to a C.Uchar8 pointer
// func cbyteptr(ptr *uint8) *C.Uint8 {
// 	return (*C.Uint8)(unsafe.Pointer(ptr))
// }

/*
// cstring converts an int to a C int *. This allocates memory,
// so don't forget to add a "defer s.free()"
func cintptrNew(self int) (* C.char) {
      return (*C.int) unsafe.Pointer(C.malloc(C.size_t())))
        return C.CString(self)
}
*/
