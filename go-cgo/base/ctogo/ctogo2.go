package main

/*
#include <stdio.h>
#include <stdlib.h>

void print(char* s) {
	printf("print: %s\n", s);
}
*/
import "C"
import (
	"fmt"
	"reflect"
	"runtime"
	"time"
	"unsafe"
)

type Slice struct {
	Data []byte
	data *c_slice_t
}

type c_slice_t struct {
	p unsafe.Pointer
	n int
}

func newSlice(p unsafe.Pointer, n int) *Slice {
	data := &c_slice_t{p, n}
	runtime.SetFinalizer(data, func(data *c_slice_t) {
		println("gc:", data.p)
		C.free(data.p)
	})
	s := &Slice{data: data}
	h := (*reflect.SliceHeader)((unsafe.Pointer(&s.Data)))
	h.Cap = n
	h.Len = n
	h.Data = uintptr(p)
	return s
}

func testSlice() {
	msg := "hello world!"
	p := C.calloc((C.size_t)(len(msg)+1), 1)
	println("malloc:", p)

	s := newSlice(p, len(msg)+1)
	copy(s.Data, []byte(msg))

	fmt.Printf("fmt.Printf: %s\n", string(s.Data))
	C.print((*C.char)(p))
}

func main() {
	testSlice()
	runtime.GC()
	runtime.Gosched()
	time.Sleep(1e9)
}
