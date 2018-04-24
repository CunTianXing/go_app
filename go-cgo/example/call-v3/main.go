package main

/*
#include <errno.h>
static void seterrno(int v) {
    errno = v;    
}
*/
import "C"
import "fmt"

func main(){
    _, err := C.seterrno(90033)
    fmt.Println(err)
}
//即使没有返回值，依然可以通过第二个返回值获取errno
//对应void类型的函数，第一个返回值可以用_占位
