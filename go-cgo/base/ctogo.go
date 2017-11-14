package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

func main() {
	fmt.Println(getData())
	//fmt.Println(getDataP())
	fmt.Println(getDataPS())
}

func getData() []byte {
	size := 1024
	buf := C.malloc(C.size_t(size))
	defer C.free(buf)
	return C.GoBytes(buf, C.int(size))
}

// C向Go语言返回内存块

// 如果是C语言向Go返回内存块, 一般是先创建一个对应的Go的切片. 有现成的函数C.GoBytes()可以基于C的内存块构造切片.

// 比如获取C返回的内存块数据:
// 代码并不复杂. 但是效率并不理想: 其中需要新创建一个Go的切片, 并进行一次冗余的复制操作.

func getDataP() []byte {
	size := 1024
	buf := C.malloc(C.size_t(size))
	defer C.free(buf)
	var s []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	sh.Data = uintptr(buf)
	sh.Cap = size
	sh.Len = size
	return s
}

// type Pointer *ArbitraryType
// type SliceHeader struct {
//     Data uintptr
//     Len  int
//     Cap  int
// }

// SliceHeader代表一个运行时的切片。它不保证使用的可移植性、安全性；它的实现在未来的版本里也可能会改变。而且，Data字段也不能保证它指向的数据不会被当成垃圾收集，因此程序必须维护一个独立的、类型正确的指向底层数据的指针。

// 如果想去掉冗余的复制操作, 就需要基于C的内存块构造切片. 这个需要依赖Go语言的反射技术.

// 返回的s是基于C语言内存块构造的切片. 没有冗余的内存复制操作.

// 但是, 上面的代码却有内存泄漏的问题. Go语言的GC并不会自动释放C.malloc释放的内存.

// 如果需要Go语言的GC自动管理C语言返回的内存, 需要基于之前讲过的

// 简而言之, 就是要将C语言的内存块绑定到一个Go语言的内存资源, 然后依靠runtime.SetFinalizer的技术管理C语言的内存块.

// 核心代码如下:

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
		C.free(data.p)
	})
	s := &Slice{data: data}
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&s.Data))
	sh.Cap = n
	sh.Len = n
	sh.Data = uintptr(p)
	return s
}

func getDataPS() []byte {
	size := 10 * 1
	buf := C.malloc(C.size_t(size))
	s := newSlice(buf, size)
	//fmt.Println(s)
	return s.Data
}
