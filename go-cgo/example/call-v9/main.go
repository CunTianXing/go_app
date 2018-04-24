package main

//int sum(int a, int b);
import "C"

func main() {
	C.sum(1, 2)
}

//export sum
func sum(a, b C.int) C.int {
	return a + b
}
