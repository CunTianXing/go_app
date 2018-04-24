package main

import "C"
import "unsafe"
import "fmt"

 func main(){ 
    //int32 ===> *C.char
    var x = int32(99334)
    var p *C.char = (*C.char)(unsafe.Pointer(uintptr(x)))
    fmt.Printf("p type %T, %#v\n",p,p)//p type *main._Ctype_char, (*main._Ctype_char)(0x18406)

    // *C.char ===> int32
    var y *C.char
    var q int32 = int32(uintptr(unsafe.Pointer(y)))
    fmt.Printf("q type %T, %#v\n",q,q)
}
