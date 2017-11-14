package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	print("xingcuntian")
	d := make([]byte, 10)
	fmt.Println(d)
	s := []byte("123456")
	scopy(d, s, len(s))
	fmt.Println(d)
	fmt.Println(string(d))
}

func print(s string) {
	cs := C.CString(s)
	C.fputs(cs, (*C.FILE)(C.stdout))
	C.free(unsafe.Pointer(cs))
}

// Go语言和C语言通讯交互主要是通过传递参数和返回值. 其中参数和返回值除了基本的 数据类型外, 最重要的是如何相互传递/共享二进制的内存块.
// Go向C语言传递内存块
// 因为C语言的字符串结尾有\0, Go语言字符串没有\0,
// 因此需要重新构造一个C字符串. 其中 C.CString(s) 是构造一个C的字符串, 然后复制字符串并传入 C.fputs.
// 用完之后不要忘记调用C.free释放新创建的C字符串(可以用defer释放).

func scopy(dst, src []byte, size int) {
	C.memcpy(unsafe.Pointer(&dst[0]), unsafe.Pointer(&src[0]), C.size_t(size))
}

// 这个代码并没有涉及内存的创建/复制/删除等额外的操作, 是比较理想的集成方式.

// 注意: 在C语言使用该资源期间要防止Go语言的GC提前释放被C语言使用的Go内存!
