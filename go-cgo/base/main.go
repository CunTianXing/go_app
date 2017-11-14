package main

/*
#include <stdlib.h>
*/
import "C"
import "fmt"

func main() {
	fmt.Println(random())
}

func random() int {
	return int(C.random())
}
