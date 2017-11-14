package main

/*
#include<stdio.h>
#include<stdlib.h>
double* c_func(int n_rows) {
	double* result;
	result = calloc(n_rows, sizeof(double));
	for (int i = 0; i < n_rows; ++i) {
	result[i] = (double)i;
	}
	return result;
}
*/
import "C"
import "unsafe"

func doubleToFloat(in *C.double, size int) []float64 {
	defer C.free(unsafe.Pointer(in))
	out := (*[1 << 30]float64)(unsafe.Pointer(in))[:size:size]
	return out
}
