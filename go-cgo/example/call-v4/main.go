package main

// static void noreturn(){}
import "C"
import "fmt"

func main(){
    x, _ := C.noreturn()
    fmt.Printf("%#v\n",x)

}
//main._Ctype_void{}
//可以获取一个void类型函数的返回值
//返回值类型：type _Ctype_void [0]byte
