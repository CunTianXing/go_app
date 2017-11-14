package main

/*
#include<stdio.h>
#include<stdlib.h>
double* cd_func(int n_rows,double *result) {
	for(int i = 0; i < n_rows; ++i) {
		result[i] = (double)i;
	}
	return result;
}
double* cdd_func(int n_rows) {
    double* result;
    result = calloc(n_rows, sizeof(double));
    for (int i = 0; i < n_rows; ++i) {
        result[i] = (double)i;
    }
    return result;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	size := 100
	f := doubleToFloats(C.cdd_func(C.int(size)), size)
	fmt.Println(f)
	doubleSize := unsafe.Sizeof(C.double(0))
	fmt.Printf("double length bytes:%d\n", doubleSize) //double length bytes:8
	buf := (*C.double)(C.calloc(C.size_t(size), C.size_t(doubleSize)))
	//farr := doubleToFloats(C.cd_func(C.int(size), buf), size)

	//fmt.Println(farr)

	data := (*[1 << 30]C.double)(unsafe.Pointer(C.cd_func(C.int(size), buf)))[:size:size]
	//data := (*[1 << 30]C.double)(unsafe.Pointer(C.cd_func(C.int(size), buf)))[:size]
	fmt.Println(data)
	fmt.Println("C *double to Go []float64")
	gf := doubleToFloatd(C.cd_func(C.int(size), buf), size)
	fmt.Println(gf)
}

func doubleToFloats(in *C.double, length int) []float64 {
	out := make([]float64, length, length)
	start := unsafe.Pointer(in)
	size := unsafe.Sizeof(C.double(0))
	for i := 0; i < length; i++ {
		val := *(*C.double)(unsafe.Pointer(uintptr(start) + size*uintptr(i)))
		out[i] = float64(val)
	}
	return out
}

func doubleToFloatd(in *C.double, size int) []float64 {
	defer C.free(unsafe.Pointer(in))
	out := (*[1 << 30]float64)(unsafe.Pointer(in))[:size:size]
	return out
}
