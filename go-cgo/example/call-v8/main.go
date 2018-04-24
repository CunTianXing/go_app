package main
/*
static int c_add(int a, int b){
    return a+b;    
}

static int go_add_proxy(int a, int b){
    extern int GoAdd(int a, int b);
    return GoAdd(a, b);
}
*/
import "C"

func main(){
    C.c_add(1,1)
}

//export GoAdd
func GoAdd(a, b C.int) C.int {
    return a + b
}
//深度调用: Go => C => Go => C (A)
//go_add_proxy 调用Go导出的 GoAdd
//Go:main => C:go_add_proxy => Go:GoAdd
