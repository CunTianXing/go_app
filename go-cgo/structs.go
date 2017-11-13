package main

/*
#include <stdlib.h>
struct point {
	char* name;
    int x;
    int y;
};
struct Person {
	char* name;
	int age;
	int height;
	int weight;
};
*/
import "C"
import "fmt"

func main() {
	p := C.struct_point{}
	p.name = C.CString("dddd")
	p.x = 99
	p.y = 42
	fmt.Printf("type:   %T\n", p)
	//type:   main._Ctype_struct_point
	fmt.Printf("struct: %+v\n", p)

	ps := C.struct_Person{}
	ps.name = C.CString("xingcuntian")
	ps.age = 23
	ps.height = 1000
	ps.weight = 12999
	fmt.Println(ps)
	fmt.Println(C.GoString(ps.name))
	fmt.Println(ps.age)
	fmt.Println(ps.height)
	fmt.Println(ps.weight)
}
