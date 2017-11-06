package main

/*
#include <stdio.h>
#include <stdint.h>
int ic = 5;
unsigned int uic = 7;
int16_t is = 12345;
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	var ig int = 10
	// C to GO
	igc := int(C.ic)
	fmt.Println("value:", igc, "type:", reflect.TypeOf(igc))
	// value: 5 type: int
	// GO to C
	icg := C.int(ig)
	fmt.Println("value:", icg, "type:", reflect.TypeOf(icg))
	//value: 10 type: main._Ctype_int

	icp := (*C.int)(unsafe.Pointer(&ig))
	fmt.Println("value:", reflect.ValueOf(icp), "type:", reflect.TypeOf(icp))
	//value: 0xc420016108 type: *main._Ctype_int

	uigc := uint(C.uic)
	fmt.Println("value: ", uigc, "type:", reflect.TypeOf(uigc))
	//value:  7 type: uint

	i16t := int16(C.is)
	fmt.Println("value:", i16t, reflect.TypeOf(i16t))
	//value: 12345 int16
}
