package main

import "unsafe"
import "sort"
import "fmt"

func main(){
    // []float64 强制类型转换为[]int
    var a = []float64{4,5,2,1,7,8,9,1}
    var b []int = ((*[1 << 20]int))(unsafe.Pointer(&a[0]))[:len(a):cap(a)]
    //以int 方式给float64排序
    sort.Ints(b)
    fmt.Println(b)
}
