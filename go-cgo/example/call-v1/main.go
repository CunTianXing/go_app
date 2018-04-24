package main

/*
static int add(int a, int b) {
    return a+b;    
}
*/
import "C"

func main(){
    C.add(1,1)
}
//C.add通过C的虚拟包访问
//最终会转为_Cfunc_add名字
