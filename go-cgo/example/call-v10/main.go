package main

//#include <math.h>
import "C"

func main(){
    println(int(C.pow(2,3)))
}
