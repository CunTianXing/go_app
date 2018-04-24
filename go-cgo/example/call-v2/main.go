package main

/*
static int add(int a, int b){
    return a+b;    
}
*/
import "C"
import "fmt"

func main(){
    v, err := C.add(1,1)
    fmt.Println(v,err)
}
//任何C函数都可以带2个返回值
//第二个返回值是errno,对应error接口类型
