package main

////////export int GoAdd(int a, int b);
//#include "add.h"
import "C"

func main() {
	//C.GoAdd(1, 1)
	C.c_add(2, 2)
}

//无法在Go文件应用导出的头文件，因为还未生成
//GoAdd是Go导出函数，无法通过_cgo_export.h引用
//c_add是C定义函数，可以通过add.h头文件引用
//可手写函数声明，不会形成循环依赖；
