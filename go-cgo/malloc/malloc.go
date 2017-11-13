package main

/*
#include <stdlib.h>
#include <stdint.h>
*/
import "C"
import "fmt"
import "unsafe"

func main() {
	size := 1024
	buf := (*C.char)(C.malloc(C.size_t(size)))
	defer C.free(unsafe.Pointer(buf))
	fmt.Println(unsafe.Sizeof(buf))
	fmt.Println(unsafe.Sizeof('a'))
	charSize := unsafe.Sizeof(new(C.char))
	fmt.Println(charSize)
}

// buf := (*C.char)(C.malloc((C.size_t)(len(data))))
// func cgo_malloc(n int) unsafe.Pointer {
// 	return C.malloc(C.size_t(n))
// }

// func cgo_free(ptr unsafe.Pointer) {
// 	C.free(ptr)
// }
// cpath := C.CString(path)
// cbuffer := (*C.char)(C.malloc(bufferSize))
// cbufferLen := C.int(bufferSize)
// defer C.free(unsafe.Pointer(cpath))
// defer C.free(unsafe.Pointer(cbuffer))

// buf := C.malloc(C.size_t(bufSize))
// defer C.free(buf)

// func cString(str string) (*C.char, C.size_t) {
// 	return C.CString(str), C.size_t(len(str))
// }

//cstring := (*C.char)(C.malloc(size))
//defer C.free(unsafe.Pointer(cstring))

//cpath := (*C.char)(C.malloc(C.size_t(128 * unsafe.Sizeof('a'))))
//defer C.free(unsafe.Pointer(cpath))

//key="ddffffffffffff"
// cs_key, key_len := cString(key)
// defer C.free(unsafe.Pointer(cs_key))
// cs_value, value_len := cString(string(buffer))
// defer C.free(unsafe.Pointer(cs_value))

// func cString(str string) (*C.char, C.size_t) {
// 	return C.CString(str), C.size_t(len(str))
// }

// var keys []string
// keys=.........
// char_size := unsafe.Sizeof(new(C.char))
// cs_keys := C.malloc(C.size_t(len(keys)) * C.size_t(char_size))
// defer C.free(cs_keys)

// len_size := unsafe.Sizeof(C.size_t(0))
// key_lens := C.malloc(C.size_t(len(keys)) * C.size_t(len_size))
// defer C.free(key_lens)

// buf := C.malloc(C.size_t(bufSize))
// defer C.free(buf)
