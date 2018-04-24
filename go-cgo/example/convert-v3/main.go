package main

import "reflect"
import "fmt"
import "unsafe"

func main() {
	var p []string
	var q []int // q = p

	pHdr := (*reflect.SliceHeader)(unsafe.Pointer(&p))
	qHdr := (*reflect.SliceHeader)(unsafe.Pointer(&q))

	pHdr.Data = qHdr.Data
	pHdr.Len = qHdr.Len * int(unsafe.Sizeof(q[0])) / int(unsafe.Sizeof(p[0]))
	pHdr.Cap = qHdr.Cap * int(unsafe.Sizeof(q[0])) / int(unsafe.Sizeof(p[0]))
	fmt.Println(unsafe.Sizeof(q[0]))
	fmt.Println(unsafe.Sizeof(p[0]))
	fmt.Printf("type %T,value %#v\n", pHdr, pHdr)
	fmt.Printf("type %T,value %#v\n", qHdr, qHdr)
}
