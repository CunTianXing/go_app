package main

import "fmt"
import "unsafe"

func main(){
    var p  *string
    var q  *int

    q = (*int)(unsafe.Pointer(p)) // *string => *int
    fmt.Printf("q type %T, value %+v\n",q,q)
    p = (*string)(unsafe.Pointer(q)) // *int => *string
    fmt.Printf("p type %T, value %+v\n",p,p)
}
//q type *int, value <nil>
//p type *string, value <nil>
