指针 - unsafe 包的灵魂
Go版无类型指针和数值化的指针:

var p unsafe.Pointer = nil        // unsafe
var q uintptr        = uintptr(p) // builtin
C版无类型指针和数值化的指针:

void     *p = NULL;
uintptr_t q = (uintptr_t)(p); // <stdint.h>
unsafe.Pointer 是 Go指针 和 C指针 转换的中介
uintptr 是 Go 中 数值 和 指针 转换的中介

实战: int32 和 *C.char 相互转换(B)
// int32 => *C.char
var x = int32(9527)
var p *C.char = (*C.char)(unsafe.Pointer(uintptr(x)))

// *C.char => int32
var y *C.char
var q int32 = int32(uintptr(unsafe.Pointer(y)))
第一步: int32 => uintprt
第二步: uintptr => unsafe.Pointer
第三步: unsafe.Pointer => *C.char
反之亦然

实战: *X 和 *Y 相互转换(B)
var p *X
var q *Y

q = (*Y)(unsafe.Pointer(p)) // *X => *Y
p = (*X)(unsafe.Pointer(q)) // *Y => *X
第一步: *X => unsafe.Pointer
第二步: unsafe.Pointer => *Y
反之亦然

实战: []X 和 []Y 相互转换(B)
var p []X
var q []Y // q = p

pHdr := (*reflect.SliceHeader)(unsafe.Pointer(&p))
qHdr := (*reflect.SliceHeader)(unsafe.Pointer(&q))

pHdr.Data = qHdr.Data
pHdr.Len = qHdr.Len * unsafe.Sizeof(q[0]) / unsafe.Sizeof(p[0])
pHdr.Cap = qHdr.Cap * unsafe.Sizeof(q[0]) / unsafe.Sizeof(p[0])
所有切片拥有相同的头部 reflect.SliceHeader
重新构造切片头部即可完成转换

Go调用C函数
C调用Go导出函数
深度调用: Go => C => Go => C

CGO内部机制
CGO生成的中间文件
内部调用流程: Go -> C
内部调用流程: C -> Go
